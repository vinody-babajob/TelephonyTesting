package progressiveJobApplication

import (
	"BabajobIvrTesting/models"
	"BabajobIvrTesting/utils"
	bson "gopkg.in/mgo.v2/bson"
	"os"
	"testing"
	"time"

	"net/http"

	"encoding/json"

	"fmt"
)

func TestJobApplicaion(t *testing.T) {
	basePath := "/src/BabajobIvrTesting/scenarios/progressiveJobApplication/"

	configurationReader := utils.NewConfigurationReader(
		os.Getenv("TELEPHONY_ENV"),
		basePath,
	)

	userid, uerr := GetUserId("07338466702", configurationReader, t)
	if uerr != nil {
		t.Errorf("Error getting user %q", uerr.Error())
	}

	t.Log(userid)
	vn, verr := GetVirtualNumber(userid, configurationReader, t)
	if verr != nil {
		t.Log(verr.Error())
	}

	toNum := vn
	ctime := time.Now().UTC()
	callsid := "2344"
	callErr := MakeInBoundCall(toNum, callsid, configurationReader, t)
	if callErr != nil {
		t.Errorf("Failed while making inbound call %q", callErr.Error())
	}

	time.Sleep(time.Second * 10)

	_, mCallErr := IsInboundCallMadeForCallSid(callsid, configurationReader, t)
	if mCallErr != nil {
		t.Errorf("No inbound call found because %q", mCallErr.Error())
	}

	time.Sleep(time.Second * 20)
	_, ocallErr := CheckIfOutboundCallMade(ctime, configurationReader, t)
	if ocallErr != nil {
		t.Errorf("No outbound call found because %q", ocallErr.Error())
	}
}

func CheckIfOutboundCallMade(timeCheck time.Time, configurationReader utils.ConfigurationReader, t *testing.T) (bool, error) {
	outboundCollection := configurationReader.GetValue("outboundCallCollection")
	mongoConf := configurationReader.GetMapValue("mongo")
	mongoDb := utils.NewMongoDBWithConfig(mongoConf)

	outboundSession, outboundCallsRepo, outboundErr := mongoDb.GetCollection(outboundCollection)
	if outboundErr != nil {
		defer outboundSession.Close()
		t.Errorf("Trouble Getting repo for outboundCall with errors %q", outboundErr.Error())
		return false, outboundErr
	}
	outboundCalls := make([]models.OutboundCallRequest, 0, 1)
	selector := bson.M{
		"$and": []interface{}{
			bson.M{"mobilenumber": "07338466702"},
			bson.M{"CreatedAt": bson.M{"$gte": timeCheck}},
		},
	}
	err := outboundCallsRepo.Find(selector).All(&outboundCalls)
	if err != nil {
		defer outboundSession.Close()
		t.Errorf("Trouble Getting the  outbound call:%q with errors %q", "07338466702", err.Error())
		return false, err
	}
	defer outboundSession.Close()

	if len(outboundCalls) <= 0 {
		t.Errorf("No Entry for the  outbound call:%q", "07338466702")
		return false, fmt.Errorf("No outbound Made for users %q", "07338466702")
	}

	return true, nil
}

func MakeInBoundCall(toNum string, callsid string, configurationReader utils.ConfigurationReader, t *testing.T) error {
	inboundUrl := configurationReader.GetValue("inboundCallUrl")

	httpProvider := utils.NewHTTPProvider(inboundUrl, map[string]string{})
	queryParams := map[string]string{
		"From":          "07338466702",
		"To":            toNum,
		"CallSid":       callsid,
		"TransactionId": "2223",
	}
	_, err := httpProvider.Get(queryParams)
	if err != nil {
		t.Errorf("Error while Making inbound call with error %q", err.Error())
		return err
	}

	return nil
}

func IsInboundCallMadeForCallSid(callSid string, configurationReader utils.ConfigurationReader, t *testing.T) (bool, error) {
	mongoConf := configurationReader.GetMapValue("mongo")
	mongoDb := utils.NewMongoDBWithConfig(mongoConf)

	inboundCollection := configurationReader.GetValue("inboundCallCollection")
	session, inboundCallsRepo, err := mongoDb.GetCollection(inboundCollection)
	inboundCalls := make([]models.InboundCall, 0, 1)
	selector := bson.M{"providerCallIdentifier": callSid}
	err = inboundCallsRepo.Find(selector).All(&inboundCalls)
	if err != nil {
		defer session.Close()
		t.Errorf("Trouble Getting the callsid inboundcall:%q with errors %q", callSid, err.Error())
		return false, err
	}

	defer session.Close()

	if len(inboundCalls) <= 0 {
		t.Errorf("No Entry for the callsid inboundcall:%q", callSid)
		return false, fmt.Errorf("No InboundCalls Made for callsid %q", callSid)
	}

	return true, nil
}

func GetVirtualNumber(userId string, configReader utils.ConfigurationReader, t *testing.T) (string, error) {
	type User struct {
		BabaJobUserTypeId int
		BabajobUserType   int
		Id                int
	}

	type VirtualAssoc struct {
		Caller        User
		Callee        User
		VirtualNumber string `json:"virtualNumber"`
	}

	virtualNumUrl := configReader.GetValue("virtualNumberUrl")
	employerId := configReader.GetValue("employerId")
	jobId := configReader.GetValue("jobId")
	data := map[string]interface{}{
		"caller": map[string]interface{}{
			"id":                userId,
			"babaJobUserTypeId": 2,
		},
		"callee": map[string]interface{}{
			"id":                employerId,
			"babaJobUserTypeId": 1,
		},
		"purposeId":    1,
		"babajobJobId": jobId,
	}

	httpProvider := utils.NewHTTPProvider(virtualNumUrl, map[string]string{})
	postData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	response, rerr := httpProvider.Post(postData)
	if rerr != nil {
		return "", rerr
	}

	if response.StatusCode == http.StatusOK {
		vnum := []VirtualAssoc{}
		t.Log(string(response.Content))
		err = json.Unmarshal(response.Content, &vnum)
		if err != nil {
			return "", err
		}

		return vnum[0].VirtualNumber, nil
	}

	return "", nil
}

func GetUserId(userNumber string, configReader utils.ConfigurationReader, t *testing.T) (string, error) {
	type User struct {
		UserId string `json:"userId"`
		Role   string `json:"role"`
	}

	userUrl := configReader.GetValue("getUserUrl")
	consumerKey := configReader.GetValue("consumerKey")

	headers := map[string]string{
		"Consumer-Key": consumerKey,
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	urlParams := map[string]string{
		"mobileNumber": userNumber,
	}

	t.Log(userUrl)

	httpProvider := utils.NewHTTPProvider(userUrl, headers)
	res, err := httpProvider.Get(urlParams)
	if err != nil {
		t.Errorf("Error while getting user id for user with phonenumber %q with error %q", userNumber, err.Error())
		return "", err
	}

	user := User{}
	t.Log(string(res.Content))

	err = json.Unmarshal(res.Content, &user)
	if err != nil {
		t.Errorf("Json Error while getting user id for user with phonenumber %q with error %q", userNumber, err.Error())
		return "", err
	}

	return user.UserId, nil
}

package progressiveRegistration

import (
	"BabajobIvrTesting/models"
	"BabajobIvrTesting/utils"
	bson "gopkg.in/mgo.v2/bson"
	"os"
	"testing"
	"time"

	"fmt"
)

func TestSeekerRegistrationCall(t *testing.T) {
	basePath := "/src/BabajobIvrTesting/scenarios/progressiveRegistration/"

	configurationReader := utils.NewConfigurationReader(
		os.Getenv("TELEPHONY_ENV"),
		basePath,
	)

	ctime := time.Now().UTC()
	toNum := configurationReader.GetValue("seekerRegNumber")
	callsid := "23e3"
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

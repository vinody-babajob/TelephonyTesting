package progressiveJobApplication

import (
	"BabajobIvrTesting/models"
	"BabajobIvrTesting/utils"
	bson "gopkg.in/mgo.v2/bson"
	"os"
	"testing"
	"time"
)

func TestJobApplicaion(t *testing.T) {
	basePath := "/src/BabajobIvrTesting/scenarios/progressiveRegistration/"

	configurationReader := utils.NewConfigurationReader(
		os.Getenv("TELEPHONY_ENV"),
		basePath,
	)

	inboundUrl := configurationReader.GetValue("inboundCallUrl")
	toNum := configurationReader.GetValue("jobNumber")

	httpProvider := utils.NewHTTPProvider(inboundUrl, map[string]string{})
	queryParams := map[string]string{
		"From":          "07338466702",
		"To":            toNum,
		"CallSid":       "23e3",
		"TransactionId": "2223",
	}
	_, err := httpProvider.Get(queryParams)
	if err != nil {
		t.Errorf("Error while Making inbound call with error %q", err.Error())
	}

	time.Sleep(time.Second * 10)

	ctime := time.Now().UTC()

	mongoConf := configurationReader.GetMapValue("mongo")
	mongoDb := utils.NewMongoDBWithConfig(mongoConf)

	inboundCollection := configurationReader.GetValue("inboundCallCollection")
	session, inboundCallsRepo, err := mongoDb.GetCollection(inboundCollection)
	inboundCalls := make([]models.InboundCall, 0, 1)
	selector := bson.M{"providerCallIdentifier": "23e3"}
	err = inboundCallsRepo.Find(selector).All(&inboundCalls)
	if err != nil {
		defer session.Close()
		t.Errorf("Trouble Getting the callsid inboundcall:%q with errors %q", "23e3", err.Error())
	}

	defer session.Close()

	if len(inboundCalls) <= 0 {
		t.Errorf("No Entry for the callsid inboundcall:%q", "23e3")
	}

	time.Sleep(time.Second * 20)

	outboundCollection := configurationReader.GetValue("outboundCallCollection")
	outboundSession, outboundCallsRepo, outboundErr := mongoDb.GetCollection(outboundCollection)
	if outboundErr != nil {
		defer outboundSession.Close()
		t.Errorf("Trouble Getting repo for outboundCall with errors %q", outboundErr.Error())
	}
	outboundCalls := make([]models.OutboundCallRequest, 0, 1)
	selector = bson.M{"MobileNumber": "07338466702", "CreatedAt": bson.M{"$gt": ctime.String()}}
	err = outboundCallsRepo.Find(selector).All(&outboundCalls)
	if err != nil {
		defer outboundSession.Close()
		t.Errorf("Trouble Getting the  outbound call:%q with errors %q", "07338466702", err.Error())
	}
	defer outboundSession.Close()

	if len(outboundCalls) <= 0 {
		t.Errorf("No Entry for the  outbound call:%q", "07338466702")
	}
}

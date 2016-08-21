package models

import (
	"time"
)

type InboundCall struct {
	FromNumber             string    `bson:"fromNumber"`
	ToNumber               string    `bson:"toNumber"`
	ProviderID             int       `bson:"providerID"`
	ProviderCallIdentifier string    `bson:"providerCallIdentifier"`
	CreatedAt              time.Time `bson:"createdAt"`
	TransactionId          string    `bson:transactionId` //This is the GUID generated by us
}

type Purpose struct {
	ID   int    `bson:"PurposeID"`
	Name string `bson:"PurposeName"`
}

type OutboundCallRequest struct {
	BabajobUserID string `bson:"BabajobUserId"`
	MobileNumber  string
	Purpose       Purpose   `bson:"Purpose"`
	CreatedAt     time.Time `bson:"CreatedAt"`
	TransactionId string    `bson:"TransactionId"`
	Delay         int       // This is in milli seconds
}

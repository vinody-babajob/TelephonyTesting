package ProgressiveRegistration

import (
	"net/http"
	"testing"
)

func TestInboundCall(t *testing.T) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", newUrl.String(), nil)
	if err != nil {
		return httpResponse, errors.New("Could not create the GET request : " + err.Error())
	}
}

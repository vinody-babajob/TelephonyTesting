package httpProvider

import (
	"Babajob/Telephony.Utilities/logger"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HTTPProvider struct {
	URL     string
	Headers map[string]string
	logger  logger.Logger
}

func (httpProvider *HTTPProvider) Get(queryParams map[string]string) (HTTPResponse, error) {
	var httpResponse HTTPResponse

	newUrl, err := url.Parse(httpProvider.URL)
	if err != nil {
		return httpResponse, errors.New("Could not parse the URL :" + err.Error())
	}

	parameters := url.Values{}
	for parameterKey, parameterValue := range queryParams {
		parameters.Add(parameterKey, parameterValue)
	}
	newUrl.RawQuery = parameters.Encode()

	request, err := http.NewRequest("GET", newUrl.String(), nil)
	if err != nil {
		return httpResponse, errors.New("Could not create the GET request : " + err.Error())
	}

	for headerKey, value := range httpProvider.Headers {
		request.Header.Set(headerKey, value)
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return httpResponse, errors.New("Failed to create the response : " + err.Error())
	}

	defer response.Body.Close()

	httpResponse.StatusCode = response.StatusCode
	httpResponse.Content, err = ioutil.ReadAll(response.Body)
	httpResponse.Headers = response.Header

	return httpResponse, err
}

func (httpProvider *HTTPProvider) Post(data []byte) (HTTPResponse, error) {
	var httpResponse HTTPResponse

	request, err := http.NewRequest("POST", httpProvider.URL, bytes.NewBuffer(data))

	if err != nil {
		return httpResponse, errors.New("Could not create the POST request : " + err.Error())
	}

	for headerKey, value := range httpProvider.Headers {
		request.Header.Set(headerKey, value)
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return httpResponse, errors.New("Could not create the response : " + err.Error())
	}

	if response != nil {
		defer response.Body.Close()

		httpResponse.StatusCode = response.StatusCode
		httpResponse.Content, err = ioutil.ReadAll(response.Body)

		if err != nil {
			return httpResponse, errors.New("IO READ ERROR : " + err.Error())
		}
	}

	httpResponse.StatusCode = response.StatusCode
	httpResponse.Headers = response.Header

	if err != nil {
		httpProvider.logger.LogError(err.Error())
	}

	return httpResponse, nil
}

func (httpProvider *HTTPProvider) Delete() (HTTPResponse, error) {
	var httpResponse HTTPResponse

	request, err := http.NewRequest("DELETE", httpProvider.URL, nil)

	if err != nil {
		return httpResponse, errors.New("Could not create the DELETE request : " + err.Error())
	}

	for headerKey, value := range httpProvider.Headers {
		request.Header.Set(headerKey, value)
	}

	httpProvider.logger.LogInfo(fmt.Sprintf("%+v", request))
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return httpResponse, errors.New("Could not create the response : " + err.Error())
	}

	if response != nil {
		defer response.Body.Close()

		httpResponse.StatusCode = response.StatusCode

		httpResponse.Content, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return httpResponse, errors.New("IO READ ERROR : " + err.Error())
		}
	}

	httpResponse.StatusCode = response.StatusCode
	httpResponse.Headers = response.Header
	httpResponse.Content, _ = ioutil.ReadAll(response.Body)

	return httpResponse, nil
}

func NewHTTPProvider(url string, header map[string]string) *HTTPProvider {
	return &HTTPProvider{
		URL:     url,
		Headers: header,
		logger:  logger.NewLogger(),
	}
}

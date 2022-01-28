package logging

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/jianyuan/go-sentry/sentry"
)

const dummyAuthToken string = "Bearer grizzliesandpolarbears"

func dummyHttpResponseStruct() *http.Response {
	return &http.Response{
		Request: &http.Request{
			URL: &url.URL{
				Scheme: "https",
				Host:   "thisissentry.com",
			},
			Header: http.Header{
				"thing":         []string{"value"},
				"Authorization": []string{dummyAuthToken},
			},
		},
		Status: "200 OK",
		Header: http.Header{
			"thing": []string{"value"},
		},
	}
}

func TestLogHttpResponseLogsPredictably(t *testing.T) {
	cases := []struct {
		testCaseName       string
		responseData       interface{}
		expectedDataString string
	}{
		{
			testCaseName: "Test with random struct",
			responseData: struct {
				Stuff string `json:"stuff"`
			}{Stuff: "hello"},
			expectedDataString: `"stuff": "hello"`,
		},
		{
			testCaseName: "Test with Sentry org",
			responseData: &sentry.Team{
				ID:   "sentry",
				Slug: "sentry stuff",
				Name: "still sentry stuff",
			},
			expectedDataString: `        "id": "sentry",
        "slug": "sentry stuff",
        "name": "still sentry stuff"`,
		},
		{
			testCaseName:       "Test with nothing",
			responseData:       nil,
			expectedDataString: "null",
		},
	}
	for _, tCase := range cases {
		t.Run(tCase.testCaseName, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			resp := dummyHttpResponseStruct()
			LogHttpResponse(resp, tCase.responseData, TraceLevel)
			if !strings.Contains(buf.String(), tCase.expectedDataString) {
				t.Logf("Log string doesn't contain the expected string...")
				t.Logf("Expected to contain:\n%s", tCase.expectedDataString)
				t.Logf("Full log:\n%s", buf.String())
				t.FailNow()
			}
		})
	}

}

func TestLogHttpResponseRedactsAuthNonDestructively(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	resp := dummyHttpResponseStruct()
	var originalHeaders map[string][]string = resp.Request.Header
	dummyRespData := &sentry.Team{
		ID:   "sentry",
		Slug: "sentry stuff",
		Name: "still sentry stuff",
	}
	LogHttpResponse(resp, dummyRespData, TraceLevel)
	if !reflect.DeepEqual(originalHeaders["Authorization"], resp.Request.Header["Authorization"]) {
		t.Fatalf("The header was changed by the function call. This should not happen.")
	}
	if strings.Contains(buf.String(), dummyAuthToken) {
		t.Logf("The logs should not contain secrets like the Authorization' Header value.")
		t.Fatalf("The logs contain the test dummy token %s", dummyAuthToken)
	}
}

package logging

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jianyuan/go-sentry/sentry"
)

func mockHttpResponseStruct() *http.Response {
	return &http.Response{
		Request: &http.Request{
			URL: &url.URL{
				Scheme: "https",
				Host:   "thisissentry.com",
			},
			Header: http.Header{
				"thing": []string{"value"},
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
			resp := mockHttpResponseStruct()
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

package logging

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestLogFunctionsLogAsExpected(t *testing.T) {
	cases := []struct {
		testCaseName   string
		inputArgs      []interface{}
		expectedString string
		logFunc        func(args ...interface{})
	}{
		{
			testCaseName:   "Test Info",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[INFO] Hello world this is fun\n",
			logFunc:        Info,
		},
		{
			testCaseName:   "Test Debug",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[DEBUG] Hello world this is fun\n",
			logFunc:        Debug,
		},
		{
			testCaseName:   "Test Warning",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[WARN] Hello world this is fun\n",
			logFunc:        Warning,
		},
		{
			testCaseName:   "Test Error",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[ERROR] Hello world this is fun\n",
			logFunc:        Error,
		},
		{
			testCaseName:   "Test Trace",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[TRACE] Hello world this is fun\n",
			logFunc:        Trace,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.testCaseName, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			tCase.logFunc(tCase.inputArgs...)
			compareLogsToExpected(t, tCase.expectedString, buf.String())
		})
	}
}

func TestLogFormatFunctionsLogAsExpected(t *testing.T) {
	cases := []struct {
		testCaseName   string
		format         string
		inputArgs      []interface{}
		expectedString string
		logFunc        func(format string, args ...interface{})
	}{
		{
			testCaseName:   "Test Infof",
			format:         "%s %s %s",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[INFO] Hello world this is fun\n",
			logFunc:        Infof,
		},
		{
			testCaseName:   "Test Debugf",
			format:         "%s %s %s",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[DEBUG] Hello world this is fun\n",
			logFunc:        Debugf,
		},
		{
			testCaseName:   "Test Warningf",
			format:         "%s %s %s",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[WARN] Hello world this is fun\n",
			logFunc:        Warningf,
		},
		{
			testCaseName:   "Test Errorf",
			format:         "%s %s %s",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[ERROR] Hello world this is fun\n",
			logFunc:        Errorf,
		},
		{
			testCaseName:   "Test Tracef",
			format:         "%s %s %s",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[TRACE] Hello world this is fun\n",
			logFunc:        Tracef,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.testCaseName, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			tCase.logFunc(tCase.format, tCase.inputArgs...)
			compareLogsToExpected(t, tCase.expectedString, buf.String())
		})
	}
}

func compareLogsToExpected(t *testing.T, expected, actual string) {
	if !strings.EqualFold(expected, actual) {
		t.Logf("Log string doesn't match the expected string...")
		t.Logf("Expected:\n%s", expected)
		t.Logf("Actual:\n%s", actual)
		t.FailNow()
	}
}

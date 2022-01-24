package logging

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestLogFunctionsLogAsExpected(t *testing.T) {
	cases := []struct {
		name           string
		inputArgs      []interface{}
		expectedString string
		logFunc        func(args ...interface{})
	}{
		{
			name:           "Test Info",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[INFO] Hello world this is fun\n",
			logFunc:        Info,
		},
		{
			name:           "Test Debug",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[DEBUG] Hello world this is fun\n",
			logFunc:        Debug,
		},
		{
			name:           "Test Warning",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[WARN] Hello world this is fun\n",
			logFunc:        Warning,
		},
		{
			name:           "Test Error",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[ERROR] Hello world this is fun\n",
			logFunc:        Error,
		},
		{
			name:           "Test Trace",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[TRACE] Hello world this is fun\n",
			logFunc:        Trace,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			tCase.logFunc(tCase.inputArgs...)
			compareLogsToExpected(t, tCase.expectedString, buf.String())
		})
	}
}

func TestLogFormatFunctionsLogAsExpected(t *testing.T) {
	cases := []struct {
		name           string
		format         string
		inputArgs      []interface{}
		expectedString string
		logFunc        func(format string, args ...interface{})
	}{
		{
			name:           "Test Infof",
			format:         "%s %s %s",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[INFO] Hello world this is fun\n",
			logFunc:        Infof,
		},
		{
			name:           "Test Debugf",
			format:         "%s %s %s",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[DEBUG] Hello world this is fun\n",
			logFunc:        Debugf,
		},
		{
			name:           "Test Warningf",
			format:         "%s %s %s",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[WARN] Hello world this is fun\n",
			logFunc:        Warningf,
		},
		{
			name:           "Test Errorf",
			format:         "%s %s %s",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[ERROR] Hello world this is fun\n",
			logFunc:        Errorf,
		},
		{
			name:           "Test Tracef",
			format:         "%s %s %s",
			inputArgs:      []interface{}{"Hello world", "this is", "fun"},
			expectedString: "[TRACE] Hello world this is fun\n",
			logFunc:        Tracef,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			tCase.logFunc(tCase.format, tCase.inputArgs...)
			compareLogsToExpected(t, tCase.expectedString, buf.String())
		})
	}
}

func compareLogsToExpected(t *testing.T, expected, actual string) {
	if !strings.EqualFold(expected, actual) {
		t.Logf("Log strings don't match the expected strings...")
		t.Logf("Expected:\n%s", expected)
		t.Logf("Actual:\n%s", actual)
		t.FailNow()
	}
}

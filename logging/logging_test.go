package logging

import (
	"io"
	"log"
	"os"
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
			reader, writer := getReadWriter(t)
			log.SetOutput(writer)
			tCase.logFunc(tCase.inputArgs...)
			compareLogsToExpected(t, tCase.expectedString, reader)
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
			reader, writer := getReadWriter(t)
			log.SetOutput(writer)
			tCase.logFunc(tCase.format, tCase.inputArgs...)
			compareLogsToExpected(t, tCase.expectedString, reader)
		})
	}
}

func compareLogsToExpected(t *testing.T, expected string, reader io.Reader) {
	buffer := make([]byte, 50) // len of Hello world this is fun
	actualLen, err := reader.Read(buffer)
	if err != nil {
		t.Fatalf("couldn't read into buffer: %v", err)
	}
	actualText := string(buffer[:actualLen])
	if !strings.EqualFold(expected, actualText) {
		t.Logf("Log strings don't match the expected strings...")
		t.Logf("Expected:\n%s", expected)
		t.Logf("Actual:\n%s", actualText)
		t.FailNow()
	}
}

func getReadWriter(t *testing.T) (reader io.Reader, writer io.Writer) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("couldn't get os Pipe: %v", err)
	}
	return reader, writer
}

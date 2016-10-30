package helpers

import (
	"net/http"
	"net/http/httputil"
)

type Logger interface {
	Printf(format string, a ...interface{})
}

func LogMark(mark string, logger Logger) {
	logger.Printf("// %s\n", mark)
}

func DumpResponse(resp *http.Response, logger Logger) {
	dump, err := httputil.DumpResponse(resp, false)
	if err != nil {
		logger.Printf("dumpResponse err: %s\n", err.Error())
		return
	}
	logger.Printf("Response dump: %s\n", string(dump))
}

func DumpRequest(req *http.Request, logger Logger) {
	dump, err := httputil.DumpRequestOut(req, false)
	if err != nil {
		logger.Printf("dumpRequest err: %s", err.Error())
		return
	}
	logger.Printf("Request dump: %s\n", string(dump))
}

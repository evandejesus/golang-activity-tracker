package server

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/kjk/common/siser"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

var (
	muLogHTTP sync.Mutex
)

// LogReqInfo describes info about HTTP request
type HTTPReqInfo struct {
	// GET etc.
	method  string
	uri     string
	referer string
	ipaddr  string
	// response code, like 200, 404
	code int
	// number of bytes of the response sent
	size int64
	// how long did it take to
	duration  time.Duration
	userAgent string
}

var Logger zerolog.Logger

func newLogger() zerolog.Logger {
	z := zerolog.New(&lumberjack.Logger{
		Filename:   "http_access.log", // File name
		MaxSize:    100,               // Size in MB before file gets rotated
		MaxBackups: 5,                 // Max number of files kept before being overwritten
		MaxAge:     30,                // Max number of days to keep the files
		Compress:   true,              // Whether to compress log files using gzip
	})
	return z.With().Caller().Timestamp().Logger()

}

func logRequestHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		ri := HTTPReqInfo{
			method:    r.Method,
			uri:       r.URL.String(),
			referer:   r.Header.Get("Referer"),
			userAgent: r.Header.Get("User-Agent"),
		}

		ri.ipaddr = requestGetRemoteAddress(r)
		m := httpsnoop.CaptureMetrics(h, w, r)

		ri.code = m.Code
		ri.size = m.Written
		ri.duration = m.Duration
		logHTTPReq(&ri)

	}
	return http.HandlerFunc(fn)
}

func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// requestGetRemoteAddress returns ip address of the client making the request,
// taking into account http proxies
func requestGetRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	}
	return hdrRealIP
}

func logHTTPReq(ri *HTTPReqInfo) error {
	w, err := createWriter()
	if err != nil {
		return err
	}

	var rec siser.Record
	rec.Name = "httplog"
	rec.Write("method", ri.method)
	rec.Write("uri", ri.uri)
	if ri.referer != "" {
		rec.Write("referer", ri.referer)
	}
	rec.Write("ipaddr", ri.ipaddr)
	rec.Write("code", strconv.Itoa(ri.code))
	rec.Write("size", strconv.FormatInt(ri.size, 10))
	dur := ri.duration / time.Millisecond
	rec.Write("duration", strconv.FormatInt(int64(dur), 10))
	rec.Write("ua", ri.userAgent)

	muLogHTTP.Lock()
	defer muLogHTTP.Unlock()
	_, err = w.WriteRecord(&rec)
	return err
}

func createWriter() (*siser.Writer, error) {
	f, err := os.Create("http_access.log")
	if err != nil {
		return nil, err
	}
	w := siser.NewWriter(f)
	return w, nil
}

package auth

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

// DebugTransport :
type DebugTransport struct {
	Base    http.RoundTripper
	Verbose bool
}

// RoundTrip :
func (t *DebugTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	b, err := httputil.DumpRequest(request, t.Verbose)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stderr, string(b))
	response, err := t.Base.RoundTrip(request)
	if err != nil {
		return nil, err
	}
	b2, err := httputil.DumpResponse(response, t.Verbose)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stderr, string(b2))
	return response, err
}

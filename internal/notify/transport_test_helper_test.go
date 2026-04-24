package notify

import (
	"net/http"
	"net/url"
)

// rewriteTransport returns an http.RoundTripper that redirects all requests
// to the given base URL, preserving path and query. It is used across notifier
// tests to intercept outbound HTTP calls without modifying production code.
func rewriteTransport(baseURL string) http.RoundTripper {
	return &hostRewriter{base: baseURL}
}

type hostRewriter struct {
	base string
}

func (h *hostRewriter) RoundTrip(req *http.Request) (*http.Response, error) {
	base, err := url.Parse(h.base)
	if err != nil {
		return nil, err
	}
	cloned := req.Clone(req.Context())
	cloned.URL.Scheme = base.Scheme
	cloned.URL.Host = base.Host
	return http.DefaultTransport.RoundTrip(cloned)
}

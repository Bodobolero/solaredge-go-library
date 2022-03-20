package solaredge

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	DEFAULT_URL = "https://monitoringapi.solaredge.com"
)

// SEOpts is a options type for the client.
type SEOpt func(sec *SEClient)

// A SEClient can call solaredge API's.
type SEClient struct {
	apikey  string
	baseurl string
	client  *http.Client
}

// SiteClient wraps a site and contains site specific methods.
type SiteClient struct {
	*SEClient
	siteid string
}

// NewClient returns a SEClient. You must supply an API key to access the solaredge API.
func NewClient(apikey string) *SEClient {
	return &SEClient{
		apikey: apikey,
		client: http.DefaultClient,
	}
}

// WithBaseURL is an option for the SEClient to change the URL of the solaredge API.
func WithBaseURL(u string) SEOpt {
	return func(c *SEClient) {
		c.baseurl = u
	}
}

// NewSite returns a SiteClient with the given site-ID.
func (sec *SEClient) NewSite(sid string) *SiteClient {
	return &SiteClient{
		SEClient: sec,
		siteid:   sid,
	}
}

// SiteFromIDs return a SiteClient with the given apikey and siteid.
func SiteFromIDs(apikey, siteid string, opts ...SEOpt) (*SiteClient, error) {
	cl := NewClient(apikey)
	for _, o := range opts {
		o(cl)
	}
	if cl.baseurl == "" {
		cl.baseurl = DEFAULT_URL
	}
	return cl.NewSite(siteid), nil
}

func (sec *SEClient) get(path string, parms url.Values, target any) error {
	if parms == nil {
		parms = make(url.Values)
	}
	parms.Add("api_key", sec.apikey)
	url := fmt.Sprintf("%s%s?%s", sec.baseurl, path, parms.Encode())
	rq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}
	rsp, err := sec.client.Do(rq)
	if err != nil {
		return fmt.Errorf("cannot invoke request: %w", err)
	}
	defer rsp.Body.Close()
	return json.NewDecoder(rsp.Body).Decode(target)
}

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

type SEOpt func(sec *SEClient)

type SEClient struct {
	apikey  string
	baseurl string
	client  *http.Client
}

type SiteClient struct {
	*SEClient
	siteid string
}

func NewClient(apikey string) *SEClient {
	return &SEClient{
		apikey: apikey,
		client: http.DefaultClient,
	}
}

func WithBaseURL(u string) SEOpt {
	return func(c *SEClient) {
		c.baseurl = u
	}
}

func (sec *SEClient) NewSite(sid string) *SiteClient {
	return &SiteClient{
		SEClient: sec,
		siteid:   sid,
	}
}

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

func (sec *SEClient) Get(path string, parms url.Values, target any) error {
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

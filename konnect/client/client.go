package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	// FormEncoded              = "application/x-www-form-urlencoded"
	// ApplicationJson          = "application/json"
	// ApplicationXml           = "application/xml"
	// IdSeparator              = ":"
	// Basic                    = "Basic"
	Bearer        = "Bearer"
	KonnectDomain = "api.konghq.com"
	GlobalRegion  = "global"
	FilterName    = "filter[name]"
	// Search                   = "$search"
	// SearchValue              = "\"%s:%s\""
	// Filter                   = "$filter"
	// FilterValue              = "%s eq '%s'"
	// FilterAnd                = " and "
	// Public                   = "Public"
	// Private                  = "Private"
	// WaitNotExists            = "NotExists"
	// WaitFound                = "Found"
	// WaitError                = "Error"
)

type Client struct {
	pat        string
	region     string
	httpClient *http.Client
}

func NewClient(ctx context.Context, pat string, region string) (client *Client, err error) {
	c := &Client{
		pat:        pat,
		region:     region,
		httpClient: &http.Client{},
	}
	return c, nil
}

func (c *Client) HttpRequest(ctx context.Context, isRegion bool, method string, path string, query url.Values, headerMap http.Header, body *bytes.Buffer) (response *bytes.Buffer, err error) {
	req, err := http.NewRequest(method, c.RequestPath(isRegion, path), body)
	if err != nil {
		return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	//Handle query values
	if query != nil {
		requestQuery := req.URL.Query()
		for key, values := range query {
			for _, value := range values {
				requestQuery.Add(key, value)
			}
		}
		req.URL.RawQuery = requestQuery.Encode()
	}
	//Handle header values
	if headerMap != nil {
		for key, values := range headerMap {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}
	//Handle authentication
	if c.pat != "" {
		req.Header.Set(headers.Authorization, Bearer+" "+c.pat)
	}
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		tflog.Info(ctx, "Konnect API:", map[string]any{"error": err})
	} else {
		tflog.Info(ctx, "Konnect API: ", map[string]any{"request": string(requestDump)})
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	defer resp.Body.Close()
	respBody := new(bytes.Buffer)
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		return nil, &RequestError{StatusCode: resp.StatusCode, Err: err}
	}
	if (resp.StatusCode < http.StatusOK) || (resp.StatusCode >= http.StatusMultipleChoices) {
		return nil, &RequestError{StatusCode: resp.StatusCode, Err: fmt.Errorf("%s", respBody.String())}
	}
	return respBody, nil
}

func (c *Client) RequestPath(isRegion bool, path string) string {
	var host string
	if isRegion {
		host = GlobalRegion
	} else {
		host = c.region
	}
	return fmt.Sprintf("https://%s.%s/v2/%s", host, KonnectDomain, path)
}

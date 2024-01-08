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
	ApplicationJson              = "application/json"
	Bearer                       = "Bearer"
	KonnectDomain                = "api.konghq.com"
	GlobalRegion                 = "global"
	IdSeparator                  = ":"
	FilterName                   = "filter[name]"
	FilterNameContains           = "filter[name][contains]"
	FilterFullName               = "filter[full_name]"
	FilterFullNameContains       = "filter[full_name][contains]"
	FilterEmail                  = "filter[email]"
	FilterEmailContains          = "filter[email][contains]"
	FilterActive                 = "filter[active]"
	FilterRoleName               = "filter[role_name]"
	FilterRoleNameContains       = "filter[role_name][contains]"
	FilterEntityTypeName         = "filter[entity_type_name]"
	FilterEntityTypeNameContains = "filter[entity_type_name][contains]"
)

type EntityId struct {
	Id string `json:"id"`
}

type Client struct {
	pat    string
	region string
	//defaultTags []string
	httpClient *http.Client
}

// func NewClient(pat string, region string, defaultTags []string) (*Client, error) {
func NewClient(pat string, region string) (*Client, error) {
	c := &Client{
		pat:    pat,
		region: region,
		//defaultTags: defaultTags,
		httpClient: &http.Client{},
	}
	return c, nil
}

func (c *Client) HttpRequest(ctx context.Context, isRegion bool, method string, path string, query url.Values, headerMap http.Header, body *bytes.Buffer) (*bytes.Buffer, error) {
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
		host = c.region
	} else {
		host = GlobalRegion
	}
	return fmt.Sprintf("https://%s.%s/%s", host, KonnectDomain, path)
}

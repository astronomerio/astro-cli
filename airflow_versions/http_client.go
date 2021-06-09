package airflowversions

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/astronomer/astro-cli/pkg/httputil"
	"github.com/pkg/errors"
)

// Client containers the logger and HTTPClient used to communicate with the HoustonAPI
type Client struct {
	HTTPClient *httputil.HTTPClient
}

// NewClient returns a new Client with the logger and HTTP client setup.
func NewClient(c *httputil.HTTPClient) *Client {
	return &Client{
		HTTPClient: c,
	}
}

type Request struct{}

// Do (request) is a wrapper to more easily pass variables to a client.Do request
func (r *Request) DoWithClient(api *Client) (*Response, error) {
	doOpts := httputil.DoOptions{
		Headers: map[string]string{
			"Accept": "application/json",
		},
	}

	return api.Do(doOpts)
}

// Do (request) is a wrapper to more easily pass variables to a client.Do request
func (r *Request) Do() (*Response, error) {
	return r.DoWithClient(NewClient(httputil.NewHTTPClient()))
}

// Do executes a query against the updates astronomer API, logging out any errors contained in the response object
func (c *Client) Do(doOpts httputil.DoOptions) (*Response, error) {
	var response httputil.HTTPResponse
	// FIXME: move to config file
	httpResponse, err := c.HTTPClient.Do("GET", "https://updates.astronomer.io/astronomer-certified", &doOpts)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	// strings.NewReader(jsonStream)
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	response = httputil.HTTPResponse{
		Raw:  httpResponse,
		Body: string(body),
	}
	decode := Response{}
	err = json.NewDecoder(strings.NewReader(response.Body)).Decode(&decode)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to JSON decode Houston response")
	}

	return &decode, nil
}
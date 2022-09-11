package cpe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	http   *http.Client
	apiKey string
}

func (c *Client) FetchAll(query string) (*CpeResponse, error) {
	endpoint := fmt.Sprintf("https://services.nvd.nist.gov/rest/json/cpes/1.0/?keyword=%s", url.QueryEscape(query))

	resp, err := c.http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	res := &CpeResponse{}
	return res, json.Unmarshal(body, res)
}

type CpeResponse struct {
	ResultsPerPage int    `json:"resultsPerPage"`
	StartIndex     int    `json:"startIndex"`
	TotalResults   int    `json:"totalResults"`
	Result         Result `json:"result"`
}

type Result struct {
	DataType      string `json:"dataType"`
	FeedVersion   string `json:"feedVersion"`
	CpeCount      int    `json:"cpeCount"`
	FeedTimestamp string `json:"feedTimestamp"`
	Cpes          []Cpe  `json:"cpes"`
}

type Cpe struct {
	Deprecated       bool          `json:"deprecated"`
	Cpe23URI         string        `json:"cpe23Uri"`
	LastModifiedDate string        `json:"lastModifiedDate"`
	Titles           []Title       `json:"titles"`
	Refs             []interface{} `json:"refs"`
	DeprecatedBy     []interface{} `json:"deprecatedBy"`
	Vulnerabilities  []interface{} `json:"vulnerabilities"`
}

type Title struct {
	Title string `json:"title"`
	Lang  string `json:"lang"`
}

func (c *Cpe) GetTitle() string {
	return c.Titles[0].Title
}
// ===========================

func NewClient(httpClient *http.Client, apiKey string) *Client {
	return &Client{httpClient, apiKey}
}

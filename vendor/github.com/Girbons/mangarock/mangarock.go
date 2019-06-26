package mangarock

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIEndpoint is the mangarock api endpoint :)
const APIEndpoint = "https://api.mangarockhd.com/query/web401/"

// MangaRockInfo is used to parse a response from the `info` endpoint
// Contains the information related to a Manga
type MangaRockInfo struct {
	Code int   `json:"code"`
	Data Manga `json:"data"`
}

// MangaRockPages is used to parse a response from `pages` endpoint
// Contains the link to the images related to a manga
type MangaRockPages struct {
	Code int      `json:"code"`
	Data []string `json:"data"`
}

// Client contains only the `client` by now
// Maybe in future it can contains an ApiKey
type Client struct {
	Client  *http.Client
	Options map[string]string
}

// NewClient returns a Client instance
func NewClient() *Client {
	return &Client{
		Client:  &http.Client{},
		Options: make(map[string]string),
	}
}

// NewClientWithOptions returns a Client instance with the given options
func NewClientWithOptions(options map[string]string) *Client {
	return &Client{
		Client:  &http.Client{},
		Options: options,
	}
}

// getJson decode the response body to given struct
func getJSON(response *http.Response, target interface{}) error {
	return json.NewDecoder(response.Body).Decode(target)
}

// Set Additional options to the client instance
func (c *Client) SetOptions(options map[string]string) {
	c.Options = options
}

// prepareRequest will setup a request by using method and endpoint
// it can be used to set a future API key
func (c *Client) prepareRequest(method, endpoint string) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", APIEndpoint, endpoint), nil)
	// add country if in options
	if country, ok := c.Options["country"]; ok {
		q := req.URL.Query()
		q.Add("country", country)
		req.URL.RawQuery = q.Encode()
	}
	return req, err
}

// Get is Client Get method
func (c *Client) Get(endpoint string) (*http.Response, error) {
	request, reqErr := c.prepareRequest("GET", endpoint)

	if reqErr != nil {
		return nil, reqErr
	}

	response, err := c.Client.Do(request)
	return response, err
}

// Info returns the info related to a manga
func (c *Client) Info(comicName string) (*MangaRockInfo, error) {
	endpoint := fmt.Sprintf("info?oid=%s", comicName)

	mangaInfo := new(MangaRockInfo)
	response, err := c.Get(endpoint)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	jsonErr := getJSON(response, &mangaInfo)

	return mangaInfo, jsonErr
}

// Pages returns the pages related to a manga
func (c *Client) Pages(comicName string) (*MangaRockPages, error) {
	endpoint := fmt.Sprintf("pages?oid=%s", comicName)

	mangaPages := new(MangaRockPages)
	response, err := c.Get(endpoint)

	if err != nil {
		return mangaPages, err
	}

	defer response.Body.Close()
	jsonErr := getJSON(response, &mangaPages)

	return mangaPages, jsonErr
}

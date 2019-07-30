package imgur

import (
	"fmt"
	qs "github.com/google/go-querystring/query"
	"github.com/tbruyelle/apiclient"
	"net/http"
)

const (
	defaultBaseURL = "https://api.imgur.com/3/"
)

type Client struct {
	api *apiclient.API
}

type Response struct {
	Success bool `json:"success"`
	Status  int  `json:"status"`
}

type SearchResponse struct {
	Response
	Data []Image `json:"data"`
}

type Image struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DateTime    int    `json:"datetime"`
	Type        string `json:"type"`
	Animated    bool   `json:"animated"`
	Nsfw        bool   `json:"nsfw"`
	Link        string `json:"link"`
}

func NewClient(clientId string) *Client {
	c := &Client{
		api: apiclient.New(defaultBaseURL),
	}
	c.api.Headers["Authorization"] = "Client-Id " + clientId
	return c
}

// SearchOptions specifies the parameters to the Search method.
type SearchOptions struct {
	Query   string `url:"q"`
	All     string `url:"q_all"`
	Any     string `url:"q_any"`
	Exactly string `url:"q_exactly"`
	Not     string `url:"q_not"`
	Type    string `url:"q_type"`
	SizePx  string `url:"q_size_px"`
}

func (c *Client) Search(opt SearchOptions) (*SearchResponse, *http.Response, error) {
	params, err := qs.Values(opt)
	if err != nil {
		return nil, nil, err
	}
	uri := fmt.Sprintf("gallery/search?%s", params.Encode())
	req, err := c.api.NewRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	searchResp := new(SearchResponse)
	resp, err := c.api.Do(req, searchResp)
	if err != nil {
		return nil, resp, err
	}
	return searchResp, resp, nil
}

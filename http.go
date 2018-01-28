package x

import (
	"net/url"
	"io/ioutil"
	"net/http"
)

type HttpClient http.Client

func (c *HttpClient) DoGet(url string) (*Response, error) {
	resp, err := (*http.Client)(c).Get(url)
	r := Response(*resp)
	return &r, err
}

type Response http.Response

func (r *Response) ReadBytes() []byte {
	defer r.Body.Close()
	bytes, _ := ioutil.ReadAll(r.Body)
	return bytes
}

func BuildUrl(rawUrl string, query map[string]string) *url.URL {
	u, _ := url.Parse(rawUrl)
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u
}

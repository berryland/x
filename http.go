package x

import (
	"net/url"
	"io/ioutil"
	"net/http"
	"strconv"
)

type HttpClient struct {
	*http.Client
}

type Query map[string]interface{}

func (q Query) Encode() Query {
	for k, v := range q {
		var s string
		switch x := v.(type) {
		case uint:
			s = strconv.FormatUint(uint64(x), 10)
		case uint8:
			s = strconv.FormatUint(uint64(x), 10)
		case uint16:
			s = strconv.FormatUint(uint64(x), 10)
		case uint32:
			s = strconv.FormatUint(uint64(x), 10)
		case uint64:
			s = strconv.FormatUint(uint64(x), 10)
		case int:
			s = strconv.FormatInt(int64(x), 10)
		case int8:
			s = strconv.FormatInt(int64(x), 10)
		case int16:
			s = strconv.FormatInt(int64(x), 10)
		case int32:
			s = strconv.FormatInt(int64(x), 10)
		case int64:
			s = strconv.FormatInt(int64(x), 10)
		case float32:
			s = strconv.FormatFloat(float64(x), 'f', -1, 32)
		case float64:
			s = strconv.FormatFloat(float64(x), 'f', -1, 64)
		case string:
			s = v.(string)
		}
		q[k] = s
	}
	return q
}

func (c *HttpClient) DoGet(url string, query Query) (*Response, error) {
	resp, err := c.Get(BuildUrl(url, query).String())
	var r Response
	if err == nil {
		r = Response(*resp)
		return &r, nil
	}

	return nil, err
}

type Response http.Response

func (r Response) ReadBytes() []byte {
	defer r.Body.Close()
	bytes, _ := ioutil.ReadAll(r.Body)
	return bytes
}

func BuildUrl(rawUrl string, query Query) *url.URL {
	u, _ := url.Parse(rawUrl)
	query.Encode()
	q := u.Query()
	for k, v := range query {
		q.Set(k, v.(string))
	}
	u.RawQuery = q.Encode()
	return u
}

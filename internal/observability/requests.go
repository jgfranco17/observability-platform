package observability

import (
	"net/http"
	"net/url"
	"time"
)

type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type jar struct {
	cookies map[string][]*http.Cookie
}

func (j *jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if j.cookies == nil {
		j.cookies = make(map[string][]*http.Cookie)
	}
	j.cookies[u.Host] = cookies
}

func (j *jar) Cookies(u *url.URL) []*http.Cookie {
	var cookieList []*http.Cookie
	for _, c := range j.cookies[u.Host] {
		if c.Expires.After(time.Now()) {
			cookieList = append(cookieList, c)
		}
	}
	return cookieList
}

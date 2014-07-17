package webqq

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
	. "webqq/tools"
)

func newClient(t time.Duration) http.Client {
	jar, err := cookiejar.New(nil)
	ErrHandle(err, `x`, `obtain_cookiejar`)

	return http.Client{
		&http.Transport{
			Dial: func(network string, address string) (net.Conn, error) {
				return net.DialTimeout(network, address, t*time.Millisecond)
			},
		},
		nil,
		jar,
		t * time.Millisecond,
	}
}

type Client struct {
	id         string
	password   string
	vfwebqq    string
	psessionid string
	clientid   string
	ptwebqq    string
	client     http.Client
}

func init() {
	fmt.Printf(``)
}

func New(id, password string, capacity int, timeout int) *Client {
	return &Client{
		id:       id,
		password: password,
		clientid: strconv.Itoa(rand.Intn(90000000) + 10000000),
		client:   newClient(time.Duration(timeout)),
	}
}

func (qq *Client) Get(url string) (re *http.Response, err error) {
	return qq.get(url)
}

func (qq *Client) get(u string) (re *http.Response, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	lg.Trace("\nGET: URL: %v", u)
	req, err := http.NewRequest("GET", u, nil)
	ErrHandle(err, `p`)

	req.Header.Add(`referer`, `http://d.web2.qq.com/proxy.html?v=20110331002&callback=2&id=3`)

	re, err = qq.client.Do(req)
	ErrHandle(err, `p`)

	if qq.ptwebqq == `` {
		for _, v := range re.Cookies() {
			if v.Name == `ptwebqq` {
				qq.ptwebqq = v.Value
			}
		}
	}

	qq.client.Jar.SetCookies(req.URL, re.Cookies())
	return
}

func (qq *Client) postForm(u string, data url.Values) (re *http.Response, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	lg.Trace("\nPOST: URL: %v\nDATA: %v", u, data.Encode())
	req, err := http.NewRequest("POST", u, strings.NewReader(data.Encode()))
	ErrHandle(err, `p`)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if qq.client.Jar != nil {
		for _, cookie := range qq.client.Jar.Cookies(req.URL) {
			req.AddCookie(cookie)
		}
	}

	req.Header.Add(`referer`, `http://d.web2.qq.com/proxy.html?v=20110331002&callback=2&id=3`)

	re, err = qq.client.Do(req)
	ErrHandle(err, `p`)

	qq.client.Jar.SetCookies(req.URL, re.Cookies())
	return
}

func (*Client) timeStamp() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())[:13]
}

func (this *Client) SetPtWebqq(p string) {
	this.ptwebqq = p
}

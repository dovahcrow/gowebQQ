package webqq

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
	. "webqq/tools"
)

var lg *logs.BeeLogger

type Client struct {
	id            string
	password      string
	vfwebqq       string
	psessionid    string
	clientid      string
	ptwebqq       string
	verifysession string
	client        http.Client
	msgid         int64
	*logs.BeeLogger
}

func init() {
	fmt.Printf(``)
}

func newClient(t time.Duration) http.Client {
	jar, err := cookiejar.New(nil)
	ErrHandle(err, `x`, `obtain_cookiejar`)

	return http.Client{
		nil,
		nil,
		jar,
		t * time.Millisecond,
	}
}

func New(id, password string, timeout int, logchannellen int64) *Client {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	c := &Client{
		id:       id,
		password: password,
		clientid: strconv.Itoa(rd.Intn(90000000) + 10000000),
		client:   newClient(time.Duration(timeout)),
		msgid:    (rd.Int63n(9000) + 1000) * 10000,
	}
	c.BeeLogger = logs.NewLogger(logchannellen)
	lg = c.BeeLogger
	return c
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
				continue
			}
			if v.Name == "verifysession" {
				qq.verifysession = v.Value
				continue
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

package qqclient

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"github.com/doomsplayer/Xgo-webqq/cl"
	"strconv"
	"strings"
	"sync"
	"time"
	. "github.com/doomsplayer/Xgo-webqq/tools"
)

type Client struct {
	id          string
	password    string
	vfwebqq     string
	psessionid  string
	clientid    string
	ptwebqq     string
	PollMutex   sync.Mutex
	MessagePool chan *PollMessage
	client      http.Client
}

func init() {
	fmt.Printf(``)
}

func New(id, password string) *Client {
	return &Client{MessagePool: make(chan *PollMessage, 100), id: id, password: password, clientid: strconv.Itoa(rand.Intn(90000000) + 10000000), client: cl.Client(20000)}
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

	req, err := http.NewRequest("GET", u, nil)
	ErrHandle(err, `p`)

	req.Header.Add(`referer`, `http://d.web2.qq.com/proxy.html?v=20110331002&callback=2&id=3`)
	re, err = qq.client.Do(req)
	if err != nil {
		panic(err)
	}
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

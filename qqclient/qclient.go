package qqclient

import (
	"io"
	"io/ioutil"
	"logind"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type QQClient struct {
	Id         string
	Password   string
	Vfwebqq    string
	Psessionid string
	Clientid   string
	Friend     struct {
		Friendlist map[string]*Friends
		Categories map[int]string
	}
	Group map[string]*Groups
	http.Client
}

func NewClient(id, password string) *QQClient {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := &QQClient{Id: id, Password: password, Vfwebqq: ``, Psessionid: ``, Clientid: strconv.Itoa(seed.Intn(90000000) + 10000000), http.Client: logind.Client()}
	p.Friend.Categories = make(map[int]string)
	p.Friend.Friendlist = make(map[string]*Friends)
	p.Group = make(map[string]*Groups)
	return p
}

func (qq *QQClient) pForm(ur string, data url.Values) (*http.Response, error) {
	req, _ := http.NewRequest("POST", ur, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if qq.Jar != nil {
		for _, cookie := range qq.Jar.Cookies(req.URL) {
			req.AddCookie(cookie)
		}
	}
	req.Header.Add(`referer`, `http://d.web2.qq.com/proxy.html?v=20110331002&callback=2`)
	re, err := qq.Do(req)
	qq.Jar.SetCookies(nil, re.Cookies())
	return re, err
}

func ehandle(err error) {
	if err != nil {
		panic(err)
	}
}

func sRead(re io.Reader) string {
	bBody, _ := ioutil.ReadAll(re)
	return string(bBody)
}

func FindCookies(cookies []*http.Cookie, name string) *http.Cookie {
	for _, n := range cookies {
		if n.Name == name {
			return n
		}
	}
	return nil
}
func (qq *QQClient) UintoNick(uin string) (string, error) {
	var rt error = nil
	defer func() {
		if err := recover(); err != nil {
			rt = err.(error)
		}
	}()
	return qq.Friend.Friendlist[uin].Nick, rt
}
func (qq *QQClient) UintoMarkname(uin string) (string, error) {
	var rt error = nil
	defer func() {
		if err := recover(); err != nil {
			rt = err.(error)
		}
	}()
	return qq.Friend.Friendlist[uin].Markname, rt
}
func (qq *QQClient) NicktoUin(nick string) string {
	for i, f := range qq.Friend.Friendlist {
		if f.Nick == nick {
			return i
		}
	}
	return ``
}

func (qq *QQClient) MarknametoUin(mkname string) string {
	for i, f := range qq.Friend.Friendlist {
		if f.Markname == mkname {
			return i
		}
	}
	return ``
}

package qqclient

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type retu struct {
	Retcode int `json:"retcode"`
	Result  struct {
		Uin        int    `json:"uin"`
		Clip       int    `json:"clip"`
		Index      int    `json:"index"`
		Port       int    `json:"port"`
		Status     string `json:"status"`
		Vfwebqq    string `json:"vfweqq"`
		Psessionid string `json:"psessionid"`
		F          int    `json:"f"`
		User_state int    `json:"user_state"`
	}
}

var sep = func(a rune) bool {
	if a == rune('\\') {
		return true
	}
	return false
}

func (qq *QQClient) Login() (rerr error) {
	vcs := ``
	vcl := ``
	defer func() {
		if err := recover(); err != nil {
			rerr = err.(error)
		}
	}()

	re, ok := qq.Get(`http://check.ptlogin2.qq.com/check?uin=` + qq.Id + `&appid=1003903&r=0.2667082343145449`)
	ehandle(ok)
	sBody := sRead(re.Body)
	if sBody[14] == uint8('0') {
		fmt.Print("不需要验证码\n")
		p := regexp.MustCompile(`ptui_checkVC\('0','(.*?)','(.*?)'\);`).FindStringSubmatch(sBody)
		sV := make([]byte, 0)
		for _, tm := range strings.FieldsFunc(p[2], sep) {
			t, _ := hex.DecodeString(tm[1:3])
			sV = append(sV, t...)
		}

		h1 := md5.New()
		h2 := md5.New()
		h3 := md5.New()
		fmt.Fprint(h1, qq.Password)
		fmt.Fprintf(h2, "%s", h1.Sum(nil))
		fmt.Fprintf(h2, "%s", sV)
		fmt.Fprintf(h3, "%s", fmt.Sprintf("%X", (h2.Sum(nil))))
		fmt.Fprintf(h3, "%s", strings.ToUpper(p[1]))
		vcl = fmt.Sprintf("%X", h3.Sum(nil))
		vcs = strings.ToUpper(p[1])
	} else if sBody[14] == uint8('1') {
		fmt.Print("需要验证码\n")
		p := regexp.MustCompile(`ptui_checkVC\('1','(.*?)','(.*?)'\);`).FindStringSubmatch(sBody)
		qq.picHandle(qq.Id, p[1])
		sChk := ``
		fmt.Printf("请输入验证码:  ")
		fmt.Scanf("%s", &sChk)
		sV := make([]byte, 0)
		for _, tm := range strings.FieldsFunc(p[2], sep) {
			t, _ := hex.DecodeString(tm[1:3])
			sV = append(sV, t...)
		}

		h1 := md5.New()
		h2 := md5.New()
		h3 := md5.New()
		fmt.Fprint(h1, qq.Password)
		fmt.Fprintf(h2, "%s", h1.Sum(nil))
		fmt.Fprintf(h2, "%s", sV)
		fmt.Fprintf(h3, "%s", fmt.Sprintf("%X", (h2.Sum(nil))))
		fmt.Fprintf(h3, "%s", strings.ToUpper(sChk))
		vcl = fmt.Sprintf("%X", h3.Sum(nil))
		vcs = strings.ToUpper(sChk)
	} else {
		panic(errors.New("登陆错误"))
	}

	re, ok = qq.Get(`http://ptlogin2.qq.com/login?u=` + qq.Id + `&p=` + vcl + `&verifycode=` + vcs + `&webqq_type=10&remember_uin=1&login2qq=1&aid=1003903&u1=http%3A%2F%2Fweb.qq.com%2Floginproxy.html%3Flogin2qq%3D1%26webqq_type%3D10&h=1&ptredirect=0&ptlang=2052&from_ui=1&pttype=1&dumy=&fp=loginerroralert&action=3-27-27672&mibao_css=m_webqq&t=1&g=1&js_type=0&js_ver=10020&login_sig=8t2fn380ZiJfY3qV48Ast8PCWNXLfqZzJ2r8W8nVr5d*gAaWdNbzYm2iIx2trLLo`)
	ehandle(ok)

	if sRead(re.Body)[8] == '0' {
		fmt.Println("第一次握手成功")
	} else {
		panic(errors.New(`第一次握手失败（密码错误？）`))
	}

	v := url.Values{}
	v.Set(`clientid`, qq.Clientid)
	v.Set(`psessionid`, `null`)
	c, ok := json.Marshal(
		map[string]interface{}{
			`status`:     `online`,
			`ptwebqq`:    FindCookies(qq.Jar.Cookies(nil), `ptwebqq`).Value,
			`passwd_sig`: ``,
			`clientid`:   qq.Clientid,
			`psessionid`: nil})
	v.Set(`r`, string(c))
	ehandle(ok)
	re, ok = qq.pForm(`http://d.web2.qq.com/channel/login2`, v)
	if ok != nil {
		panic(errors.New("第二次握手失败（网络错误？）"))
	}
	ehandle(ok)
	ret := make(map[string]interface{})
	json.Unmarshal([]byte(sRead(re.Body)), &ret)
	if i := ret[`retcode`].(float64); i == float64(0) {
		fmt.Println("第二次握手成功")
	} else {
		panic(fmt.Errorf("第二次握手失败,错误码：%v", i))
	}
	qq.Vfwebqq = ret[`result`].(map[string]interface{})[`vfwebqq`].(string)
	qq.Psessionid = ret[`result`].(map[string]interface{})[`psessionid`].(string)
	fmt.Println("口令获得，登陆成功")
	return

}

func (qq *QQClient) picHandle(id, vc string) {
	re, ok := qq.Get(`http://captcha.qq.com/getimage?aid=1003903&r=0.85623475915069394&uin=` + id + `&vc_type=` + strings.ToUpper(vc))
	ehandle(ok)
	toFile(re.Body, `pic.png`)
}

func toFile(re io.Reader, name string) {
	f, err := os.Create(name)
	ehandle(err)
	s, err := ioutil.ReadAll(re)
	ehandle(err)
	_, err = f.Write(s)
	ehandle(err)
	f.Close()
	t := ``
errh:
	fmt.Printf("请输入你的图片浏览工具名\n")
	fmt.Scanf("%s", &t)
	err = exec.Command(t, `pic.png`).Run()
	if err != nil {
		fmt.Println(err)
		goto errh
	}
}

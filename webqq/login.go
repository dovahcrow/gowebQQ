package webqq

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	//	"os/exec"
	"regexp"
	"strings"
	. "webqq/tools"
	"webqq/tools/simplejson"
)

type LoginRet struct {
	needPic bool
	PicAddr string
	p1      string
	p2      []byte
	pass    string
}

func (s *LoginRet) SetVFCode(i string) error {
	if s.needPic {
		s.p1 = i
		return nil
	} else {
		return fmt.Errorf(`No Need For Setting Verify Code!`)
	}
}

func (s *LoginRet) gen() (vcl string, vcs string) {
	lg.Trace("generating vcl and vcs")

	sum1 := md5.Sum([]byte(s.pass))
	sum2 := md5.Sum([]byte(append(sum1[:], s.p2...)))
	sum3 := md5.Sum([]byte(fmt.Sprintf("%X", sum2) + strings.ToUpper(s.p1)))

	vcl = fmt.Sprintf("%X", sum3)
	vcs = strings.ToUpper(s.p1)
	return
}

func (s *LoginRet) NeedPic() bool {
	return s.needPic
}

func (qq *Client) LoginStep1() (ret LoginRet, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	re, err := qq.get(`https://ssl.ptlogin2.qq.com/check?uin=` + qq.id + `&appid=1003903&js_ver=10040&js_type=0&login_sig=glCTiV1*UGC58vTwRS3f-xyFDmTfq45dfLQxy2IMjw8BGt1UUldhM9fq2EXdSamn&u1=http%3A%2F%2Fweb2.qq.com%2Floginproxy.html&r=0.` + fmt.Sprint(rand.Int63n(10000000000000000)))
	ErrHandle(err, `p`)
	defer re.Body.Close()

	sBody := ReadString(re.Body)

	lg.Debug("ptui_checkVC is %v", sBody)

	p := regexp.
		MustCompile(`^ptui_checkVC\('(\d)','(.*?)','(.*?)', '(.*?)'\);$`).
		FindStringSubmatch(sBody)

	switch p[1] {
	case `0`:
		{
			lg.Info("No need for captcha")
			ret.needPic = false
			ret.p1 = p[2]
		}
	case `1`:
		{
			lg.Info("Need captcha")
			ret.needPic = true
			ret.PicAddr = `https://ssl.captcha.qq.com/getimage?aid=1003903&r=0.10382554663612781&uin=` + qq.id
		}
	default:
		{
			err = fmt.Errorf(`login error`)
			return
		}
	}

	cs := strings.Split(p[3], `\x`)
	for _, v := range cs {
		token, _ := hex.DecodeString(v)
		ret.p2 = append(ret.p2, token...)
	}

	ret.pass = qq.password
	return
}

func (qq *Client) LoginStep2(ret LoginRet) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	if ret.p1 == `` {
		lg.Warn("p1 not set")
		return fmt.Errorf(`Please Set Verify Code First`)
	}

	vcl, vcs := ret.gen()

	re, err := qq.get(`https://ssl.ptlogin2.qq.com/login?u=` + qq.id + `&p=` + vcl + `&verifycode=` + vcs + `&webqq_type=10&remember_uin=1&login2qq=1&aid=1003903&u1=http%3A%2F%2Fweb2.qq.com%2Floginproxy.html%3Flogin2qq%3D1%26webqq_type%3D10&h=1&ptredirect=0&ptlang=2052&daid=164&from_ui=1&pttype=1&dumy=&fp=loginerroralert&action=3-29-12615&mibao_css=m_webqq&t=1&g=1&js_type=0&js_ver=10040&login_sig=glCTiV1*UGC58vTwRS3f-xyFDmTfq45dfLQxy2IMjw8BGt1UUldhM9fq2EXdSamn`)
	if err != nil {
		panic(fmt.Errorf("第一次握手失败,网络错误: %v", err))
	}
	defer re.Body.Close()

	sBody := ReadString(re.Body)
	lg.Debug("ptuiCB is %v", sBody)

	reg := regexp.MustCompile(`ptuiCB\('0','0','(.*)','0','登录成功！', '.*'\);`)

	if !reg.MatchString(sBody) {

		panic(errors.New(`第一次握手失败（密码错误？）`))
	}

	ssl := reg.FindStringSubmatch(sBody)
	re, err = qq.get(ssl[1])
	ErrHandle(err, `x`)
	defer re.Body.Close()

	v := url.Values{}
	v.Set(`clientid`, qq.clientid)
	v.Set(`psessionid`, `null`)

	c, err := json.Marshal(
		map[string]interface{}{
			`status`:     `online`,
			`ptwebqq`:    qq.ptwebqq,
			`passwd_sig`: ``,
			`clientid`:   qq.clientid,
			`psessionid`: nil})
	v.Set(`r`, string(c))

	re, err = qq.postForm(`http://d.web2.qq.com/channel/login2`, v)

	if err != nil {
		panic(fmt.Errorf("第二次握手失败:%v", err))
	}
	defer re.Body.Close()

	retb := ReadByte(re.Body)
	lg.Debug("online info is %s", retb)

	js, err := simplejson.NewJson(retb)
	if err != nil {
		panic(err)
	}

	if i := js.Get(`retcode`).MustFloat64(); i != float64(0) {
		panic(fmt.Errorf("第二次握手失败,错误码:%v", i))
	}

	qq.vfwebqq = js.Get(`result`).Get(`vfwebqq`).MustString()
	qq.psessionid = js.Get(`result`).Get(`psessionid`).MustString()
	return nil
}

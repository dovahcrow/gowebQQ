package qqclient

import (
	"encoding/json"
	"fmt"
	//"net/http"
	"net/url"
	"runtime"
	"strconv"
	"time"
	. "tools"
	"tools/simplejson"
)

func init() {
	fmt.Printf("")
}

var ELoginAgain = fmt.Errorf(`login again`)
var ENOMSG = fmt.Errorf(`No msg for 10 minutes`)

type PollMessage struct {
	Type     string
	Body     []string
	From_uin string
	T        time.Time
}

func (qq *Client) pollSafe() {
	f := func() chan struct{} {
		c := make(chan struct{})
		go func() {
			ret, err := qq.Poll()
			ErrHandle(err, `n`, `pollmsg`)
			for _, v := range ret {
				qq.MessagePool <- v
			}
			c <- struct{}{}
		}()
		return c
	}
	for {
		select {
		case <-f():
			{
			}
		case <-time.After(60 * time.Second):
			{
			}
		}
	}
}

func (qq *Client) Poll() (ret []*PollMessage, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	for len(ret) == 0 && err == nil {

		var r []byte
		r, err = qq.pollraw()
		ErrHandle(err, `p`)
		var js *simplejson.Json
		js, err = simplejson.NewJson(r)

		if err != nil {
			panic(fmt.Errorf("parse json error,%v", err))
		}

		var retcode int
		retcode, err = js.Get(`retcode`).Int()

		if err != nil {
			panic(fmt.Errorf("parse json error,%v", err))
		}
		fmt.Println(string(r))
		switch retcode {
		case 0:
			{
				result := js.Get(`result`)

				for i := 0; i < len(result.MustArray()); i++ {

					resulti := result.GetIndex(i)

					var poll_type string

					poll_type, err = resulti.Get(`poll_type`).String()

					if err != nil {
						panic(fmt.Errorf("parse poll_type error,%v", err))
					}

					switch poll_type {
					case `message`:
						{
							fuin := strconv.FormatInt(int64(resulti.Get(`value`).Get(`from_uin`).MustFloat64()), 10)
							msg := []string{}

							for i := 1; i < len(resulti.Get(`value`).Get(`content`).MustArray()); i++ {
								content, err := resulti.Get(`value`).Get(`content`).GetIndex(i).String()
								if err != nil || content == ` ` {
									continue
								}
								msg = append(msg, content)
							}
							t := int64(resulti.Get(`value`).Get(`time`).MustFloat64())
							ret = append(ret, &PollMessage{Type: `buddyMsg`, Body: msg, From_uin: fuin, T: time.Unix(t, 0)})
						}
					case `system_message`:
						{
							switch resulti.Get(`value`).Get(`type`).MustString() {
							case `added_buddy_sig`:
								{
									t := int64(resulti.Get(`value`).Get(`time`).MustFloat64())
									fuin := strconv.FormatInt(int64(resulti.Get(`value`).Get(`from_uin`).MustFloat64()), 10)
									ret = append(ret, &PollMessage{Type: `addBuddy`, From_uin: fuin, T: time.Unix(t, 0)})
								}
							}
						}

					case `kick_message`:
						{ //您的帐号在另一地点登录，您已被迫下线。
							ret = append(ret, &PollMessage{Type: `kicked`, From_uin: `10000`, T: time.Now()})
							qq.PollMutex.Lock()
						}
					case `group_message`:
						{
							fuin := strconv.FormatInt(int64(resulti.Get(`value`).Get("send_uin").MustFloat64()), 10) + `@` + strconv.FormatInt(int64(resulti.Get(`value`).Get(`from_uin`).MustFloat64()), 10)
							msg := []string{}

							for i := 1; i < len(resulti.Get(`value`).Get(`content`).MustArray()); i++ {
								content, err := resulti.Get(`value`).Get(`content`).GetIndex(i).String()
								if err != nil || content == ` ` {
									continue
								}
								msg = append(msg, content)
							}
							t := int64(resulti.Get(`value`).Get(`time`).MustFloat64())
							ret = append(ret, &PollMessage{Type: `groupMsg`, Body: msg, From_uin: fuin, T: time.Unix(t, 0)})
						}
					case `input_status`:
						{
						}
					case `buddies_status_change`:
						{

						}
					case `tips`:
						{
							//news
						}
					case `input_notify`:
						{
						}
					case `ok`:
						{

						}
					default:
						{
							//fmt.Println(poll_type, `potp`)
						}
					}
				}
			}
		case 102, 116:
			{
			}
		case 103, 121, 100006:
			{ //断线
				ret = append(ret, &PollMessage{Type: `offline`, T: time.Now()})
				qq.PollMutex.Lock()
			}
		default:
			{
				err = fmt.Errorf("错误！：%v", string(r))
			}
		}
	}
	return
}

func (qq *Client) pollraw() (retu []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
		qq.PollMutex.Unlock()
	}()
	runtime.Gosched()
	qq.PollMutex.Lock()

	c, _ := json.Marshal(map[string]interface{}{
		`clientid`:   qq.clientid,
		`psessionid`: qq.psessionid,
		`key`:        0,
		`ids`:        [0]interface{}{}})
	v := url.Values{}
	v.Set(`clientid`, qq.clientid)
	v.Set(`psessionid`, qq.psessionid)
	v.Set(`r`, string(c))
	re, err := qq.postForm(`http://d.web2.qq.com/channel/poll2`, v)
	ErrHandle(err, `p`)
	defer re.Body.Close()
	retu = ReadByte(re.Body)
	return

}

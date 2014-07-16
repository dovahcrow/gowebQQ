package webqq

import (
	"encoding/json"
	"fmt"
	"time"
	//"net/http"
	"net/url"
	"strconv"
	. "webqq/tools"
	"webqq/tools/simplejson"
)

func init() {
	fmt.Printf("")
}

var ELoginAgain = fmt.Errorf(`login again`)
var ENOMSG = fmt.Errorf(`No msg for 10 minutes`)

type PollType string

const (
	PT_Message             = PollType("message")
	PT_SystemMessage       = PollType("system_message")
	PT_KickMessage         = PollType("kick_message")
	PT_GroupMessage        = PollType("group_message")
	PT_InputStatus         = PollType("input_status")
	PT_BuddiesStatusChange = PollType("buddies_status_change")
	PT_Tips                = PollType("tips")
	PT_OK                  = PollType("ok")
	PT_InputNotify         = PollType("input_notify")
)

type PollMessage struct {
	retCode int
	//if retcode == 0
	pollType PollType
	value    []byte

	//if retcode != 0
	errMsg string

	//Public
	t time.Time
}

func (this *PollMessage) RetCode() int {
	return this.retCode
}
func (this *PollMessage) PollType() (PollType, error) {
	if this.retCode != 0 {
		return ``, fmt.Errorf("ret code is not 0")
	}
	return this.pollType, nil
}

type BuddyMessage struct {
	PollType    PollType
	MessageId   int64
	FromUin     string
	ToUin       string
	MessageId2  int64
	MessageType int64
	ReplyIp     int64
	Time        time.Time
	Font        struct {
		Size  int
		Color string
		Style [3]int
		Name  string
	}
	Content []string
}

func (this *PollMessage) IsBuddyMessage() (ret *BuddyMessage, is bool) {
	if this.retCode != 0 || this.pollType != PT_Message {
		return nil, false
	}

	js, err := simplejson.NewJson(this.value)
	if err != nil {
		return nil, false
	}
	fromUin := strconv.FormatInt(int64(js.Get(`from_uin`).MustFloat64()), 10)
	toUin := strconv.FormatInt(int64(js.Get(`to_uin`).MustFloat64()), 10)
	msgId := int64(js.Get("msg_id").MustFloat64())
	msgId2 := int64(js.Get("msg_id2").MustFloat64())
	msgType := int64(js.Get("msg_type").MustFloat64())
	replyIp := int64(js.Get("reply_ip").MustFloat64())
	t := time.Unix(int64(js.Get("time").MustFloat64()), 0)
	contentjs := js.Get("content")
	fontjs := contentjs.GetIndex(0).GetIndex(1)
	fontSize := int(fontjs.Get("size").MustFloat64())
	color := fontjs.Get("color").MustString("000000")
	fontName := fontjs.Get("name").MustString("")
	content := []string{}
	for i := 1; i < len(contentjs.MustArray([]interface{}{})); i++ {
		content = append(content, fmt.Sprint(contentjs.GetIndex(i)))

	}

	ret = new(BuddyMessage)
	ret.FromUin = fromUin
	ret.ToUin = toUin
	ret.MessageId = msgId
	ret.MessageId2 = msgId2
	ret.MessageType = msgType
	ret.ReplyIp = replyIp
	ret.Time = t
	ret.PollType, _ = this.PollType()
	ret.Font.Color = color
	ret.Font.Name = fontName
	ret.Font.Size = fontSize
	ret.Font.Style = [3]int{0, 0, 0}
	ret.Content = content
	return ret, true

}

type BuddyStatusChange struct {
	Uin        string
	Status     string
	ClientType int
}

func (this *PollMessage) IsBuddyStatusChange() (ret *BuddyStatusChange, is bool) {

	if this.retCode != 0 || this.pollType != PT_BuddiesStatusChange {
		return nil, false
	}
	js, err := simplejson.NewJson([]byte(this.value))
	if err != nil {
		return nil, false
	}
	uin := strconv.FormatInt(int64(js.Get(`uin`).MustFloat64()), 10)
	status := js.Get("status").MustString("offline")
	clientType := int(js.Get("client_type").MustFloat64())

	ret = new(BuddyStatusChange)
	ret.Uin = uin
	ret.Status = status
	ret.ClientType = clientType
	return ret, true
}

func (this *PollMessage) IsNothing() bool {
	if this.retCode == 102 {
		return true
	}
	return false
}
func (this *PollMessage) IsKick() bool {
	if this.retCode != 0 || this.pollType != PT_KickMessage {
		return false
	}
	return true
}
func ParseRawPoll(retu []byte) (ret []*PollMessage, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}

	}()
	lg.Debug("parse raw poll: msg is %s", retu)
	js, err := simplejson.NewJson(retu)
	if err != nil {
		panic(fmt.Errorf("parse poll error,%v", err))
	}

	retcode, err := js.Get(`retcode`).Int()

	if err != nil {
		panic(fmt.Errorf("parse poll error,%v", err))
	}

	lg.Info("ret code is %d", retcode)

	switch retcode {
	case 0:
		{

			result := js.Get(`result`)

			for i := 0; i < len(result.MustArray()); i++ {
				r := new(PollMessage)
				r.retCode = 0
				r.t = time.Now()

				resulti := result.GetIndex(i)

				poll_type := resulti.Get(`poll_type`).MustString("")

				r.value, _ = resulti.Get(`value`).MarshalJSON()

				switch PollType(poll_type) {
				case PT_Message:
					r.pollType = PT_Message

				case PT_SystemMessage:
					r.pollType = PT_SystemMessage

					// switch resulti.Get(`value`).Get(`type`).MustString() {
					// case `added_buddy_sig`:
					// 	{
					// 		t := int64(resulti.Get(`value`).Get(`time`).MustFloat64())
					// 		fuin := strconv.FormatInt(int64(resulti.Get(`value`).Get(`from_uin`).MustFloat64()), 10)
					// 		ret = append(ret, &PollMessage{Type: `addBuddy`, FromUin: fuin, T: time.Unix(t, 0)})
					// 	}
					// }

				case PT_KickMessage:
					r.pollType = PT_KickMessage //您的帐号在另一地点登录，您已被迫下线。

					// ret = append(ret, &PollMessage{Type: `kicked`, FromUin: `10000`, T: time.Now()})

				case PT_GroupMessage:
					r.pollType = PT_GroupMessage

					// fuin := strconv.FormatInt(int64(resulti.Get(`value`).Get("send_uin").MustFloat64()), 10) + `@` + strconv.FormatInt(int64(resulti.Get(`value`).Get(`from_uin`).MustFloat64()), 10)
					// msg := []string{}

					// for i := 1; i < len(resulti.Get(`value`).Get(`content`).MustArray()); i++ {
					// 	content, err := resulti.Get(`value`).Get(`content`).GetIndex(i).String()
					// 	if err != nil || content == ` ` {
					// 		continue
					// 	}
					// 	msg = append(msg, content)
					// }
					// t := int64(resulti.Get(`value`).Get(`time`).MustFloat64())
					// ret = append(ret, &PollMessage{Type: `groupMsg`, Body: msg, FromUin: fuin, T: time.Unix(t, 0)})

				case PT_InputStatus:
					r.pollType = PT_InputStatus

				case PT_BuddiesStatusChange:
					r.pollType = PT_BuddiesStatusChange
				case PT_Tips:
					r.pollType = PT_Tips

				case PT_InputNotify:

					r.pollType = PT_InputNotify

				case PT_OK:
					r.pollType = PT_OK

				default:
					lg.Critical("doesn't expected poll type", poll_type)

				}
				ret = append(ret, r)
			}
		}
	case 102:
		{
			r := new(PollMessage)
			r.retCode = 102
			r.t = time.Now()
			r.errMsg = js.Get("errmsg").MustString("")
			ret = append(ret, r)
		}
	case 116:
		{
			r := new(PollMessage)
			r.retCode = 102
			r.t = time.Now()
			r.errMsg = js.Get("errmsg").MustString("")
			ret = append(ret, r)
		}
	case 103:
		{
			r := new(PollMessage)
			r.retCode = 103
			r.t = time.Now()
			r.errMsg = js.Get("errmsg").MustString("")
			ret = append(ret, r)
		}
	case 121:
		{
			r := new(PollMessage)
			r.retCode = 121
			r.t = time.Now()
			r.errMsg = js.Get("errmsg").MustString("")
			ret = append(ret, r)
		}
	case 100006:
		{
			r := new(PollMessage)
			r.retCode = 1000056
			r.t = time.Now()
			r.errMsg = js.Get("errmsg").MustString("")
			ret = append(ret, r)
		}

	default:
		{
			err = fmt.Errorf("unknown ret code：%v", retcode)
		}
	}
	return
}

func (qq *Client) RawPoll() (retu []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}

	}()

	c, _ := json.Marshal(map[string]interface{}{
		`clientid`:   qq.clientid,
		`psessionid`: qq.psessionid,
		`key`:        0,
		`ids`:        [0]interface{}{}})

	v := url.Values{}
	v.Set(`clientid`, qq.clientid)
	v.Set(`psessionid`, qq.psessionid)
	v.Set(`r`, string(c))

	lg.Debug("poll req is %v", v.Encode())

	re, err := qq.postForm(`http://d.web2.qq.com/channel/poll2`, v)
	ErrHandle(err, `p`)
	defer re.Body.Close()
	retu = ReadByte(re.Body)
	lg.Debug("poll raw msg is %s", retu)
	return

}

// func (qq *Client) pollSafe() {
// 	f := func() chan struct{} {
// 		c := make(chan struct{})
// 		go func() {
// 			ret, err := qq.Poll()
// 			ErrHandle(err, `n`, `pollmsg`)
// 			for _, v := range ret {
// 				qq.MessagePool <- v
// 			}
// 			c <- struct{}{}
// 		}()
// 		return c
// 	}
// 	for {
// 		select {
// 		case <-f():
// 			{
// 			}
// 		case <-time.After(60 * time.Second):
// 			{
// 			}
// 		}
// 	}
// }

// func (qq *Client) Poll() (ret []*PollMessage, err error) {
// 	defer func() {
// 		if e := recover(); e != nil {
// 			err = e.(error)
// 		}
// 	}()
// 	for len(ret) == 0 && err == nil {

// 		var r []byte
// 		r, err = qq.pollraw()
// 		ErrHandle(err, `p`)
// 		var js *simplejson.Json
// 		js, err = simplejson.NewJson(r)

// 		if err != nil {
// 			panic(fmt.Errorf("parse json error,%v", err))
// 		}

// 		var retcode int
// 		retcode, err = js.Get(`retcode`).Int()

// 		if err != nil {
// 			panic(fmt.Errorf("parse json error,%v", err))
// 		}
// 		fmt.Println(string(r))
// 		switch retcode {
// 		case 0:
// 			{
// 				result := js.Get(`result`)

// 				for i := 0; i < len(result.MustArray()); i++ {

// 					resulti := result.GetIndex(i)

// 					var poll_type string

// 					poll_type, err = resulti.Get(`poll_type`).String()

// 					if err != nil {
// 						panic(fmt.Errorf("parse poll_type error,%v", err))
// 					}

// 					switch poll_type {
// 					case `message`:
// 						{
// 							fuin := strconv.FormatInt(int64(resulti.Get(`value`).Get(`from_uin`).MustFloat64()), 10)
// 							msg := []string{}

// 							for i := 1; i < len(resulti.Get(`value`).Get(`content`).MustArray()); i++ {
// 								content, err := resulti.Get(`value`).Get(`content`).GetIndex(i).String()
// 								if err != nil || content == ` ` {
// 									continue
// 								}
// 								msg = append(msg, content)
// 							}
// 							t := int64(resulti.Get(`value`).Get(`time`).MustFloat64())
// 							ret = append(ret, &PollMessage{Type: `buddyMsg`, Body: msg, FromUin: fuin, T: time.Unix(t, 0)})
// 						}
// 					case `system_message`:
// 						{
// 							switch resulti.Get(`value`).Get(`type`).MustString() {
// 							case `added_buddy_sig`:
// 								{
// 									t := int64(resulti.Get(`value`).Get(`time`).MustFloat64())
// 									fuin := strconv.FormatInt(int64(resulti.Get(`value`).Get(`from_uin`).MustFloat64()), 10)
// 									ret = append(ret, &PollMessage{Type: `addBuddy`, FromUin: fuin, T: time.Unix(t, 0)})
// 								}
// 							}
// 						}

// 					case `kick_message`:
// 						{ //您的帐号在另一地点登录，您已被迫下线。
// 							ret = append(ret, &PollMessage{Type: `kicked`, FromUin: `10000`, T: time.Now()})
// 							qq.PollMutex.Lock()
// 						}
// 					case `group_message`:
// 						{
// 							fuin := strconv.FormatInt(int64(resulti.Get(`value`).Get("send_uin").MustFloat64()), 10) + `@` + strconv.FormatInt(int64(resulti.Get(`value`).Get(`from_uin`).MustFloat64()), 10)
// 							msg := []string{}

// 							for i := 1; i < len(resulti.Get(`value`).Get(`content`).MustArray()); i++ {
// 								content, err := resulti.Get(`value`).Get(`content`).GetIndex(i).String()
// 								if err != nil || content == ` ` {
// 									continue
// 								}
// 								msg = append(msg, content)
// 							}
// 							t := int64(resulti.Get(`value`).Get(`time`).MustFloat64())
// 							ret = append(ret, &PollMessage{Type: `groupMsg`, Body: msg, FromUin: fuin, T: time.Unix(t, 0)})
// 						}
// 					case `input_status`:
// 						{
// 						}
// 					case `buddies_status_change`:
// 						{

// 						}
// 					case `tips`:
// 						{
// 							//news
// 						}
// 					case `input_notify`:
// 						{
// 						}
// 					case `ok`:
// 						{

// 						}
// 					default:
// 						{
// 							//fmt.Println(poll_type, `potp`)
// 						}
// 					}
// 				}
// 			}
// 		case 102, 116:
// 			{
// 			}
// 		case 103, 121, 100006:
// 			{ //断线
// 				ret = append(ret, &PollMessage{Type: `offline`, T: time.Now()})
// 				qq.PollMutex.Lock()
// 			}
// 		default:
// 			{
// 				err = fmt.Errorf("错误！：%v", string(r))
// 			}
// 		}
// 	}
// 	return
// }
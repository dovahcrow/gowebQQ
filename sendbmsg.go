package qqclient

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	. "github.com/doomsplayer/Xgo-webqq/tools"
	"github.com/doomsplayer/Xgo-webqq/tools/simplejson"
)

func (qq *Client) buddyMsgStructer(uiuin string, msg_id int, msg, fontname, fontsize, fontcolor string, fontstyle []int) url.Values {
	uin, _ := strconv.Atoi(uiuin)
	v := url.Values{}
	v.Set(`clientid`, qq.clientid)
	v.Set(`psessionid`, qq.psessionid)
	ms := [2]interface{}{
		msg,
		[2]interface{}{
			"font",
			map[string]interface{}{
				"name":  fontname,
				"size":  fontsize,
				"style": fontstyle,
				"color": fontcolor}}}
	byts, _ := json.Marshal(ms)
	m := map[string]interface{}{
		"to":         uin,
		"face":       0,
		"content":    string(byts),
		"msg_id":     msg_id,
		"clientid":   qq.clientid,
		"psessionid": qq.psessionid,
	}

	byts, _ = json.Marshal(m)
	v.Set(`r`, string(byts))
	return v
}

func (qq *Client) EasyBuddySend(uin string, msg_id int, msg string) (err error) {
	return qq.Sendbmsg(uin, msg_id, msg, `黑体`, `12`, `000000`, []int{0, 0, 0})
}

func (qq *Client) Sendbmsg(uin string, msg_id int, msg, fontname, fontsize, fontcolor string, fontstyle []int) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	v := qq.buddyMsgStructer(uin, msg_id, msg, fontname, fontsize, fontcolor, fontstyle)
	re, err := qq.postForm(`http://d.web2.qq.com/channel/send_buddy_msg2`, v)
	if err != nil {
		panic(err)
	}
	ret, err := simplejson.NewJson(ReadByte(re.Body))
	if err != nil {
		panic(err)
	}
	if i := ret.Get(`retcode`).MustInt(); i == 0 {
		return nil
	} else {
		panic(fmt.Errorf("发送个人消息:%v 失败，错误代码：%v", msg, i))
	}
	return
}

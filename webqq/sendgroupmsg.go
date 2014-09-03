package webqq

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	. "webqq/tools"
	"webqq/tools/simplejson"
)

func (qq *Client) groupMsgStructer(uiuin string, msg_id int64, msg, fontname, fontsize, fontcolor string, fontstyle [3]int) url.Values {
	uin, _ := strconv.Atoi(uiuin)
	v := url.Values{}
	v.Set(`clientid`, qq.clientid)
	v.Set(`psessionid`, qq.psessionid)

	ms := []interface{}{
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
		"group_uin":  uin,
		"content":    string(byts),
		"msg_id":     msg_id,
		"clientid":   qq.clientid,
		"psessionid": qq.psessionid,
	}

	byts, _ = json.Marshal(m)
	v.Set(`r`, string(byts))
	return v
}

func (qq *Client) SendGroupMsgEasy(uin string, msg_id int64, msg string) (err error) {
	return qq.SendGroupMsg(uin, msg_id, msg, `宋体`, `15`, `000000`, [3]int{0, 0, 0})
}

func (qq *Client) SendGroupMsg(uin string, msg_id int64, msg, fontname, fontsize, fontcolor string, fontstyle [3]int) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
		qq.msgid++
	}()
	fmt.Println(qq.msgid)
	v := qq.groupMsgStructer(uin, qq.msgid, msg, fontname, fontsize, fontcolor, fontstyle)
	re, err := qq.postForm(`http://d.web2.qq.com/channel/send_qun_msg2`, v)

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
		panic(fmt.Errorf("发送群消息:%v 失败，错误代码:%v", msg, i))
	}
	return
}

//r={"content":"[[\"font\",{\"name\":\"宋体\",\"size\":\"10\",\"style\":[0,0,0],\"color\":\"000000\"}]]"}

//r={"content":"[[\"font\",{\"color\":\"000000\",\"name\":\"宋体\",\"size\":\"15\",\"style\":[0,0,0]}]]","face":0}

// {
//     "clientid": "66979012",
//     "psessionid": "8368046764001d636f6e6e7365727665725f77656271714031302e3133332e34312e383400001a5e00000163026e04008a26454d6d0000000a40385775676e7446777a6d000000286798f9b9ed9385e4a86a8933ea02994fd55e386d7fa95b8190110c82f80416edbe540258ebd1980c",
//     "key": 0,
//     "ids": [
//         "58705"
//     ]
// },
// {
//     "group_uin": 1461871639,
//     "content": "[\"l\",[\"font\",{\"name\":\"宋体\",\"size\":\"10\",\"style\":[0,0,0],\"color\":\"000000\"}]]",
//     "msg_id": 89790003,
//     "clientid": "66979012",
//     "psessionid": "8368046764001d636f6e6e7365727665725f77656271714031302e3133332e34312e383400001a5e00000163026e04008a26454d6d0000000a40385775676e7446777a6d000000286798f9b9ed9385e4a86a8933ea02994fd55e386d7fa95b8190110c82f80416edbe540258ebd1980c"
// }

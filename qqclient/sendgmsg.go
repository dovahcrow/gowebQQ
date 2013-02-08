package qqclient

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

func init() {
	fmt.Printf("")
}

var Groupname = 1
var Groupid = 2

func (qq *QQClient) Sendgmsg(tp int, in, msg string) (rerr error) {
	defer func() {
		if e := recover(); e != nil {
			rerr = e.(error)
		}
	}()
	guin := ``
	switch tp {
	case 1:
		{
			found := false
			for i, b := range qq.Group {
				if b.Name == in {
					guin = i
					found = true
					break
				}
			}
			if found == false {
				panic(fmt.Errorf(`cannot find specific name`))
			}
		}
	case 2:
		guin = in
	default:
		panic(fmt.Errorf(`not a correct send type`))
	}
	if qq.Group[guin].Msg_id == 0 {
		qq.Group[guin].Msg_id = (rand.New(rand.NewSource(time.Now().UnixNano())).Intn(9000)+1000)*10000 + 1
	}
	v := url.Values{}
	v.Set(`clientid`, qq.Clientid)
	v.Set(`psessionid`, qq.Psessionid)

	//json.Marshal(map[string]interface{}{
	//	`group_uin`:         tid,
	//	`msg_id`:     23500401,
	//	`clientid`:   clientid,
	//	`psessionid`: qq.Psessionid,
	//	`content`: [2]interface{}{
	//		msg,
	//		[2]interface{}{
	//			"font",
	//			map[string]interface{}{
	//				"name":  "宋体",
	//				"size":  "10",
	//				"style": [3]int{0, 0, 0},
	//				"color": "993366"}}}})

	f := `{"group_uin":` + guin + `,"content":"[\"` + msg + `\",[\"font\",{\"name\":\"楷体_GB2312\",\"size\":\"10\",\"style\":[1,0,0],\"color\":\"808080\"}]]","msg_id":` + strconv.Itoa(qq.Group[guin].Msg_id) + `,"clientid":"` + qq.Clientid + `","psessionid":"` + qq.Psessionid + `"}`
	v.Set(`r`, f)
	re, err := qq.pForm(`http://d.web2.qq.com/channel/send_qun_msg2`, v)
	ehandle(err)
	qq.Group[guin].Msg_id++
	ret := make(map[string]interface{})
	q := sRead(re.Body)
	json.Unmarshal([]byte(q), &ret)
	if i := ret[`retcode`].(float64); i == float64(0) {
	} else {
		panic(fmt.Errorf("发送失败，错误代码：%v", i))
	}
	return
}

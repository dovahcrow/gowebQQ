package qqclient

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

//map[string]interface{}{
//	`to`:         tid,
//	`face`:       0,
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
//				"color": "993366"}}}}

func init() {
	fmt.Printf("")
}

type ftype int

var SNick ftype = 1
var SMarkName ftype = 2
var SUin ftype = 3

func (qq *QQClient) Sendbmsg(tp ftype, in, msg string) (rerr error) {
	defer func() {
		if e := recover(); e != nil {
			rerr = e.(error)
		}
	}()
	tid := ``
	switch tp {
	case 1:
		{
			found := false
			for i, b := range qq.Friend.Friendlist {
				if b.Nick == in {
					tid = i
					found = true
					break
				}
			}
			if found == false {
				panic(fmt.Errorf(`cannot find specific nick`))
			}
		}
	case 2:
		{
			found := false
			for i, b := range qq.Friend.Friendlist {
				if b.Markname == in {
					tid = i
					found = true
					break
				}
			}
			if found == false {
				panic(fmt.Errorf(`cannot find specific markname`))
			}

		}
	case 3:
		tid = in
	default:
		panic(fmt.Errorf(`not a correct send type`))
	}

	if qq.Friend.Friendlist[tid].Msg_id == 0 {
		qq.Friend.Friendlist[tid].Msg_id = (rand.New(rand.NewSource(time.Now().UnixNano())).Intn(9000)+1000)*10000 + 1
	}
	v := url.Values{}
	v.Set(`clientid`, qq.Clientid)
	v.Set(`psessionid`, qq.Psessionid)
	f := `{"to":` + tid + `,"face":0,"content":"[\"` + msg + `\",[\"font\",{\"name\":\"楷体_GB2312\",\"size\":\"10\",\"style\":[1,0,0],\"color\":\"808080\"}]]","msg_id":` + strconv.Itoa(qq.Friend.Friendlist[tid].Msg_id) + `,"clientid":"` + qq.Clientid + `","psessionid":"` + qq.Psessionid + `"}`
	v.Set(`r`, f)
	re, err := qq.pForm(`http://d.web2.qq.com/channel/send_buddy_msg2`, v)
	ehandle(err)
	qq.Friend.Friendlist[tid].Msg_id++
	ret := make(map[string]interface{})
	q := sRead(re.Body)
	json.Unmarshal([]byte(q), &ret)
	if i := ret[`retcode`].(float64); i == float64(0) {
	} else {
		panic(fmt.Errorf("发送失败，错误代码：%v", i))
	}
	return
}

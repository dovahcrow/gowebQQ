package qqclient

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func init() {
	fmt.Printf("")
}

type poll struct {
	Retcode int
	Result  []struct {
		Poll_type string
		Value     struct {
			Msg_id     int
			From_uin   int
			To_uin     int
			Msg_id2    int
			Msg_type   int
			Reply_ip   int
			Group_code int
			Send_uin   int
			Seq        int
			Time       int
			Info_seq   int
			Content    struct {
				Vd struct {
					Font string
					Vs   struct {
						size  int
						color string
						style []int
						name  string
					}
				}
				Msg string
			}
		}
	}
}

type Msg struct {
	Msgtype string
	Fromuin string
	MsgBody string
	Time    time.Time
	Guin    string
}

func (qq *QQClient) Pollmsg(ch chan Msg, pol chan int) {
	defer func() {
		if e := recover(); e != nil {

			fmt.Printf("%v\n", e.(error))
		}
		<-pol
	}()
	pol <- 1
	c, _ := json.Marshal(map[string]interface{}{
		`clientid`:   qq.Clientid,
		`psessionid`: qq.Psessionid,
		`key`:        0,
		`ids`:        [0]interface{}{}})
	v := url.Values{}
	v.Set(`clientid`, qq.Clientid)
	v.Set(`psessionid`, qq.Psessionid)
	v.Set(`r`, string(c))
	//for {
	re, err := qq.pForm(`http://d.web2.qq.com/channel/poll2`, v)
	ehandle(err)
	ret := make(map[string]interface{})
	json.Unmarshal([]byte(sRead(re.Body)), &ret)
	if ret[`retcode`].(float64) != float64(0) {
	}
	for _, k := range ret[`result`].([]interface{}) {
		switch k.(map[string]interface{})[`poll_type`].(string) {
		case `message`:
			{
				fuin := strconv.FormatFloat(k.(map[string]interface{})[`value`].(map[string]interface{})[`from_uin`].(float64), 'g', 20, 64)
				//fname, _ := qq.UintoNick(fuin)
				msg := ``
				switch k.(map[string]interface{})[`value`].(map[string]interface{})[`content`].([]interface{})[1].(type) {
				case []interface{}:
					msg = k.(map[string]interface{})[`value`].(map[string]interface{})[`content`].([]interface{})[0].(string)
				case string:
					msg = k.(map[string]interface{})[`value`].(map[string]interface{})[`content`].([]interface{})[1].(string)
				}
				msgtp := `bmsg`
				time := time.Unix(int64(k.(map[string]interface{})[`value`].(map[string]interface{})[`time`].(float64)), 0)
				ch <- Msg{msgtp, fuin, msg, time, ``}
			}

		case `group_message`:
			{
				msgtp := `gmsg`
				fuin := strconv.FormatFloat(k.(map[string]interface{})[`value`].(map[string]interface{})[`from_uin`].(float64), 'g', 20, 64)
				suin := strconv.FormatFloat(k.(map[string]interface{})[`value`].(map[string]interface{})["send_uin"].(float64), 'g', 20, 64)
				msg := ``
				switch k.(map[string]interface{})[`value`].(map[string]interface{})[`content`].([]interface{})[1].(type) {
				case []interface{}:
					msg = k.(map[string]interface{})[`value`].(map[string]interface{})[`content`].([]interface{})[0].(string)
				case string:
					msg = k.(map[string]interface{})[`value`].(map[string]interface{})[`content`].([]interface{})[1].(string)
				}
				time := time.Unix(int64(k.(map[string]interface{})[`value`].(map[string]interface{})[`time`].(float64)), 0)
				ch <- Msg{msgtp, suin, msg, time, fuin}
			}

		}
		//}
	}
}

//map[result:[map[poll_type:message value:map[from_uin:3.26513531e+09 content:[[font map[name:楷体_GB2312 color:808080 style:[1 0 0] size:12]] xz ] to_uin:1.296377482e+09 time:1.360055642e+09 reply_ip:1.76498346e+08 msg_id:62008 msg_id2:418681 msg_type:9]]] retcode:0]
//map[value:map[send_uin:1.975494613e+09 info_seq:1.08593811e+08 from_uin:5.41199791e+08 group_code:3.104313843e+09 msg_type:43 to_uin:1.73165159e+08 content:[[font map[color:000000 style:[1 0 0] name:新宋体 size:15]] 草木长的好好的，能不伤，就不伤 ] msg_id2:753421 seq:597706 msg_id:3795 time:1.360321115e+09 reply_ip:1.76498138e+08] poll_type:group_message]
//[map[poll_type:group_message value:map[group_code:2.662513848e+09 msg_type:43 content:[[font map[name:微软雅黑 style:[1 1 0] size:22 color:ff0000]] ? ] seq:3945 msg_id:10320 msg_id2:536385 time:1.360321128e+09 reply_ip:1.76722727e+08 send_uin:4.087066516e+09 to_uin:1.73165159e+08 info_seq:2.61490394e+08 from_uin:3.399522231e+09]]]
//[map[value:map[to_uin:1.73165159e+08 content:[[font map[color:800000 size:22 style:[1 1 0] name:新宋体]] [cface map[key:naRNJPH5kYZnevrt file_id:2.581137474e+09 server:123.138.154.215:8000 name:{5CC65E22-0D95-4D5A-254F-65B25FCCDD07}.jpg]]  ] msg_id2:891838 time:1.360321155e+09 send_uin:4.104993788e+09 from_uin:3.399522231e+09 msg_id:19921 seq:3946 group_code:2.662513848e+09 msg_type:43 reply_ip:1.76722709e+08 info_seq:2.61490394e+08] poll_type:group_message]]

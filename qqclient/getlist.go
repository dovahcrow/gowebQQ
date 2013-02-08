package qqclient

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

func init() {
	fmt.Print("")
}

//friends是好友列表，flag是在线状态，uin是临时QQ码，注意不是QQ号！
//markname是好友备注名，categories是好友分类组。
//vipinfo是会员信息，categories是好友分类信息
type friends struct {
	Retcode int
	Result  struct {
		Friends []struct {
			Flag       int `json:"flag"`
			Uin        int `json:"uin"`
			Categories int `json:"categories"`
		}
		Marknames []struct {
			Uin      int
			Markname string
		}
		Categories []struct {
			Index int
			Sort  int
			Name  string
		}
		Vipinfo []struct {
			Viplevel int
			U        int
			Isvip    int
		}
		Info []struct {
			Face int
			Flag int
			Nick string
			Uin  int
		}
	}
}
type Friends struct {
	Nick       string
	Markname   string
	Uin        string
	Categories int
	Msg_id     int
	Status     string
}

type Groups struct {
	Guin   string
	Name   string
	Msg_id int
}

func (qq *QQClient) GetFriendList() (rerr error) {
	defer func() {
		if e := recover(); e != nil {
			rerr = e.(error)
		}
	}()
	fmt.Printf("获取好友列表中\n")
	v := url.Values{}
	h, err := json.Marshal(
		map[string]string{
			`h`:       `hello`,
			`vfwebqq`: qq.Vfwebqq})
	ehandle(err)
	v.Set(`r`, string(h))

	re, err := qq.pForm(`http://s.web2.qq.com/api/get_user_friends2`, v)
	ehandle(err)
	fmt.Println("获得好友列表成功，解析中")
	frd := make(map[string]interface{})
	q := sRead(re.Body)
	err = json.Unmarshal([]byte(q), &frd)
	ehandle(err)
	for _, t := range frd[`result`].(map[string]interface{})[`categories`].([]interface{}) {
		qq.Friend.Categories[int(t.(map[string]interface{})[`index`].(float64))] = t.(map[string]interface{})[`name`].(string)
	}
	for _, t := range frd[`result`].(map[string]interface{})[`friends`].([]interface{}) {
		qq.Friend.Friendlist[strconv.FormatFloat(t.(map[string]interface{})[`uin`].(float64), 'g', 20, 64)] = &Friends{Uin: strconv.FormatFloat(t.(map[string]interface{})[`uin`].(float64), 'g', 20, 64), Categories: int(t.(map[string]interface{})[`categories`].(float64))}
	}
	for _, t := range frd[`result`].(map[string]interface{})[`info`].([]interface{}) {
		qq.Friend.Friendlist[strconv.FormatFloat(t.(map[string]interface{})[`uin`].(float64), 'g', 20, 64)].Nick = t.(map[string]interface{})[`nick`].(string)
	}
	for _, t := range qq.Friend.Friendlist {
		t.Markname = t.Nick
	}
	for _, t := range frd[`result`].(map[string]interface{})[`marknames`].([]interface{}) {
		ptr, ok := qq.Friend.Friendlist[strconv.FormatFloat(t.(map[string]interface{})[`uin`].(float64), 'g', 20, 64)]
		if ok == true {
			ptr.Markname = t.(map[string]interface{})[`markname`].(string)
		}
	}
	fmt.Printf("解析好友列表完成\n")
	return
}

func (qq *QQClient) GetGroupList() (rerr error) {
	defer func() {
		if e := recover(); e != nil {
			rerr = e.(error)
		}
	}()
	fmt.Printf("获取群列表中\n")
	v := url.Values{}
	c, err := json.Marshal(map[string]interface{}{
		`vfwebqq`: qq.Vfwebqq})
	ehandle(err)
	v.Set(`r`, string(c))
	re, err := qq.pForm(`http://s.web2.qq.com/api/get_group_name_list_mask2`, v)
	ehandle(err)
	fmt.Println("获得群列表成功，解析中")
	p := make(map[string]interface{})
	json.Unmarshal([]byte(sRead(re.Body)), &p)
	for _, con := range p[`result`].(map[string]interface{})[`gnamelist`].([]interface{}) {
		gid := strconv.FormatFloat(con.(map[string]interface{})[`gid`].(float64), 'g', 20, 64)
		name := con.(map[string]interface{})[`name`].(string)
		qq.Group[gid] = &Groups{Guin: gid, Name: name}
	}
	fmt.Println("群列表解析解析完毕")
	return
}

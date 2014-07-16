package webqq

import (
	"encoding/json"
	"fmt"
	. "webqq/tools"
	"webqq/tools/simplejson"
)

func (qq *Client) GetId(uin string, b ...*BuddyInfo) (str string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	if len(b) > 1 {
		return ``, fmt.Errorf(`too many buddyinfo objects`)
	}
	re, err := qq.get(`http://s.web2.qq.com/api/get_friend_uin2?tuin=` +
		uin + `&verifysession=&type=1&code=&vfwebqq=` +
		qq.vfwebqq +
		`&t=` + qq.timeStamp())

	if err != nil {
		panic(err)
	}
	defer re.Body.Close()
	j, err := simplejson.NewJson(ReadByte(re.Body))
	if err != nil {
		panic(err)
	}
	acc, err := j.Get(`result`).Get(`account`).Int()
	if err != nil {
		panic(err)
	}
	str = fmt.Sprintf("%d", acc)
	if len(b) == 1 {
		b[0].Id = str
	}
	return
}

type BuddyInfo struct {
	Nick     string `json:"nick"`
	markname string
	Face     int    `json:"face"`
	Id       string `json:"id"`
	Uin      int    `json:"uin"`
	status   string
	Birthday struct {
		Month int `json:"month"`
		Year  int `json:"year"`
		Day   int `json:"day"`
	} `json:"birthday"`
	CZodiac     int    `json:"shengxiao"`
	Occupation  string `json:"occupation"`
	Phone       string `json:"phone"`
	college     string `json:"college"`
	Constel     int    `json:"constel"`
	Homepage    string `json:"homepage"`
	Country     string `json:"country"`
	City        string `json:"city"`
	Email       string `json:"email"`
	Province    string `json:"province"`
	Gender      string `json:"gender"`
	Mobile      string `json:"mobile"`
	online      int
	Vip_info    int    `json:"vip_info"`
	Personal    string `json:"personal"`
	Stat        int    `json:"stat"`
	Blood       int    `json:"blood"`
	Client_type int    `json:"client_type"`
}

type cc struct {
	Retcode int
	Result  *BuddyInfo
}

func (qq *Client) GetInfo(uin string, b *BuddyInfo) error {
	re, err := qq.Get(`http://s.web2.qq.com/api/get_friend_info2?tuin=` + uin +
		`&verifysession=&code=&vfwebqq=` + qq.vfwebqq + `&t=` + qq.timeStamp())
	ErrHandle(err, `n`)
	byt := ReadByte(re.Body)
	err = json.Unmarshal(byt, &cc{0, b})
	return err
}

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
	PT_SysGMessage         = PollType("sys_g_msg")
	PT_GroupWebMessage     = PollType("group_web_message")
	PT_BuddylistChange     = PollType("buddylist_change")
	PT_DiscuMessage        = PollType("discu_message")
	PT_SessMessage         = PollType("sess_message")
)

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

	lg.Debug("ret code is %d", retcode)

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
				case PT_SysGMessage:
					r.pollType = PT_SysGMessage
				case PT_GroupWebMessage:
					r.pollType = PT_GroupWebMessage
				case PT_BuddylistChange:
					r.pollType = PT_BuddylistChange
				case PT_DiscuMessage:
					r.pollType = PT_DiscuMessage
				case PT_SessMessage:
					r.pollType = PT_SessMessage
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

			p := js.Get("p").MustString("")

			r := new(PollMessage)
			r.retCode = 116
			r.t = time.Now()
			r.value = []byte(p)
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
			r.retCode = 100006
			r.t = time.Now()
			r.errMsg = js.Get("errmsg").MustString("")
			ret = append(ret, r)
		}
	case 108:
		{
			r := new(PollMessage)
			r.retCode = 108
			r.t = time.Now()
			r.errMsg = js.Get("errmsg").MustString("")
			ret = append(ret, r)
		}
	case 112:
		{
			r := new(PollMessage)
			r.retCode = 112
			r.t = time.Now()
			r.errMsg = js.Get("errmsg").MustString("")
			ret = append(ret, r)
		}
	default:
		{
			err = fmt.Errorf("unknown ret code:%v", retcode)
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

	re, err := qq.postForm(`http://d.web2.qq.com/channel/poll2`, v)
	ErrHandle(err, `p`)
	defer re.Body.Close()
	retu = ReadByte(re.Body)
	lg.Debug("poll raw msg is %s", retu)
	return

}

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
		content = append(content, fmt.Sprint(contentjs.GetIndex(i).MustString("不支持文字以外的消息")))

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

type GroupMessage struct {
	PollType    PollType
	MessageId   int64
	FromUin     string
	ToUin       string
	MessageId2  int64
	MessageType int64
	ReplyIp     int64
	GroupCode   string
	SendUin     string
	Seq         int64
	Time        time.Time
	InfoSeq     int64
	Font        struct {
		Size  int
		Color string
		Style [3]int
		Name  string
	}
	Content []string
}

func (this *PollMessage) IsGroupMessage() (*GroupMessage, bool) {

	if this.retCode != 0 || this.pollType != PT_GroupMessage {
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
	groupCode := strconv.FormatInt(int64(js.Get("group_code").MustFloat64()), 10)
	sendUin := strconv.FormatInt(int64(js.Get("send_uin").MustFloat64()), 10)
	seq := int64(js.Get("seq").MustFloat64())
	infoSeq := int64(js.Get("info_seq").MustFloat64())
	t := time.Unix(int64(js.Get("time").MustFloat64()), 0)
	contentjs := js.Get("content")
	fontjs := contentjs.GetIndex(0).GetIndex(1)
	fontSize := int(fontjs.Get("size").MustFloat64())
	color := fontjs.Get("color").MustString("000000")
	fontName := fontjs.Get("name").MustString("")
	content := []string{}
	for i := 1; i < len(contentjs.MustArray([]interface{}{})); i++ {
		content = append(content, fmt.Sprint(contentjs.GetIndex(i).MustString("不支持文字以外的消息")))

	}

	ret := new(GroupMessage)
	ret.FromUin = fromUin
	ret.ToUin = toUin
	ret.MessageId = msgId
	ret.MessageId2 = msgId2
	ret.MessageType = msgType
	ret.ReplyIp = replyIp
	ret.InfoSeq = infoSeq
	ret.Seq = seq
	ret.SendUin = sendUin
	ret.GroupCode = groupCode
	ret.Time = t
	ret.PollType, _ = this.PollType()
	ret.Font.Color = color
	ret.Font.Name = fontName
	ret.Font.Size = fontSize
	ret.Font.Style = [3]int{0, 0, 0}
	ret.Content = content
	return ret, true

}

func (this *PollMessage) IsNewPtwebqq() (string, bool) {
	if this.retCode == 116 {
		return string(this.value), true
	} else {
		return ``, false
	}
}

func (this *PollMessage) IsOffline() bool {
	if this.retCode == 108 && this.retCode == 112 {
		return true
	} else {
		return false
	}
}

func (this *PollMessage) IsTips() (ret string, is bool) {
	if this.retCode != 0 || this.pollType != PT_Tips {
		return ``, false
	}
	return string(this.value), true
}
func (this *PollMessage) IsSystemMessage() (ret string, is bool) {
	if this.retCode != 0 || this.pollType != PT_SystemMessage {
		return ``, false
	}
	return string(this.value), true
}

// {"msg_id":6026,"from_uin":1502113816,"to_uin":173165159,"msg_id2":497528,
// "msg_type":45,"reply_ip":180028749,"group_code":2288613488,"group_type":1,
// "ver":1,"send_uin":291028157,
// "xml":"\u0001\u0000\u0001 \u0000\u0018m\u00BA\u0000\u000F\u00C8\u02FC\u00E4\u00B5\u00DA\u04BB\u00C7\u00E9-\u00B9\u00A8\u00ABh\u0000\u00EEtencent://miniplayer/?cmd=1\u0026fuin=125846471\u0026frienduin=125846471\u0026groupid=2188573954\u0026groupcode=198573954\u0026action=\u0027accept\u0027\u0026mdlurl=\u0027http://scenecgi.chatshow.qq.com/fcgi-bin/gm_qry_music_info.fcg?songcount=1\u0026songidlist=1600954\u0026version=207\u0026cmd=1\u0027\u0000\u00EEtencent://miniplayer/?cmd=1\u0026fuin=125846471\u0026frienduin=125846471\u0026groupid=2188573954\u0026groupcode=198573954\u0026action=\u0027refuse\u0027\u0026mdlurl=\u0027http://scenecgi.chatshow.qq.com/fcgi-bin/gm_qry_music_info.fcg?songcount=1\u0026songidlist=1600954\u0026version=207\u0026cmd=1\u0027"}}]}
func (this *PollMessage) IsGroupWebMessage() (ret string, is bool) {
	if this.retCode != 0 || this.pollType != PT_GroupWebMessage {
		return ``, false
	}
	//TODO
	return string(this.value), true
}
func (this *PollMessage) IsBuddylistChange() (ret string, is bool) {
	if this.retCode != 0 || this.pollType != PT_BuddylistChange {
		return ``, false
	}
	return string(this.value), true
}

type SysGMessage struct {
	AdminUin    string
	FromUin     string
	Gcode       string
	Message     string
	MessageId   int64
	MessageId2  int64
	MessageType int64
	ReplyIp     int64
	TGcode      string
	ToUin       string
	Type        string
}

func (this *PollMessage) IsSysGMessage() (ret *SysGMessage, is bool) {
	if this.retCode != 0 || this.pollType != PT_SysGMessage {
		return nil, false
	}
	ret = new(SysGMessage)
	js, _ := simplejson.NewJson(this.value)
	ret.AdminUin = fmt.Sprint(int64(js.Get("admin_uin").MustFloat64()))
	ret.FromUin = fmt.Sprint(int64(js.Get("from_uin").MustFloat64()))
	ret.Gcode = fmt.Sprint(int64(js.Get("gcode").MustFloat64()))
	ret.Message = js.Get("msg").MustString("")
	ret.MessageId = int64(js.Get("msg_id").MustFloat64())
	ret.MessageId2 = int64(js.Get("msg_id2").MustFloat64())
	ret.MessageType = int64(js.Get("msg_type").MustFloat64())
	ret.ReplyIp = int64(js.Get("admin_uin").MustFloat64())
	ret.TGcode = fmt.Sprint(int64(js.Get("t_gcode").MustFloat64()))
	ret.ToUin = fmt.Sprint(int64(js.Get("to_uin").MustFloat64()))
	ret.Type = js.Get("type").MustString("")
	//{"admin_uin":1.366473281e+09,"from_uin":4.217094352e+09,"gcode":3.762527505e+09,
	//"msg":"","msg_id":32284,"msg_id2":22284,"msg_type":36,"reply_ip":1.76488537e+08,
	//"t_gcode":2.13871351e+08,"to_uin":1.73165159e+08,"type":"group_request_join_agree"}

	return ret, true
}

type SessMessage struct {
	MessageId   int64
	FromUin     string
	ToUin       string
	MessageId2  int64
	MessageType int64
	ReplyIp     int64
	Time        time.Time
	Id          string
	Ruin        string
	ServiceType int64
	Flags       struct {
		Text  bool
		Pic   bool
		File  bool
		Audio bool
		Video bool
	}
	Font struct {
		Size  int
		Color string
		Style [3]int
		Name  string
	}
	Content []string
}

func (this *PollMessage) IsSessMessage() (ret *SessMessage, is bool) {
	if this.retCode != 0 || this.pollType != PT_SessMessage {
		return nil, false
	}
	js, _ := simplejson.NewJson(this.value)
	ret = new(SessMessage)
	ret.MessageId = int64(js.Get("msg_id").MustFloat64())
	ret.MessageId2 = int64(js.Get("msg_id2").MustFloat64())
	ret.ToUin = fmt.Sprint(int64(js.Get("to_uin").MustFloat64()))
	ret.FromUin = fmt.Sprint(int64(js.Get("from_uin").MustFloat64()))
	ret.MessageType = int64(js.Get("msg_type").MustFloat64())
	ret.ReplyIp = int64(js.Get("admin_uin").MustFloat64())
	ret.Time = time.Unix(int64(js.Get("time").MustFloat64()), 0)
	ret.Id = fmt.Sprint(int64(js.Get("id").MustFloat64()))
	ret.Ruin = fmt.Sprint(int64(js.Get("ruin").MustFloat64()))
	ret.ServiceType = int64(js.Get("service_type").MustFloat64())

	ret.Flags.Text, _ = strconv.ParseBool(fmt.Sprint(js.Get("flags").Get("text").MustInt(1)))
	ret.Flags.Audio, _ = strconv.ParseBool(fmt.Sprint(js.Get("flags").Get("audio").MustInt(1)))
	ret.Flags.File, _ = strconv.ParseBool(fmt.Sprint(js.Get("flags").Get("file").MustInt(1)))
	ret.Flags.Pic, _ = strconv.ParseBool(fmt.Sprint(js.Get("flags").Get("pic").MustInt(1)))
	ret.Flags.Video, _ = strconv.ParseBool(fmt.Sprint(js.Get("flags").Get("video").MustInt(1)))

	ret.Font.Size = js.Get("content").GetIndex(0).GetIndex(1).Get("size").MustInt(15)

	ret.Font.Color = js.Get("content").GetIndex(0).GetIndex(1).Get("color").MustString("000000")
	ret.Font.Name = js.Get("content").GetIndex(0).GetIndex(1).Get("name").MustString("宋体")
	for i := 1; i < len(js.Get("content").MustArray([]interface{}{})); i++ {
		ret.Content = append(ret.Content, fmt.Sprint(js.Get("content").GetIndex(i)))
	}
	return
}

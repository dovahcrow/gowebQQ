package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
	"webqq"
)

type cfg struct {
	Username      string
	Password      string
	DailTimeout   int
	VfCodeFile    string
	VfCodeTrigger string
	LogChannelLen int64
	KickRelogin   bool
	FileLog       struct {
		Use      bool
		Filename string
		Level    int
	}
	ConsoleLog struct {
		Use   bool
		Level int
	}
}

func main() {
	defer time.Sleep(1 * time.Second)
	config := new(cfg)
	_, err := toml.DecodeFile(`config.toml`, config)

	if err != nil {
		log.Fatalln("read config error:", err)
	}

	client := webqq.New(config.Username, config.Password, config.DailTimeout, config.LogChannelLen) //qq Number and Password,obtain a client object
	if config.FileLog.Use {
		client.SetLogger("file", fmt.Sprintf("{\"level\":%d,\"filename\":\"%s\"}", config.FileLog.Level, config.FileLog.Filename))
	}
	if config.ConsoleLog.Use {
		client.SetLogger("console", fmt.Sprintf("{\"level\":%d}", config.ConsoleLog.Level))
	}

L:
	ret, err := client.LoginStep1() //test login to know if a validation image is necessary
	if err != nil {
		client.Error("login step 1 fail: error: %v", err)
		client.Info("relogin after 1 minute")
		time.Sleep(1 * time.Minute)
		goto L
	}

	if ret.NeedPic() {
		re, err := client.Get(ret.PicAddr) //get pic addr
		if err != nil {
			client.Error("get login pic fail: error: %v", err)
		}

		f, err := os.Create(config.VfCodeFile) //write image to file
		if err != nil {
			client.Critical("create vfcode pic fail: error: %v", err)
			return
		}
		io.Copy(f, re.Body)
		f.Close()

		re.Body.Close()

		trigger := strings.Fields(config.VfCodeTrigger)

		cmd := exec.Command(trigger[0], trigger[1:]...)
		out, _ := cmd.StdoutPipe()
		cmd.Stdin = os.Stdin

		s := ``
		go func() { fmt.Fscan(out, &s) }()

		err = cmd.Run()
		if err != nil {
			client.Critical("cannot run vf trigger: error: %v", err)
			return
		}

		ret.SetVFCode(s) //set validation code
	}

	//if ret.NeedPic()==false ret.SetVFCode() is unnecessary

	err = client.LoginStep2(ret) //true login
	if err != nil {
		client.Error("login step2 fail: error: %v", err)
		client.Info("wait 10 second to relogin")
		time.Sleep(10 * time.Second)
		client.Warn(`relogin`)
		goto L
	}

	for {
		msgraw, err := client.RawPoll() //qq use poll to get message.
		if err != nil {
			client.Error("poll message fail: error: %v", err)
			client.Warn("relogin")
			goto L
		}
		msg, err := webqq.ParseRawPoll(msgraw)
		if err != nil {
			client.Error("parse raw poll fail: error: %v", err)
			continue
		}

		for _, v := range msg {
			if v.IsNothing() {
				continue
			}

			if m, ok := v.IsBuddyMessage(); ok {

				err := client.SendBuddyMsgEasy(m.FromUin, m.MessageId, "你好"+m.Content[0])
				if err != nil {
					fmt.Println(err)
				}
				gid, err := client.GetId(m.FromUin)
				client.Info("收到个人消息: %v ,id: %v", m.Content, gid)
				continue
			}
			if _, ok := v.IsBuddyStatusChange(); ok {
				continue
			}
			if m, ok := v.IsGroupMessage(); ok {
				err = client.SendGroupMsgEasy(m.FromUin, 55440000, "机器人重复说:"+m.Content[0])
				if err != nil {
					client.Error("reply group message fail: error: %v", err)
				}

				gid, err := client.GetId(m.GroupCode)
				if err != nil {
					client.Warn("get gid fail %v", err)
				}
				id, err := client.GetId(m.FromUin)
				if err != nil {
					client.Warn("get id fail %v", err)
				}
				client.Info("收到群消息: %v, 群号: %v 来自号码: %v", m.Content, gid, id)
				continue
			}
			if ptwebqq, ok := v.IsNewPtwebqq(); ok {
				client.SetPtWebqq(ptwebqq)
				continue
			}
			if ok := v.IsKick(); ok {
				client.Warn("你被踢下线了")
				if config.KickRelogin == true {
					goto L
				}
				return
			}
			if _, ok := v.IsSysGMessage(); ok {
				continue
			}
			if _, ok := v.IsTips(); ok {
				continue
			}
			if _, ok := v.IsGroupWebMessage(); ok {
				continue
			}
			if _, ok := v.IsSystemMessage(); ok {
				continue
			}
			if _, ok := v.IsBuddylistChange(); ok {
				continue
			}
			client.Warn("unexpected message: %v", v)
		}

	}

}

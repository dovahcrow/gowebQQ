package qqserver

import (
	"fmt"
	"os"
	"sender/qq/client"
	"sender/qq/session"
	"strings"
)

type qqServer struct {
	id, password string
	ch           chan *session.Session
}

var self = qqServer{id: `1296377482`, password: `N62G'luck`, ch: make(chan *session.Session, 10)}

var pool = map[string]*session.Session{}

func GetSession() *session.Session { return <-self.ch }

func Delete() {
	return
}

func cloze(id string) func() {
	return func() {
		delete(pool, id)
	}
}

func Start() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(`qqServer`, e)
			os.Exit(1)
		}
	}()
ALL:
	for ; ; fmt.Println(`leak4`) {

		client := qqclient.New(self.id, self.password)
		err := client.Login()
		if err != nil {
			panic(fmt.Errorf("登陆错误：%v", err))
		}
		for t_count := 0; ; fmt.Println(`leak3`) {
			if t_count > 10 {
				continue ALL
			}

			mp, err := client.PollSafe()
			if err != nil {
				if strings.Contains(err.Error(), "i/o timeout") {
					t_count += 1
					continue
				} else if err == qqclient.ELoginAgain {
					fmt.Println(`QQrelogin....`)
					continue ALL
				} else if err == qqclient.ENOMSG {
					continue
				} else {
					fmt.Printf("%v,%v\n", err, `poll`)
				}
			}

			for _, v := range mp {
				fmt.Println(`leak5`)
				s, ok := pool[v.From_uin]
				if ok == true {
					s.GotMsg(strings.Join(v.Body, ","))
				} else {
					id, err := client.GetId(v.From_uin)
					if err != nil {
						session.New(client, id, v.From_uin, cloze(v.From_uin)).ReplMsg(`对不起,无法获取你的QQ号`)
						break
					}

					sess := session.New(client, id, v.From_uin, cloze(v.From_uin))
					if v.Type == `addBuddy` {

					} else {
						sess.GotMsg(strings.Join(v.Body, ","))
					}
					pool[v.From_uin] = sess
					self.ch <- sess
				}

			}
		}
	}
}

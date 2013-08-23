package session

import (
	"fmt"
	"math/rand"
	"sender/qq/client"
	. "tools"
)

//qq

var me = `session`

func init() {
	fmt.Printf("")

}

type Session struct {
	msg_id  int
	uin     string
	id      string
	client  *qqclient.Client
	cloze   func()
	msgpool chan string
}

func (self *Session) Close() {
	self.cloze()
}
func (self *Session) GetMe() (string, error) {
	return self.id, nil
}
func New(cl *qqclient.Client, id string, uin string, cloze func()) *Session {
	return &Session{client: cl, msg_id: (rand.Intn(9000)+1000)*10000 + 1, uin: uin, id: id, msgpool: make(chan string, 100), cloze: cloze}
}

func (self *Session) RecvMsg() (msg string) {
	return <-self.msgpool
}

func (self *Session) GotMsg(msg string) {
	self.msgpool <- ParseChinese(msg)
}

func (self *Session) ReplMsg(msg string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	err = self.client.EasyBuddySend(self.uin, self.msg_id, msg)
	ErrHandle(err, `p`)
	self.msg_id++
	return nil
}

func (seld *Session) Id() string {
	return seld.id
}

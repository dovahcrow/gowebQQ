package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
	"webqq"
)

func main() {

	cl := webqq.New(`123`, `123`, 1000, 20000000) //qq Number and Password,obtain a client object
L:
	ret, err := cl.LoginStep1() //test login to know if a validation image is necessary
	if err != nil {
		log.Fatalln(err)

	}
	if ret.NeedPic() {
		re, err := cl.Get(ret.PicAddr) //get pic addr
		if err != nil {
			log.Fatalln(err)

		}

		f, err := os.Create(`cool.png`) //write image to file
		if err != nil {
			log.Fatalln(err)

		}
		io.Copy(f, re.Body)
		f.Close()

		re.Body.Close()

		exec.Command(`gwenview`, `cool.png`).Run() //show the image
		var s string
		fmt.Scanf("%s", &s) //recognize it
		ret.SetVFCode(s)    //set validation code
	}

	//if ret.NeedPic()==false ret.SetVFCode() is unnecessary

	err = cl.LoginStep2(ret) //true login
	if err != nil {
		time.Sleep(10 * time.Second)
		goto L

	}

	for {
		msgraw, err := cl.RawPoll() //qq use poll to get message.
		if err != nil {
			log.Fatalln(err)
		}

		msg, err := webqq.ParseRawPoll(msgraw)
		if err != nil {
			log.Fatalln(err)
		}
		for _, v := range msg {
			if v.IsNothing() {
				continue
			}

			if _, ok := v.IsBuddyMessage(); ok {
				// fmt.Printf("%+v\n", bmsg)
				// err := cl.SendBuddyMsgEasy(bmsg.FromUin, bmsg.MessageId+1, "got"+bmsg.Content[0])
				// if err != nil {
				// 	fmt.Println(err)
				// }
				continue
			}
			if _, ok := v.IsBuddyStatusChange(); ok {
				continue
			}
			if _, ok := v.IsGroupMessage(); ok {
				// err := cl.SendGroupMsgEasy(gmsg.FromUin, gmsg.MessageId, "fa")
				// if err != nil {
				// 	fmt.Println(err)
				// }
				// fmt.Printf("%+v\n", gmsg)
				continue
			}
			if ptwebqq, ok := v.IsNewPtwebqq(); ok {
				cl.SetPtWebqq(ptwebqq)
				continue
			}
			if ok := v.IsKick(); ok {
				return
			}
			if s, ok := v.IsSysGMessage(); ok {
				fmt.Printf("%+v\n", s)
				continue
			}
			if tips, ok := v.IsTips(); ok {
				fmt.Println("tips", tips)
				continue
			}
			if gwm, ok := v.IsGroupWebMessage(); ok {
				fmt.Println("gwm", gwm)
				continue
			}
			if gwm, ok := v.IsSystemMessage(); ok {
				fmt.Printf("%+v\n", gwm)
				continue
			}
			if gwm, ok := v.IsBuddylistChange(); ok {
				fmt.Println("blc", gwm)
				continue
			}
			fmt.Println(v)
		}

	}

}

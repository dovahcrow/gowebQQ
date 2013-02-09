package main

import (
	"fmt"
	//"strconv"
	"qqclient"
	"time"
	//"github.com/tncardoso/gocurses"
)

var pol chan int

func init() {
	fmt.Print("")
	pol = make(chan int, 1)
}

func main() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("错误：%v\n退出\n", e)
		}
	}()
	qq := qqclient.NewClient(`173165159`, `R9:T4K6@`)
	err := qq.Login()
	ehandle(err)
	err = qq.GetFriendList()
	ehandle(err)
	err = qq.GetGroupList()
	ehandle(err)
	ch := make(chan qqclient.Msg, 100)
	go func() {
		for {
			if len(pol) == 0 {
				go qq.Pollmsg(ch, pol)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	go func() {
		for {
			fmt.Printf("%+v\n", <-ch)
		}
	}()

	//t := ``
	for {

		//getinfo:
		toname := ``
		fmt.Scanf("%s", &toname)
		fmt.Printf("switch to %s\n", toname)
		qq.Sendgmsg(1, `不是凝聚`, toname)
		//for {
		//	fmt.Scanf("%s", &t)
		//	switch {
		//	case t == `:ls`:
		//		{
		//			for i, f := range qq.Friend.Friendlist {
		//				fmt.Printf("uin: %s		Name:%s\n", i, f.Markname)
		//			}
		//		}
		//	case t == `:ch`:
		//		goto getinfo

		//	case t[:5] == `:send`:
		//		{
		//			qq.Sendbmsg(qqclient.SNick, toname, t[5:])
		//		}
		//	}
		//}
	}

}

func ehandle(err error) {
	if err != nil {
		panic(err)
	}
}

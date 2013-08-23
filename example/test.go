package main

import (
	"fmt"
	"github.com/doomsplayer/Xgo-webqq/"
	"io"
	"os"
	"os/exec"
)

func main() {
	cl := qqclient.New(`123456789`, `123456789`) //qq Number and Password,obtain a client object
	ret, err := cl.TestLogin()                   //test login to know if a validation image is necessary
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if ret.NeedPic() {
		re, err := cl.Get(ret.PicAddr) //get pic addr
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		f, err := os.Create(`cool.png`) //write image to file
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
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

	err = cl.TrueLogin(ret) //true login
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fl, _ := cl.GetFriendList() //get friend list and catagory list
	for _, v := range fl.FriendMap {
		fmt.Printf("%+v\n", v)
	}
	for i, v := range fl.CataMap {
		fmt.Printf("%v:%v\n", i, v)
	}
	for {
		msg := <-cl.MessagePool //qq use poll to get message.
		//So this go version client starts a goroutine when login is succeeded automatically,
		//and put the message to the message pool.
		fmt.Println(msg)
		if msg.Type == `buddyMsg` && len(msg.Body) > 0 {
			cl.EasyBuddySend(msg.From_uin, 000000, msg.Body[0]) //second argument is a pseudo-random code and increases per message.000000 for first message and 000001 for second message etc.
		}

		if msg.Type == `kicked` && `offline` {
			//Relogin
			//if the poll goroutine get the kick message or offline message,it automantically locks itsself.
			//after relogin,you should unlock the poll lock manually.
			//when necessary,you can mannually lock it to stop poll process
			cl.PollMutex.Unlock()
		}
	}

}

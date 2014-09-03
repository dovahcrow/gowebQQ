package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func main() {

	f, err := os.Open("../log.log.2014-07-18.001")
	if err != nil {
		log.Fatalln(err)
	}
	b, _ := ioutil.ReadAll(f)
	s := string(b)
	reg := regexp.MustCompile(`(?m)parse raw poll: msg is {"retcode":0,"result":\[{"poll_type":"\w+"`)
	sub := reg.FindAllStringSubmatch(s, -1)
	for _, v := range sub {
		fmt.Println(v)
	}
}

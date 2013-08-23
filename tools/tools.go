package tools

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
)

func init() {

}

func ReadString(re io.Reader) string {
	bBody, _ := ioutil.ReadAll(re)
	return string(bBody)
}

func ReadByte(re io.Reader) []byte {
	bBody, _ := ioutil.ReadAll(re)
	return bBody
}

func ErrHandle(err error, c string, des ...string) {
	if err != nil {
		switch c {
		case `p`:

			{
				for _, v := range des {
					fmt.Println(v)
				}
				panic(err)
			}
		case `x`:
			{
				for _, v := range des {
					fmt.Println(v)
				}
				fmt.Println(err)
				os.Exit(0)
			}
		case `n`:
			{
				for _, v := range des {
					fmt.Println(v)
				}
				fmt.Println(err)
			}
		}
	}
}

func Md5Encode(in string) (out string) {
	h := md5.New()
	h.Write([]byte(in))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func JsMarshal(in interface{}) (out string) {
	b, err := json.Marshal(in)
	ErrHandle(err, `n`, `jsmarshal`)
	return string(b)
}

func PrintlnStruct(i interface{}) {
	t := reflect.TypeOf(i)
	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Field(i).Name, t.Field(i).Type.String())
	}

}

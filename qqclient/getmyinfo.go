package  qqclient

import (
	//"encoding/json"
	"fmt"
)

func init() {
	fmt.Print("")

}

func (qq *QQClient) GetOwnerInfo() {
	re, err := qq.Get(`http://s.web2.qq.com/api/get_friend_info2?tuin=` + qq.Id + `&verifysession=&code=&vfwebqq=` + qq.Vfwebqq + `&t=1338859742796`)
	ehandle(err)
	//c := make(map[string]interface{})
	//json.Unmarshal([]byte(sRead(re.Body)), &c)
	fmt.Printf("%+v\n", sRead(re.Body))

}

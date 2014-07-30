package webqq

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	. "webqq/tools"
)

func (this *Client) SearchQQ(id string) (bdinfo *BuddyInfo, err error) {
	re, err := this.Get(`http://captcha.qq.com/getimage?aid=1003901&0.3768496420234442`)

	f, _ := os.OpenFile("pp.png", os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	io.Copy(f, re.Body)
	f.Close()
	s := ``
	fmt.Scan(&s)

	url := fmt.Sprintf(
		"http://s.web2.qq.com/api/search_qq_by_uin2?tuin=%s&verifysession=%s&code=%s&vfwebqq=%s&t=%s",
		id,
		this.verifysession,
		s,
		this.vfwebqq,
		this.timeStamp())
	fmt.Println(url)
	re, err = this.Get(url)

	ErrHandle(err, `n`)
	byt := ReadByte(re.Body)
	fmt.Println(string(byt))
	bdinfo = new(BuddyInfo)
	err = json.Unmarshal(byt, &cc{0, bdinfo})
	return bdinfo, err
}

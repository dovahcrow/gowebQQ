package qqclient

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	. "tools"
	"tools/simplejson"
)

func init() {
	fmt.Print("")
}

type Friend struct {
	Uin        string
	Categories int
	Nick       string
	Markname   string
	Face       int
	Online     int
	Id         string
}

type FriendMap map[string]*Friend
type CataMap map[int]string

type FriendList struct {
	FriendMap
	CataMap
}

func (qq *Client) GetFriendList() (flst *FriendList, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	v := url.Values{}
	h, err := json.Marshal(
		map[string]string{
			`h`:       `hello`,
			`vfwebqq`: qq.vfwebqq,
			`hash`:    hash(qq.id, qq.ptwebqq)})

	v.Set(`r`, string(h))

	re, err := qq.postForm(`http://s.web2.qq.com/api/get_user_friends2`, v)
	if err != nil {
		panic(err)
	}
	defer re.Body.Close()
	flst, err = friendListParser(ReadByte(re.Body))
	if err != nil {
		panic(err)
	}
	for i := range flst.FriendMap {
		va, err := qq.GetId(i)
		if err != nil {
			panic(err)
		}
		flst.FriendMap[i].Id = va
	}

	return
}

func friendListParser(in []byte) (f *FriendList, err error) {
	js, err := simplejson.NewJson(in)
	if err != nil {
		return nil, err
	}
	f = &FriendList{}
	f.FriendMap = make(FriendMap)
	f.CataMap = make(CataMap)
	result := js.Get(`result`)
	for i, cat := 0, result.Get(`categories`); i < len(cat.MustArray()); i++ {
		f.CataMap[cat.GetIndex(i).Get(`index`).MustInt()] = cat.GetIndex(i).Get(`name`).MustString()
	}
	for i, frd := 0, result.Get(`friends`); i < len(frd.MustArray()); i++ {
		uin := fmt.Sprint(int64(frd.GetIndex(i).Get(`uin`).MustFloat64()))
		categories := frd.GetIndex(i).Get(`categories`).MustInt()
		f.FriendMap[uin] = &Friend{Uin: uin, Categories: categories}
	}
	for i, inf := 0, result.Get(`info`); i < len(inf.MustArray()); i++ {
		uin := fmt.Sprint(int64(inf.GetIndex(i).Get(`uin`).MustFloat64()))
		nick := inf.GetIndex(i).Get(`nick`).MustString()
		face := inf.GetIndex(i).Get(`face`).MustInt()
		flag := inf.GetIndex(i).Get(`flag`).MustInt()
		if ptr, ok := f.FriendMap[uin]; ok {
			ptr.Nick = nick
			ptr.Markname = nick
			ptr.Face = face
			ptr.Online = flag
		}
	}
	for i, mkn := 0, result.Get(`marknames`); i < len(mkn.MustArray()); i++ {
		markname := mkn.GetIndex(i).Get(`markname`).MustString()
		uin := fmt.Sprint(int64(mkn.GetIndex(i).Get(`uin`).MustFloat64()))
		if ptr, ok := f.FriendMap[uin]; ok {
			ptr.Markname = markname
		}
	}
	return
}

type b struct {
	s, e int
}

func nb(c, i int) b {
	return b{s: c | 0, e: i | 0}
}
func hash(is, a string) string {

	r := [4]int{}
	i, _ := strconv.Atoi(is)
	r[0] = i >> 24 & 255
	r[1] = i >> 16 & 255
	r[2] = i >> 8 & 255
	r[3] = i >> 0 & 255
	j := []int{}
	for e := 0; e < len(a); e++ {
		j = append(j, int(a[e]))
	}
	e := []b{}
	for e = append(e, nb(0, len(j)-1)); len(e) > 0; {
		var c = e[len(e)-1]
		e = e[:len(e)-1]
		if !(c.s >= c.e || c.s < 0 || c.e >= len(j)) {
			if c.s+1 == c.e {
				if j[c.s] > j[c.e] {
					j[c.s], j[c.e] = j[c.e], j[c.s]
				}
			} else {
				var f = j[c.s]
				var l = c.s
				var J = c.e
				for c.s < c.e {
					for c.s < c.e && j[c.e] >= f {
						c.e--
						r[0] = (r[0] + 3) & 255
					}
					if c.s < c.e {
						j[c.s] = j[c.e]
						c.s++
						r[1] = (r[1]*13 + 43) & 255
					}

					for c.s < c.e && j[c.s] <= f {
						c.s++
						r[2] = (r[2] - 3) & 255
					}
					if c.s < c.e {
						j[c.e] = j[c.s]
						c.e--
						r[3] = (r[0] ^ r[1] ^ r[2] ^ (r[3] + 1)) & 255
					}
				}

				j[c.s] = f
				e = append(e, nb(l, c.s-1))
				e = append(e, nb(c.s+1, J))
			}
		}
	}
	js := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
	es := ""
	for c := 0; c < len(r); c++ {
		es += js[r[c]>>4&15]
		es += js[r[c]&15]
	}
	return es
}

/*
{
    "retcode": 0,
    "result": {
        "friends": [
            {
                "flag": 0,
                "uin": 1192173647,
                "categories": 0
            },
            {
                "flag": 4,
                "uin": 2549892970,
                "categories": 1
            },
            {
                "flag": 0,
                "uin": 4069820047,
                "categories": 0
            },
            {
                "flag": 0,
                "uin": 490067261,
                "categories": 0
            },
            {
                "flag": 0,
                "uin": 4185452572,
                "categories": 0
            },
            {
                "flag": 8,
                "uin": 223290852,
                "categories": 2
            },
            {
                "flag": 8,
                "uin": 1196002108,
                "categories": 2
            }
        ],
        "marknames": [
            {
                "uin": 223290852,
                "markname": "沈超",
                "type": 0
            },
            {
                "uin": 1196002108,
                "markname": "大头",
                "type": 0
            }
        ],
        "categories": [
            {
                "index": 1,
                "sort": 1,
                "name": "道友"
            },
            {
                "index": 2,
                "sort": 2,
                "name": "同学"
            }
        ],
        "vipinfo": [
            {
                "vip_level": 0,
                "u": 1192173647,
                "is_vip": 0
            },
            {
                "vip_level": 0,
                "u": 2549892970,
                "is_vip": 0
            },
            {
                "vip_level": 2,
                "u": 4069820047,
                "is_vip": 1
            },
            {
                "vip_level": 0,
                "u": 490067261,
                "is_vip": 0
            },
            {
                "vip_level": 0,
                "u": 4185452572,
                "is_vip": 0
            },
            {
                "vip_level": 0,
                "u": 223290852,
                "is_vip": 0
            },
            {
                "vip_level": 0,
                "u": 1196002108,
                "is_vip": 0
            }
        ],
        "info": [
            {
                "face": 600,
                "flag": 29884928,
                "nick": "42",
                "uin": 1192173647
            },
            {
                "face": 12,
                "flag": 0,
                "nick": "火里栽莲",
                "uin": 2549892970
            },
            {
                "face": 0,
                "flag": 13107718,
                "nick": "林选",
                "uin": 4069820047
            },
            {
                "face": 336,
                "flag": 13107712,
                "nick": "度竹轻衣",
                "uin": 490067261
            },
            {
                "face": 558,
                "flag": 524802,
                "nick": "Qrox",
                "uin": 4185452572
            },
            {
                "face": 339,
                "flag": 13140480,
                "nick": "弎弍弌廿卅卌",
                "uin": 223290852
            },
            {
                "face": 147,
                "flag": 8388608,
                "nick": "blue~beak",
                "uin": 1196002108
            }
        ]
    }
}*/

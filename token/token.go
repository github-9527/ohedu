package Mtoken

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Login struct {
	S int `json:"s"`
	D D   `json:"d"`
}
type D struct {
	Token         string `json:"token"`
	StartTimeLong int64  `json:"start_time_long"`
	Effective     int    `json:"effective"`
}

func GetToken() string {
	attempts := 5

RECONNECT:
	req, err := http.NewRequest(http.MethodGet, "https://yun.ohedu.cn/httpapi/getToken.json", nil)
	if err != nil {
		panic(err)
	}

	params := make(url.Values)
	params.Add("account", "yzm")
	params.Add("password", "a@NuCoQ7ZMrgBUQis")

	req.URL.RawQuery = params.Encode()
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		//如果获取失败我们等待10秒
		attempts = attempts * 2
		fmt.Sprintf("登录请求发生错误%d秒后再次请求,\nerr：", attempts, err.Error())
		time.Sleep(time.Duration(attempts) * time.Second)

		//设置一个上线退出
		if attempts > 512 {
			return ""
		}
		goto RECONNECT
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		fmt.Println(err, "：token获取失败")
		return ""
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	login := &Login{}
	err = json.Unmarshal(body, login)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println("更新token=", login.D.Token)

	return login.D.Token
}

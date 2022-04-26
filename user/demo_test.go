package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

func TestFindAllUser(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://yun.ohedu.cn/httpapi/findModuleUsers.json?", nil)
	if err != nil {
		panic(err)
	}

	params := make(url.Values)
	//拼接字符
	//params.Add("apiToken", Mtoken.Token) // token每4小时获取一次
	params.Add("moduleId", "4698514001960960")
	//变更时间

	req.URL.RawQuery = params.Encode()
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		fmt.Println(err, "：用户获取失败")
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	user := &User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		fmt.Println(err)
		return
	}

	//创建ldap
	conn := Init(
		&LdapConn{
			ConnUrl:       "ldaps://192.168.1.10:636",
			BaseDn:        "OU=瓯海教育局,DC=howlink,DC=local",
			AdminAccount:  "CN=ctxadmin,CN=Users,DC=howlink,DC=local",
			Password:      "Howlink@1401",
			SslEncryption: true,
			Timeout:       5 * time.Second,
		})
	// ldap_conn, err := NewLdapConn(conn) // 建立ldap连接
	// if err != nil {
	// 	log.Printf("setup ldap connect failed,err:%v\n", err)
	// }
	//
	// //s := "cn=陈适,dc=howlink,dc=local"
	// passwordModifyRequest := ldap.NewPasswordModifyRequest("", "", "Howlink@1401")
	// result, err := ldap_conn.PasswordModify(passwordModifyRequest)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(result)

	allUser := FetchLdapUsers(conn)
	fmt.Println(allUser)

}

//测试定时器
func Test_adfs(t *testing.T) {

	// authorization := "Bearer YjFiMTA3MWM1ZDE5NDhmYTk5NzM5YmI1YjY0NGNhZmE=vWCI4KAdrqfBvfB3erb9"
	// data := strings.Split(authorization, " ")
	//
	// code := data[1][0:44]
	// token := data[1][44:]

	// //	fmt.Println(code, token)
	// 	i := strings.Index(data[1], "vWCI4KAdrqfBvfB3erb9")
	//
	// 	temp := strconv.Itoa(i)
	//
	// 	token := data[1][:strings.LastIndex(data[1], temp)]
	// 	fmt.Println(token, token)
	stat := cron.New(cron.WithSeconds())

	//定时器获取token,4小时一次
	log.Println("start")
	_, err := stat.AddFunc("0/20 * * * * ?", a)
	if err != nil {
		panic(err)
	}

	_, err = stat.AddFunc("0/20 * * * * ?", updateUser)
	if err != nil {
		panic(err)
	}
	stat.Start()

	time.Sleep(1 * time.Hour)
}

func a() {
	log.Println("b")
}

func updateUser() {
	log.Println("a")
}

func Test_reflect(t *testing.T) {
	a := make([]string, 4)
	a = append(a, "1")
	a = append(a, "2")

	b := make([]string, 4)
	b = append(b, "1")
	b = append(b, "2")
	fmt.Println(reflect.DeepEqual(a, b))
}

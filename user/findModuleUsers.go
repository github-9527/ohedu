package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/go-ldap/ldap"
	Mtoken "ohedu/token"
	"ohedu/util"
)

//获取用户信息
type User struct {
	D       []D    `json:"d"`
	ErrCode string `json:"err_code"`
	Msg     string `json:"msg"`
	S       int    `json:"s"`
	Total   int    `json:"total"`
}
type IdentityGroupList struct {
	CreateDate   string `json:"createDate"`
	DefaultGroup int    `json:"defaultGroup"`
	DisplayOrder int    `json:"displayOrder"`
	ID           int    `json:"id"`
	IdentityID   int    `json:"identityId"`
	Name         string `json:"name"`
}

type UserAuth struct {
	Ctime    string `json:"ctime"`
	ID       int    `json:"id"`
	ModuleID int    `json:"moduleId"`
	Role     int    `json:"role"`
	State    int    `json:"state"`
	UserID   int    `json:"userId"`
	Utime    string `json:"utime"`
}
type UserList struct {
	Address            string              `json:"address"`
	Birthday           string              `json:"birthday"`
	CreateDate         string              `json:"createDate"`
	DdUserid           string              `json:"dd_userid"`
	Descr              string              `json:"descr"`
	Email              string              `json:"email"`
	HeadImage          string              `json:"headImage"`
	HomePhone          string              `json:"homePhone"`
	ID                 int                 `json:"id"`
	Identity           string              `json:"identity"`
	IdentityGroupList  []IdentityGroupList `json:"identityGroupList"`
	IsDeleted          bool                `json:"isDeleted"`
	LastLoginDate      string              `json:"lastLoginDate"`
	LastLoginIP        string              `json:"lastLoginIp"`
	MobilePhone        string              `json:"mobilePhone"`
	MultiIdentity      string              `json:"multiIdentity"`
	MultiIdentityValue string              `json:"multiIdentityValue"`
	OnjobState         int                 `json:"onjob_state"`
	OrgID              int                 `json:"orgId"`
	PassWord           string              `json:"passWord"`
	Qq                 string              `json:"qq"`
	RealName           string              `json:"realName"`
	RealNameOfPingyin  string              `json:"realNameOfPingyin"`
	Role               int                 `json:"role"`
	Sex                int                 `json:"sex"`
	ShortTel           string              `json:"shortTel"`
	State              int                 `json:"state"`
	SysRole            int                 `json:"sysRole"`
	Unit               Unit                `json:"unit"`
	UnitCode           string              `json:"unitCode"`
	UpdateDate         string              `json:"updateDate"`
	UserAuth           UserAuth            `json:"userAuth"`
	UserName           string              `json:"userName"`
}
type UnitModelChilds struct {
	Code             string            `json:"code"`
	CreateDate       string            `json:"createDate"`
	DisplayOrder     int               `json:"displayOrder"`
	ID               int               `json:"id"`
	IsDeleted        bool              `json:"isDeleted"`
	IsReal           int               `json:"isReal"`
	Name             string            `json:"name"`
	NameOfPingyin    string            `json:"nameOfPingyin"`
	OfficeNature     string            `json:"officeNature"`
	ParentID         int               `json:"parentId"`
	SchoolIdentiCode string            `json:"schoolIdentiCode"`
	SchoolYear       string            `json:"schoolYear"`
	State            int               `json:"state"`
	TypeLevel1       int               `json:"typeLevel1"`
	TypeLevel2       int               `json:"typeLevel2"`
	TypeLevel3       int               `json:"typeLevel3"`
	UnitModelChilds  []UnitModelChilds `json:"unitModelChilds"`
	UpdateDate       string            `json:"updateDate"`
	UserList         []UserList        `json:"userList"`
}

type Unit struct {
	Code             string            `json:"code"`
	CreateDate       string            `json:"createDate"`
	DisplayOrder     int               `json:"displayOrder"`
	ID               int               `json:"id"`
	IsDeleted        bool              `json:"isDeleted"`
	IsReal           int               `json:"isReal"`
	Name             string            `json:"name"`
	NameOfPingyin    string            `json:"nameOfPingyin"`
	OfficeNature     string            `json:"officeNature"`
	ParentID         int               `json:"parentId"`
	SchoolIdentiCode string            `json:"schoolIdentiCode"`
	SchoolYear       string            `json:"schoolYear"`
	State            int               `json:"state"`
	TypeLevel1       int               `json:"typeLevel1"`
	TypeLevel2       int               `json:"typeLevel2"`
	TypeLevel3       int               `json:"typeLevel3"`
	UnitModelChilds  []UnitModelChilds `json:"unitModelChilds"`
	UpdateDate       string            `json:"updateDate"`
	UserList         []UserList        `json:"userList"`
}
type D struct {
	Address            string              `json:"address"`
	Birthday           string              `json:"birthday"`
	CreateDate         string              `json:"createDate"`
	DdUserid           string              `json:"dd_userid"`
	Descr              string              `json:"descr"`
	Email              string              `json:"email"`
	HeadImage          string              `json:"headImage"`
	HomePhone          string              `json:"homePhone"`
	ID                 int                 `json:"id"`
	Identity           string              `json:"identity"`
	IdentityGroupList  []IdentityGroupList `json:"identityGroupList"`
	IsDeleted          bool                `json:"isDeleted"`
	LastLoginDate      string              `json:"lastLoginDate"`
	LastLoginIP        string              `json:"lastLoginIp"`
	MobilePhone        string              `json:"mobilePhone"`
	MultiIdentity      string              `json:"multiIdentity"`
	MultiIdentityValue string              `json:"multiIdentityValue"`
	OnjobState         int                 `json:"onjob_state"`
	OrgID              int                 `json:"orgId"`
	PassWord           string              `json:"passWord"`
	Qq                 string              `json:"qq"`
	RealName           string              `json:"realName"`
	RealNameOfPingyin  string              `json:"realNameOfPingyin"`
	Role               int                 `json:"role"`
	Sex                int                 `json:"sex"`
	ShortTel           string              `json:"shortTel"`
	State              int                 `json:"state"`
	SysRole            int                 `json:"sysRole"`
	Unit               Unit                `json:"unit"`
	UnitCode           string              `json:"unitCode"`
	UpdateDate         string              `json:"updateDate"`
	UserAuth           UserAuth            `json:"userAuth"`
	UserName           string              `json:"userName"`
}

/*
afterTime 变更时间，查询在这个时间之后变更的记录格式：yyyy-MM-dd HH:mm:ss
haveAfterTime 是否要有变更时间
*/
func FindModuleUsers(afterTime time.Time, haveAfterTime bool) {

	req, err := http.NewRequest(http.MethodGet, "https://yun.ohedu.cn/httpapi/findModuleUsers.json?", nil)
	if err != nil {
		panic(err)
	}
	params := make(url.Values)
	//拼接字符
	params.Add("apiToken", Mtoken.GetToken()) //
	params.Add("moduleId", "4698514001960960")
	//变更时间
	if haveAfterTime {
		params.Add("afterTime", afterTime.String())
	}

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

	connUrl, err := util.Cfg.Section("Ldap").GetKey("connUrl")
	if err != nil {
		fmt.Println(err)
		return
	}

	baseDn, err := util.Cfg.Section("Ldap").GetKey("baseDn")
	if err != nil {
		fmt.Println(err)
		return
	}
	adminAccount, err := util.Cfg.Section("Ldap").GetKey("adminAccount")
	if err != nil {
		fmt.Println(err)
		return
	}
	password, err := util.Cfg.Section("Ldap").GetKey("password")
	if err != nil {
		fmt.Println(err)
		return
	}

	//创建ldap
	conn := Init(
		&LdapConn{
			ConnUrl:       connUrl.String(),
			BaseDn:        baseDn.String(),
			AdminAccount:  adminAccount.String(),
			Password:      password.String(),
			SslEncryption: false,
			Timeout:       5 * time.Second,
		})
	var LdapUsers []*LdapAttributes = make([]*LdapAttributes, 0, 5)
	// 读取配置文件
	domain, err := util.Cfg.Section("Ldap").GetKey("domain")
	if err != nil {
		fmt.Println(err)
		return
	}

	pwdLastSet, err := util.Cfg.Section("Ldap").GetKey("pwdLastSet")
	if err != nil {
		fmt.Println(err)
		return
	}

	accountCtl, err := util.Cfg.Section("Ldap").GetKey("accountCtl")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range user.D {

		//如果教育局接口已经删除，就我们也应该删除
		if v.IsDeleted {
			ldap_conn, err := NewLdapConn(conn) // 建立ldap连接
			if err != nil {
				fmt.Println(err)
			}
			delReq := ldap.NewDelRequest("CN="+v.RealName+","+conn.BaseDn, []ldap.Control{})
			if err := ldap_conn.Del(delReq); err != nil {
				fmt.Println(err)
			}
			continue
		}

		LdapUsers = append(LdapUsers,
			&LdapAttributes{
				DisplayName:       v.RealName,
				Dn:                "CN=" + v.RealName + "," + conn.BaseDn, // dn 				// 工号
				AccountCtl:        accountCtl.String(),
				Sam:               v.MobilePhone,
				PwdLastSet:        pwdLastSet.String(),
				Phone:             v.MobilePhone,
				Depart:            v.Unit.Name,
				UserPrincipalName: strconv.Itoa(v.ID) + "@" + domain.String(), //改成配置文件
				Name:              v.RealName,
				Status:            v.State,
			})

	}

	//新增用户
	res := AddLdapUsers(conn, LdapUsers)
	fmt.Println(res)
}

//
func TestName(t *testing.T) {

}

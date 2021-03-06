package user

import (
	"crypto/tls"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-ldap/ldap"
	"ohedu/logger"
	"ohedu/util"
)

// LdapConn LDAP服务器连接配置
type LdapConn struct {
	//gorm.Model
	// 连接地址
	ConnUrl string `json:"conn_url" gorm:"type:varchar(255);unique_index;not null;comment:连接地址 逻辑外键"`
	// SSL加密方式
	SslEncryption bool `json:"ssl_encryption" gorm:"type:tinyint;length:1;comment:SSL加密方式"`
	// 超时设置
	Timeout time.Duration `json:"timeout" gorm:"type:int;comment:超时设置"`
	// 根目录
	BaseDn string `json:"base_dn" gorm:"type:varchar(255);not null;comment:根目录"`
	// 用户名
	AdminAccount string `json:"admin_account" gorm:"type:varchar(255);not null;comment:用户名"`
	// 密码
	Password string `json:"password" gorm:"type:varchar(255);not null;comment:密码"`
}

type LdapAttributes struct {
	// ldap字段
	// Num         string `json:"employeeNumber" gorm:"type:varchar(100);unique_index"`     // 工号
	Sam        string `json:"sAMAccountName" gorm:"type:varchar(128);unique_index"`     // SAM账号
	Dn         string `json:"distinguishedName" gorm:"type:varchar(100);unique_index"`  // dn
	AccountCtl string `json:"UserAccountControl" gorm:"type:varchar(100);unique_index"` // 用户账户控制
	// Expire      string `json:"accountExpires" gorm:"type:varchar(100);unique_index"`     //  账户过期时间
	PwdLastSet string `json:"pwdLastSet" gorm:"type:varchar(100);unique_index"` //  用户下次登录必须修改密码
	// WhenCreated string `json:"whenCreated" gorm:"type:varchar(100);unique_index"`        //  创建时间
	// WhenChanged string `json:"whenChanged" gorm:"type:varchar(100);unique_index"`        //  修改时间
	DisplayName string `json:"displayName" gorm:"type:varchar(32);unique_index"` //  真实姓名
	// Sn          string `json:"sn" gorm:"type:varchar(100);unique_index"`                 //  姓
	Name string `json:"name" gorm:"type:varchar(100);unique_index"` // 姓名
	// GivenName   string `json:"givenName" gorm:"type:varchar(100);unique_index"`          // 名
	// Email       string `json:"mail" gorm:"type:varchar(128);unique_index"`               // 邮箱
	Phone string `json:"mobile" gorm:"type:varchar(32);unique_index"` // 移动电话
	// //Company     string `json:"company" gorm:"type:varchar(128);unique_index"`            // 公司
	Depart string `json:"department" gorm:"type:varchar(128);unique_index"` // 部门
	Title  string `json:"title" gorm:"type:varchar(100);unique_index"`      // 职务

	UserPrincipalName string `json:"userPrincipalName" gorm:"type:varchar(100);unique_index"` // 职务

	Status int `json:"status"`
}

var attrs = []string{
	"employeeNumber",     // 工号
	"sAMAccountName",     // SAM账号
	"distinguishedName",  // dn
	"UserAccountControl", // 用户账户控制
	"accountExpires",     // 账户过期时间
	"pwdLastSet",         // 用户下次登录必须修改密码
	"whenCreated",        // 创建时间
	"whenChanged",        // 修改时间
	"displayName",        // 显示名
	"sn",                 // 姓
	"name",
	"givenName",  // 名
	"mail",       // 邮箱
	"mobile",     // 手机号
	"company",    // 公司
	"department", // 部门
	"title",      // 职务
	"userPrincipalName",
}

// Init 实例化一个 ldapConn
func Init(c *LdapConn) *LdapConn {
	return &LdapConn{
		ConnUrl:       c.ConnUrl,
		SslEncryption: c.SslEncryption,
		Timeout:       c.Timeout,
		BaseDn:        c.BaseDn,
		AdminAccount:  c.AdminAccount,
		Password:      c.Password,
	}
}

// 获取ldap连接
func NewLdapConn(conn *LdapConn) (l *ldap.Conn, err error) {
	// 建立ldap连接
	l, err = ldap.DialURL(conn.ConnUrl)
	if err != nil {
		log.Printf("dial ldap url failed,err:%v", err)
		return
	}
	// 设置超时时间
	l.SetTimeout(time.Duration(conn.Timeout))
	if err != nil {
		log.Printf("dial ldap url failed,err:%v", err)
		return
	}
	// defer l.Close()

	// 重新连接TLS
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Printf("start tls failed,err:%v", err)
		return
	}

	// 首先与只读用户绑定
	err = l.Bind(conn.AdminAccount, conn.Password)
	if err != nil {
		log.Printf("admin user auth failed,err:%v", err)
		return
	}
	return
}

// 查询所有用户
func FetchLdapUsers(conn *LdapConn) (LdapUsers []*LdapAttributes) {
	ldap_conn, err := NewLdapConn(conn) // 建立ldap连接
	if err != nil {
		log.Printf("setup ldap connect failed,err:%v\n", err)
		return nil
	}
	defer ldap_conn.Close()

	searchRequest := ldap.NewSearchRequest(
		conn.BaseDn, // 待查询的base dn
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectclass=user)", // 过滤规则
		attrs,                // 待查询属性列表
		nil,
	)

	sr, err := ldap_conn.Search(searchRequest)
	if err != nil {
		log.Printf("查询用户出错:%v", err)
	}

	for _, entry := range sr.Entries {
		LdapUsers = append(LdapUsers,
			&LdapAttributes{
				//	Num:         entry.GetAttributeValue("employeeNumber"),
				Sam:        entry.GetAttributeValue("sAMAccountName"),
				Dn:         entry.GetAttributeValue("distinguishedName"),
				AccountCtl: entry.GetAttributeValue("UserAccountControl"),
				///	Expire:      entry.GetAttributeValue("accountExpires"),
				PwdLastSet: entry.GetAttributeValue("pwdLastSet"),
				//	WhenCreated: entry.GetAttributeValue("whenCreated"),
				//	WhenChanged: entry.GetAttributeValue("whenChanged"),
				DisplayName: entry.GetAttributeValue("displayName"),
				//	Sn:          entry.GetAttributeValue("sn"),
				Name: entry.GetAttributeValue("name"),
				//	GivenName:   entry.GetAttributeValue("givenName"),
				//	Email:       entry.GetAttributeValue("mail"),
				Phone: entry.GetAttributeValue("mobile"),
				//	Company:     entry.GetAttributeValue("company"),
				Depart: entry.GetAttributeValue("department"),
				//	Title:       entry.GetAttributeValue("title"),
			},
		)
	}
	return
}

// 批量新增用户 (AddLdapUsersRes []bool)
func AddLdapUsers(conn *LdapConn, LdapUsers []*LdapAttributes) (AddLdapUsersRes []bool) {
	ldap_conn, err := NewLdapConn(conn) // 建立ldap连接
	if err != nil {
		log.Printf("setup ldap connect failed,err:%v\n", err)
	}
	defer ldap_conn.Close()

	//查找已经创建用户，新增的时候对比，已存在账号并且有变化则更新,切片对比，要保证两个切片顺序相同

	// 批量处理  userAcountCoutr=66050禁用的意思   -2启用
	//
	for _, user := range LdapUsers {
		if user.Status == 1 {
			tempAccountCtl, err := strconv.Atoi(user.AccountCtl)
			if err != nil {
				fmt.Println(err)
				continue
			}
			tempAccountCtl = tempAccountCtl - 2
			user.AccountCtl = strconv.Itoa(tempAccountCtl)
		}

		//新增部分 "CN="+user.DisplayName+","+conn.BaseDn
		addReq := ldap.NewAddRequest(user.Dn, []ldap.Control{})                                    // 指定新用户的dn 会同时给cn name字段赋值
		addReq.Attribute("objectClass", []string{"top", "organizationalPerson", "user", "person"}) // 必填字段 否则报错 LDAP Result Code 65 "Object Class Violation"
		addReq.Attribute("sAMAccountName", []string{user.Sam})                                     // 登录名 必填
		addReq.Attribute("UserAccountControl", []string{user.AccountCtl})                          // 账号控制 544 是启用用户 		// 账号过期时间 当前时间加一个时间差并转换为NT时间
		addReq.Attribute("pwdLastSet", []string{user.PwdLastSet})                                  // 用户下次登录必须修改密码 0是永不过期
		addReq.Attribute("displayName", []string{user.DisplayName})                                // 真实姓名 某些系统需要
		addReq.Attribute("mobile", []string{user.Phone})                                           // 手机号 必填 某些系统需要
		addReq.Attribute("department", []string{user.Depart})
		addReq.Attribute("userPrincipalName", []string{user.UserPrincipalName})
		//添加
		modReq := ldap.NewModifyRequest(user.Dn, []ldap.Control{})
		modReq.Replace("sAMAccountName", []string{user.Sam})            // 登录名 必填
		modReq.Replace("UserAccountControl", []string{user.AccountCtl}) // 账号控制 544 是启用用户 		// 账号过期时间 当前时间加一个时间差并转换为NT时间
		modReq.Replace("pwdLastSet", []string{user.PwdLastSet})         // 用户下次登录必须修改密码 0是永不过期
		modReq.Replace("displayName", []string{user.DisplayName})       // 真实姓名 某些系统需要
		modReq.Replace("mobile", []string{user.Phone})                  // 手机号 必填 某些系统需要
		modReq.Replace("department", []string{user.Depart})
		modReq.Replace("userPrincipalName", []string{user.UserPrincipalName})

		if err = ldap_conn.Add(addReq); err != nil {
			if ldap.IsErrorWithCode(err, 68) {
				fetchldapUsers := FetchLdapUsers(conn)
				//获取所有的
				for _, v := range fetchldapUsers {
					if v.Dn == user.Dn {
						//修改数据
						fmt.Println(fetchldapUsers)
						if err := ldap_conn.Modify(modReq); err != nil {
							logger.Error("error enabling user account:", modReq, err)
						}
						log.Printf("User already exist: %s", err)
						break
					}
				}
				continue
				//账号已存在，数据有更新则修改数据
			} else {
				log.Printf("User insert error: %s", err)
			}
			AddLdapUsersRes = append(AddLdapUsersRes, false)
		}
		fmt.Println("创建成功", user.Name)

		accountPassword, err := util.Cfg.Section("Ldap").GetKey("accountPassword")
		if err != nil {
			logger.Error(err)
			return
		}

		//修改密码
		err = ChangePassword(user.Dn, accountPassword.String(), ldap_conn)
		if err != nil {
			fmt.Println(err)
			continue
		}

	}
	return
}

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"ohedu/config"
	"ohedu/logger"
	Mtoken "ohedu/token"
	"ohedu/user"
	"ohedu/util"
)

var (
	appId       string
	token       string
	state       string
	redirectUri string
	port        string
	code        string
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
	ErrCode     string `json:"errcode"`
	UserId      string `json:"UserId"`
	Email       string `json:"email"`
	Scope       string `json:"scope"`
	ExpiresIn   string `json:"expires_in"`
}

func main() {

	var err error
	router := gin.Default()

	stat := cron.New(cron.WithSeconds())

	//定时器获取token,3小时一次
	_, err = stat.AddFunc("0 0 0/3 * * ?", updateUser)
	if err != nil {
		panic(err)
	}
	stat.Start()
	//第一步
	router.GET("/oauth2/authorize", Authorize)
	router.GET("/oauth2/code", Getcode)
	router.POST("/oauth2/gettoken", GetToken)
	router.GET("/oauth2/userinfo", GetUserInfo)

	port, err := util.Cfg.Section("Citrix").GetKey("port")
	if err != nil {
		fmt.Println(err)
		return
	}

	router.Run(":" + port.String())

}

func GetUserInfo(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("Authorization")
	data := strings.Split(authorization, " ")

	code = data[1][:44]
	token = data[1][44:]

	result, _ := Base64DecodeString(code)
	userId, err := GetUserId(token, result)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx.JSON(200, userId)
}

func GetUserId(accessToken, code string) (*Id, error) {
	mids, err := util.Cfg.Section("Citrix").GetKey("mid")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	mid := mids.String()
	userIdUrl := "https://yun.ohedu.cn/yunapi/findUserByUserToken.json?" + "apiToken=" + accessToken + "&" + "token=" + code + "&" + "moduleId=" + mid

	resp, err := http.Get(userIdUrl)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	user1 := &Id{}
	err = json.Unmarshal(data, user1)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return user1, nil
}

// Base64DecodeString 解码
func Base64DecodeString(str string) (string, []byte) {
	resBytes, _ := base64.StdEncoding.DecodeString(str)
	return string(resBytes), resBytes
}

func GetToken(ctx *gin.Context) {
	code = ctx.Request.FormValue("code")

	access := &AccessToken{}
	access.AccessToken = Mtoken.GetToken()
	access.ErrCode = "0"
	access.Scope = "openid email profile id"

	//id_token, _:= CreateToken([]byte("secret"),"hhc", userId, user.Email)
	//access.IdToken = id_token
	access.TokenType = "Bearer"
	access.ExpiresIn = "3600"

	//token = access.AccessToken
	access.AccessToken = base64.StdEncoding.EncodeToString([]byte(code)) + access.AccessToken

	ctx.JSON(200, access)
}

func GetAccessToken(corpid, secret string) (AccessToken, error) {
	// https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=ID&corpsecret=SECRET
	url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?" + "corpid=" + corpid + "&" + "corpsecret=" + secret
	resp, err := http.Get(url)
	if err != nil {
		return AccessToken{}, err
	}
	var access AccessToken
	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	err = json.Unmarshal([]byte(data), &access)
	if err != nil {
		return access, err
	}

	return access, nil
}

func Getcode(ctx *gin.Context) {
	code = ctx.Request.FormValue("token")
	url := redirectUri + "?code=" + code + "&client_id=" + appId + "&state=" + state
	fmt.Println(url)
	//重定向
	ctx.Redirect(http.StatusMovedPermanently, url)
}

func init() {
	err := config.InitConfig("config/config.toml")
	if err != nil {
		panic(err)
	}

	err = logger.InitLogger()
	if err != nil {
		fmt.Println(err)
		return
	}
	Mtoken.GetToken()
	user.FindModuleUsers(time.Now(), false)
}

//更新ad用户的东西
func updateUser() {
	user.FindModuleUsers(time.Now(), true)
}

func Authorize(ctx *gin.Context) {
	redirectUri = ctx.Request.FormValue("redirect_uri")
	state = url.QueryEscape(ctx.Request.FormValue("state"))

	appId = ctx.Request.FormValue("client_id")

	ohedu_redirect_uri, err := util.Cfg.Section("Citrix").GetKey("OheduRedirectUri")
	if err != nil {
		logger.Error(err)
		return
	}

	url := "https://yun.ohedu.cn/aouth2/tologin" + "?" + "mid=4698514001960960" + "&" + "redirect_uri=" + ohedu_redirect_uri.String()

	//重定向
	ctx.Redirect(http.StatusMovedPermanently, url)

}

type Id struct {
	D D `json:"d"`
}
type IdentityGroupList struct {
	CreateDate   string `json:"createDate"`
	DefaultGroup int    `json:"defaultGroup"`
	DisplayOrder int    `json:"displayOrder"`
	ID           int    `json:"id"`
	IdentityID   int    `json:"identityId"`
	Name         string `json:"name"`
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
type D struct {
	ID int `json:"id"`
}

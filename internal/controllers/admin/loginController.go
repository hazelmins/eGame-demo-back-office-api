/*
 * @Description:后台登录相关方法
 */

package admin

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	services "eGame-demo-back-office-api/internal/services/admin"
	"eGame-demo-back-office-api/pkg/captcha/store"
	"eGame-demo-back-office-api/pkg/redisx"
	gstrings "eGame-demo-back-office-api/pkg/utils/strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
)

type loginController struct {
	BaseController
}

func NewLoginController() loginController {
	return loginController{}
}

func (con loginController) Routes(rg *gin.RouterGroup) {
	rg.GET("/captcha", con.captcha)
	/*******登录路由**********/
	rg.GET("/login", con.login)
	rg.POST("/login", con.login)
	rg.GET("/login_out", con.loginOut)
	rg.POST("/login_out", con.loginOut)
}

/*
* 登录
 */
func (con loginController) login(c *gin.Context) {
	if c.Request.Method == "GET" {
		con.Html(c, http.StatusOK, "home/login.html", gin.H{
			"title": "Egame backoffice",
		})
	} else {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// 为测试方便release模式才开启验证码
		if gin.Mode() == gin.ReleaseMode {

			captch := c.PostForm("captcha")
			var store = store.NewSessionStore(c, 20)
			verify := store.Verify("", captch, true)
			if !verify {
				con.Error(c, "验证码错误")
				return
			}

		}
		//setp 1 進db取玩家資料
		adminUser, err := services.NewAdminUserService().GetAdminUser(map[string]interface{}{"username": username})
		if err != nil {
			con.Error(c, "無此管理員")
			return
		}

		// 獲取GroupName
		groupName := adminUser.GroupName

		permissions, err := services.NewAdminGroupService().GetAdminGroup(groupName)
		if err != nil {
			con.Error(c, "無此管理組")
			return
		}

		// 現在您可以在permissions中獲得groupname的權限

		//判断密码是否正确
		if gstrings.Encryption(password, adminUser.Salt) == adminUser.Password {
			// 如果密碼驗證成功，創建用戶信息字典
			userInfo := make(map[string]interface{})
			userInfo["uid"] = adminUser.Uid
			userInfo["username"] = adminUser.Username
			userInfo["groupname"] = adminUser.GroupName
			userInfo["permissions"] = permissions

			// 將用戶信息序列化為 JSON 字符串
			userstr, _ := json.Marshal(userInfo)
			token := uuid.New().String()
			token = strings.Replace(token, "-", "", -1)

			// 將 JSON 字符串存儲到 Redis 中
			redisClient := redisx.GetRedisClient()
			err := redisClient.Set(token, string(userstr), time.Hour*1).Err()
			if err != nil {
				// 处理错误
				con.Error(c, "无法存储用户信息到 Redis")
				return
			}
			// 將 token 存儲到會話中，以便登出時使用
			session := sessions.Default(c)
			session.Set("token", token)
			session.Save()

			// 登录成功，重定向到 /admin/home 并显示成功消息
			con.Success2(c, "/admin/home", "登录成功", permissions)
		} else {
			// 登录失败，显示错误消息
			con.Error(c, "账号密码错误")
		}
	}

}

/**
* 登出
 */
func (con loginController) loginOut(c *gin.Context) {
	// 從會話中獲取 token
	session := sessions.Default(c)
	token := session.Get("token")
	if token != nil {
		// 從 Redis 中刪除 token
		redisClient := redisx.GetRedisClient()
		err := redisClient.Del(token.(string)).Err()
		if err != nil {
			// 处理错误
			con.Error(c, "無法刪除用戶信息從 Redis")
			return
		}
	}

	// 清除會話中的 token
	session.Delete("token")
	session.Delete("userInfo")
	session.Save()

	c.Redirect(http.StatusFound, "/admin/login")
}

/*
* 验证码
 */
func (con loginController) captcha(c *gin.Context) {

	var store = store.NewSessionStore(c, 20)
	driver := &base64Captcha.DriverString{
		Height: 60,
		Width:  150,
		Length: 4,
		Source: "abcdefghijklmnopqr234509867",
	}
	draw := base64Captcha.NewCaptcha(driver, store)
	_, b64s, err := draw.Generate()
	if err != nil {
		con.Error(c, "获取验证码失败")
	}

	i := strings.Index(b64s, ",")
	if i < 0 {
		log.Fatal("no comma")
	}
	// pass reader to NewDecoder
	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64s[i+1:]))

	io.Copy(c.Writer, dec)
}

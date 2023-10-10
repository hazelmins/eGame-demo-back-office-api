/*
 * @Description:後台中間件管理
 */
package middleware

import (
	"encoding/json"
	"net/http"

	"eGame-demo-back-office-api/pkg/casbinauth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

/*
*
用户登录验证
*/
func AdminUserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		type UserData struct { // redis中的資料結構
			Groupname   string          `json:"groupname"`
			Permissions map[string]bool `json:"permissions"`
			Token       string          `json:"token"`
			GroupUid    int             `json:"group uid"`
			Username    string          `json:"username"`
		}
		var userData UserData

		session := sessions.Default(c)
		userInfoJson := session.Get("userInfo")
		if userInfoJson == nil {
			// 取不到就是没有登录
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.String(200, ``)
			return
		}

		err := json.Unmarshal([]byte(userInfoJson.(string)), &userData)
		if err != nil {
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.String(200, ``)
			return
		}

		c.Set("username", userData.Username)
		c.Set("uid", userData.GroupUid)
		c.Set("groupname", userData.Groupname)
		c.Next()
	}
}

/*
*
用户权限验证
*/
func AdminUserPrivs() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := c.Get("username")
		if !ok {
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.String(200, ``)
			return
		}

		uri := c.FullPath()
		ok, err := casbinauth.Check(username.(string), uri, c.Request.Method)
		if !ok || err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": false,
				"msg":    "无权限禁止操作",
			})
			c.Abort()
		}
		c.Next()
	}
}

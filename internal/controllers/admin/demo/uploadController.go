//上传附件示例 上傳圖檔
package demo

import (
	"net/http"

	"eGame-demo-back-office-api/internal/controllers/admin"
	"eGame-demo-back-office-api/internal/models"
	services "eGame-demo-back-office-api/internal/services/admin"
	"eGame-demo-back-office-api/pkg/uploader"

	"github.com/gin-gonic/gin"
)

type uploadController struct {
	admin.BaseController
}

func NewUploadController() uploadController {
	return uploadController{}
}

func (con uploadController) Routes(rg *gin.RouterGroup) {
	rg.GET("/show", con.show)
	rg.POST("/upload", con.upload)
}

func (con uploadController) show(c *gin.Context) {

	con.Html(c, http.StatusOK, "demo/upload.html", gin.H{})

}

func (con uploadController) upload(c *gin.Context) {

	var (
		err error
		req models.UploadReq
	)
	err = con.FormBind(c, &req)
	if err != nil {
		con.Error(c, err.Error())
		return
	}
	req.Dst = "uploadfile"

	stor := uploader.LocalStorage{}

	filepath, err := services.NewUploadService().Save(stor, req)
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path": filepath,
	})

}

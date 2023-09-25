/*
 * @Description:上传附件服务
 * @Author: gphper
 * @Date: 2021-07-18 17:52:20
 */
package admin

import (
	"sync"

	"eGame-demo-back-office-api/internal/models"
	"eGame-demo-back-office-api/pkg/uploader"
)

type uploadService struct {
}

var (
	instanceUploadService *uploadService
	onceUploadService     sync.Once
)

func NewUploadService() *uploadService {
	onceUploadService.Do(func() {
		instanceUploadService = &uploadService{}
	})
	return instanceUploadService
}

func (ser *uploadService) Save(storage uploader.Storage, req models.UploadReq) (string, error) {
	return storage.Save(req.File, req.Dst)
}

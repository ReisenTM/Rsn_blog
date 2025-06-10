package image_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	hash2 "blogX_server/utils/hash"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

// 图片转存api
type TransferDepositRequest struct {
	Url string `json:"url" binding:"required" `
}

func (ImageApi) TransferDepositView(c *gin.Context) {
	cr := middleware.GetBind[TransferDepositRequest](c)
	//尝试请求图片
	response, err := http.Get(cr.Url)
	if err != nil {
		resp.FailWithMsg("图片请求错误", c)
		return
	}

	suffix := "png"
	switch response.Header.Get("Content-Type") {
	case "image/avif":
		suffix = "avif"
	}
	byteData, _ := io.ReadAll(response.Body)
	hash := hash2.Md5(byteData)

	filePath := fmt.Sprintf("uploads/%s/%s.%s", global.Config.Upload.UploadDir, hash, suffix)

	err = os.WriteFile(filePath, byteData, 0666)
	if err != nil {
		logrus.Error(err)
		resp.FailWithMsg("图片保存失败", c)
		return
	}
	resp.OkWithData("/"+filePath, c)
}

package image_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	fileTools "blogX_server/utils/file"
	"blogX_server/utils/hash"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
)

// UploadImageView 图片上传
// 需要满足：文件后缀；文件大小；文件不重复
func (api ImageApi) UploadImageView(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	//文件大小判断
	s := global.Config.Upload.Size * 1024 * 1024
	if fileHeader.Size > s {
		resp.FailWithMsg(fmt.Sprintf("文件大小大于%dMB", s), c)
		return
	}
	//文件后缀判断
	filename := fileHeader.Filename
	size := fileHeader.Size
	suffix, err := fileTools.ImageSuffixJudge(filename)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	//图片hash
	file, err := fileHeader.Open()
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	hashed := hash.Md5(data)
	var models model.ImageModel
	err = global.DB.Take(&models, "hash = ?", hashed).Error
	if err == nil {
		//如果找到了
		logrus.Infof("上传图片重复 %s <==> %s  %s", filename, models.Filename, hashed)
		resp.OK(models.NetPath(), "上传成功", c)
		return
	}
	//如果以前没有上传过，保存到数据库
	filePath := fmt.Sprintf("uploads/%s/%s.%s", global.Config.Upload.UploadDir, hashed, suffix)
	newModel := model.ImageModel{
		Filename: filename,
		Hash:     hashed,
		Path:     filePath,
		Size:     size,
	}
	err = global.DB.Create(&newModel).Error
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	err = c.SaveUploadedFile(fileHeader, filePath)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	resp.OK(filePath, "图片上传成功", c)
	return
}

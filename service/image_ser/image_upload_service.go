package image_ser

import (
	"fmt"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/ctype"
	"gvb_server/plugins/qiniu"
	"gvb_server/utils"
	"io"
	"mime/multipart"
	"path"
	"strings"
)

var (
	// WhiteImageList 图片上传的白名单
	WhiteImageList = []string{
		"jpg",
		"png",
		"jpeg",
		"gif",
		"ico",
		"tiff",
		"svg",
		"webp",
	}
)

type FileUploadResponse struct {
	FileName  string `json:"file_name"`  // 文件名
	IsSuccess bool   `json:"is_success"` // 是否上传成功
	Msg       string `json:"msg"`        // 消息
}

// 向数据库添加数据，如果重复就不要添加啦
// ImageUploadService 文件上传方法
func (ImageService) ImageUploadService(file *multipart.FileHeader) (res FileUploadResponse) {
	// 判断图片是否存在在白名单中
	fileName := file.Filename
	basePath := global.Config.Upload.Path
	filePath := path.Join(basePath, file.Filename) // 存储路径

	res.FileName = filePath

	// 判断文件白名单
	nameList := strings.Split(fileName, ".") // 图片名字切片
	suffix := nameList[len(nameList)-1]      // 获取后缀名
	if !utils.Inlist(suffix, WhiteImageList) {
		res.Msg = "非法文件"
		return
	}

	// 判断大小
	// 加float64是因为要进行浮点数除法
	size := float64(file.Size) / float64(1024*1024)
	if size > float64(global.Config.Upload.Size) {
		res.Msg = fmt.Sprintf("图片大小超过设定大小，当前大小为：%.2fMB，设定大小为：%dMB", size, global.Config.Upload.Size)
		return
	}

	// 读取文件内容 hash
	fileObj, err := file.Open()
	if err != nil {
		global.Log.Error(err)
	}
	byteData, err := io.ReadAll(fileObj)
	imageHash := utils.Md5(byteData)
	// 去数据库中查这个图片是否存在
	var bannerModel models.BannerModel
	err = global.DB.Take(&bannerModel, "hash = ?", imageHash).Error
	if err == nil {
		// 找到了
		res.FileName = bannerModel.Path
		res.Msg = "图片已存在"
		return
	}

	fileType := ctype.Local
	res.Msg = "图片上传成功"
	res.IsSuccess = true

	// 是否上传到七牛云
	if global.Config.QiNiu.Enable {
		filePath, err = qiniu.UploadImage(byteData, fileName, global.Config.QiNiu.Prefix)
		if err != nil {
			global.Log.Error(err)
			res.Msg = err.Error()
			return
		}
		res.FileName = filePath
		res.Msg = "上传七牛云成功"
		fileType = ctype.QiNiu
	}

	// 图片入库
	global.DB.Create(&models.BannerModel{
		Path:      filePath,
		Hash:      imageHash,
		Name:      fileName,
		ImageType: fileType,
	})
	return
}

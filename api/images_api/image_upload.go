package images_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/ctype"
	"gvb_server/models/res"
	"gvb_server/plugins/qiniu"
	"gvb_server/utils"
	"io"
	"os"
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

// 图片上传的响应
type FileUploadResponse struct {
	FileName  string `json:"file_name"`  // 文件名
	IsSuccess bool   `json:"is_success"` // 是否上传成功
	Msg       string `json:"msg"`        // 消息
}

// 上传图片，返回图片的URL
func (ImagesApi) ImageUploadView(c *gin.Context) {
	// 上传多个图片文件
	form, err := c.MultipartForm()
	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}

	fileList, ok := form.File["images"]
	if !ok {
		res.FailWithMessage("不存在的文件", c)
		return
	}

	// 判断路径是否存在
	// 截断（逐级判断）
	basePath := global.Config.Upload.Path
	// 在Go语言中，os.ReadDir(basePath)函数用于读取指定路径basePath下的所有目录条目，并将它们以文件名排序后返回。
	// 这个函数返回一个FileInfo结构的数组，每个元素代表一个目录条目。
	// 这个函数可以用来列出一个目录下的所有文件和子目录。
	_, err = os.ReadDir(basePath)
	if err != nil {
		// 在Go语言中，os.MkdirAll(basePath, os.ModePerm)函数用于创建指定路径basePath下的所有目录，包括任何必要的父目录。
		// 这个函数返回nil，或者在出错时返回一个错误。
		// 第二个参数os.ModePerm是一个权限位，用于设置所有被该函数创建的目录的读/写/执行权限。
		// 在Unix-like系统中，os.ModePerm等价于07772，意味着用户有权列出、修改和搜索目录中的文件。
		// 如果目录已经存在，os.MkdirAll()函数不会做任何事情，而是返回nil。
		// 递归创建
		err = os.MkdirAll(basePath, os.ModePerm)
		if err != nil {
			global.Log.Error(err)
		}
	}

	// 不存在就创建
	// 图片上传结果数组
	var resList []FileUploadResponse

	for _, file := range fileList {
		// 存数据库。。。。

		// 判断图片是否存在在白名单中
		fileName := file.Filename
		nameList := strings.Split(fileName, ".") // 图片名字切片
		suffix := nameList[len(nameList)-1]      // 获取后缀名
		if !utils.Inlist(suffix, WhiteImageList) {
			resList = append(resList, FileUploadResponse{
				FileName:  file.Filename,
				IsSuccess: false,
				Msg:       "非法文件",
			})
			continue
		}

		// 图片存储
		filePath := path.Join(basePath, file.Filename) // 存储路径

		// 判断大小
		// 加float64是因为要进行浮点数除法
		size := float64(file.Size) / float64(1024*1024)
		if size > float64(global.Config.Upload.Size) {
			resList = append(resList, FileUploadResponse{
				FileName:  file.Filename,
				IsSuccess: false,
				Msg:       fmt.Sprintf("图片大小超过设定大小，当前大小为：%.2fMB，设定大小为：%dMB", size, global.Config.Upload.Size),
			})
			continue
		}

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
			resList = append(resList, FileUploadResponse{
				FileName:  bannerModel.Path,
				IsSuccess: false,
				Msg:       "图片已存在",
			})
			continue
		}

		// 是否上传到七牛云
		if global.Config.QiNiu.Enable {
			filePath, err = qiniu.UploadImage(byteData, fileName, "gvb")
			if err != nil {
				global.Log.Error(err)
				continue
			}
			resList = append(resList, FileUploadResponse{
				FileName:  filePath,
				IsSuccess: true,
				Msg:       "上传七牛云成功",
			})

			// 图片入库
			global.DB.Create(&models.BannerModel{
				Path:      filePath,
				Hash:      imageHash,
				Name:      fileName,
				ImageType: ctype.QiNiu,
			})
			continue
		}

		// 上传
		err = c.SaveUploadedFile(file, filePath)
		// 上传失败
		if err != nil {
			global.Log.Error(err)
			resList = append(resList, FileUploadResponse{
				FileName:  file.Filename,
				IsSuccess: false,
				Msg:       err.Error(),
			})
			continue
		}

		// 上传成功
		resList = append(resList, FileUploadResponse{
			FileName:  filePath,
			IsSuccess: true,
			Msg:       "上传成功",
		})

		// 图片入库
		global.DB.Create(&models.BannerModel{
			Path:      filePath,
			Hash:      imageHash,
			Name:      fileName,
			ImageType: ctype.Local,
		})
	}

	res.OKWithData(resList, c)
}

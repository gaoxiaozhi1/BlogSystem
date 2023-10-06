package images_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models/res"
	"gvb_server/service"
	"gvb_server/service/image_ser"
	"os"
)

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
	var resList []image_ser.FileUploadResponse

	for _, file := range fileList {
		// 上传文件
		serviceRes := service.ServiceApp.ImageService.ImageUploadService(file)
		// 没有上传成功的话
		if !serviceRes.IsSuccess {
			resList = append(resList, serviceRes)
			continue
		}

		// 上传成功, 不上传到七牛，那么本地还需要保存一下
		if !global.Config.QiNiu.Enable {
			err = c.SaveUploadedFile(file, serviceRes.FileName)
			// 上传失败
			if err != nil {
				global.Log.Error(err)
				serviceRes.Msg = err.Error()
				serviceRes.IsSuccess = false
				resList = append(resList, serviceRes)
				continue
			}
		}

		resList = append(resList, serviceRes)
	}

	res.OKWithData(resList, c)
}

package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goTouch/TicTok_SimpleVersion/dao"
	"github.com/goTouch/TicTok_SimpleVersion/domain"
	"github.com/goTouch/TicTok_SimpleVersion/util"
	"github.com/minio/minio-go/v7"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

// 既然是发布视频，首先需要校验token，登入的问题
// Publish 获取用户投稿的视频并保存到本地
func Publish(c *gin.Context) {
	//id := c.GetInt64("userId")
	//if err != nil {
	//	log.Println(err)
	//	log.Println("Publish接口：当前用户token核验失败")
	//}
	id := "1"
	// 获取用户上传的视频
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 视频文件的后缀，也即视频的格式
	fileSuffix := filepath.Ext(data.Filename)
	// 通过用户id和当前时间戳拼接成最终存放的视频文件名
	filename := fmt.Sprintf("%x_%x%s", id, time.Now().Unix(), fileSuffix)
	// 拼接存放视频的本地路径
	//saveFile := filepath.Join("./static-server/videos/", filename)
	file, err := data.Open()
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	defer file.Close()
	miniodata, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	cclient := dao.MinioClient

	log.Println(cclient.EndpointURL())
	log.Println(dao.MinioClient.EndpointURL())
	_, err = dao.MinioClient.PutObject(
		c,
		util.VidioBucketName,
		filename,
		bytes.NewBuffer(miniodata),
		int64(len(miniodata)),
		minio.PutObjectOptions{UserMetadata: map[string]string{"x-amz-acl": "public-read"}},
	)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	reqParams := make(url.Values)
	presignedURL, err := dao.MinioClient.PresignedGetObject(c, util.VidioBucketName, filename, time.Duration(1000)*time.Second, reqParams)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(presignedURL)
	c.JSON(http.StatusOK, domain.Response{
		StatusCode: 0,
		StatusMsg:  presignedURL.String(),
	})
	return
	// 保存视频文件到本地
	//if err = c.SaveUploadedFile(data, saveFile); err != nil {
	//	c.JSON(http.StatusOK, domain.Response{
	//		StatusCode: 1,
	//		StatusMsg:  err.Error(),
	//	})
	//	return
	//}

	// 封面的文件名和最终存放的视频文件名一致，因为封面是图片，所以把后缀改为jpg
	//covername := strings.TrimSuffix(filename, fileSuffix) + ".jpg"
	// 拼接存放封面的本地路径
	//saveCover := filepath.Join("./static-server/covers/", covername)
	//log.Println(saveCover)
	//isGenerateOK := coverGenerator(saveFile, saveCover)
	//
	//if err = service.Publish(filename, covername, c.PostForm("title"), id, isGenerateOK); err != nil {
	//	c.JSON(http.StatusOK, dao.Response{
	//		StatusCode: 1,
	//		StatusMsg:  err.Error(),
	//	})
	//	return
	//}
	//c.JSON(http.StatusOK, Response{
	//	StatusCode: 0,
	//	StatusMsg:  filename + " uploaded successfully",
	//})
}

// 既然是发布视频，首先需要校验token，登入的问题

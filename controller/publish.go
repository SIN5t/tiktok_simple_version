package controller

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goForward/tictok_simple_version/config"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goForward/tictok_simple_version/dao"
	"github.com/goForward/tictok_simple_version/domain"
	"github.com/goForward/tictok_simple_version/service"
	"github.com/goForward/tictok_simple_version/util"
	"github.com/minio/minio-go/v7"
)

// Publish 获取用户投稿的视频并保存到本地
func Publish(c *gin.Context) {
	id := c.GetInt64("userId")
	title := c.PostForm("title")
	//if err != nil {
	//	log.Println(err)
	//	log.Println("Publish接口：当前用户token核验失败")
	//}
	//id := "1"
	// 获取用户上传的视频
	data, err := c.FormFile("data")
	if err != nil {
		log.Println(err.Error())
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
	//saveFile := filepath.Join("./static-server/video/", filename)
	file, err := data.Open()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	defer file.Close()
	miniodata, err := io.ReadAll(file)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	_, err = dao.MinioClient.PutObject(
		c,                          // 上下文。
		config.VideoBucketName,     // 存储桶名称。
		filename,                   // 存储对象名称。
		bytes.NewBuffer(miniodata), // 读取对象的内容。
		int64(len(miniodata)),      // 对象的大小。
		minio.PutObjectOptions{UserMetadata: map[string]string{"x-amz-acl": "miniodata-read"}}, // minio.PutObjectOptions，用户可以通过这个参数设置对象的元数据。
	)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	reqParams := make(url.Values)
	presignedURL, err := dao.MinioClient.PresignedGetObject(
		c,                               // 上下文。
		config.VideoBucketName,          // 存储桶名称。
		filename,                        // 存储对象名称。
		time.Duration(1000)*time.Second, // 过期时间，有效期为1小时。
		reqParams,                       // minio.RequestHeaders，用户可以通过这个参数设置请求头。
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(presignedURL)
	coverGenerateStatus := false
	jpeg, err := util.ReadFrameAsJpeg(presignedURL.String(), 1)

	if err == nil {
		buf := &bytes.Buffer{}
		buf.ReadFrom(jpeg)
		putSnapshotToOss(buf, config.PictureBucketName, filename+".jpg")
		coverGenerateStatus = true
	} else {
		log.Println(err)
		coverGenerateStatus = false
	}
	err = service.InsertVideos(filename, title, filename+".jpg", id, coverGenerateStatus)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, domain.Response{
		StatusCode: 0,
		StatusMsg:  "上传成功",
	})
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

func putSnapshotToOss(buf *bytes.Buffer, bucketName string, saveName string) {
	_, err := dao.MinioClient.PutObject(context.Background(),
		bucketName,
		saveName,
		buf,
		int64(buf.Len()),
		minio.PutObjectOptions{
			ContentType: "image/jpeg",
		})

	if err != nil {
		log.Fatalln("图片上传失败", err)
		return
	}
	log.Printf("图片上传成功")
}

func PublishList(c *gin.Context) {
	idStr := c.Query("user_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: "用户id解析错误"})
		return
	}
	list, err := service.QueryAuthorPublishedVideo(id)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.VideoListResponse{
		Response:  domain.Response{StatusCode: 0},
		VideoList: list,
	})
}

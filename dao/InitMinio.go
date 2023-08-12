package dao

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	endpoint          = "127.0.0.1:9000"
	accessKeyID       = "minioadmin"
	secretAccessKey   = "minioadmin"
	useSSL            = false
	VidioBucketName   = "vidio"
	PictureBucketName = "picture"
)

var (
	MinioClient *minio.Client
	err         error
	ctx         = context.Background()
)

func InitMinio() {
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL})
	if err != nil {
		log.Fatalln("minio连接错误: ", err)
	} else {
		log.Printf("%#v\n", MinioClient)
		log.Println(MinioClient.EndpointURL())
		createBucket(VidioBucketName)
		createBucket(PictureBucketName)
	}

}
func createBucket(bucketName string) {

	err = MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "local"})
	if err != nil {
		// 查询桶是否存在，如果存在则打印已存在，否则打印错误
		exists, errBucketExists := MinioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("已存在 %s\n", bucketName)
		} else {
			log.Fatalln("minio创建错误 %s", err)
		}
	} else {
		log.Printf("成功创建 %s\n", bucketName)
	}
}

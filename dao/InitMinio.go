package dao

import (
	"context"
	"github.com/goForward/tictok_simple_version/config"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	useSSL = false
)

var (
	MinioClient *minio.Client
	err         error
	ctx         = context.Background()
)

func InitMinio() {
	MinioClient, err = minio.New(config.GetMinioEndpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(config.GetMinioAccessKeyID(), config.GetMinioSecretAccessKey(), ""),
		Secure: useSSL})
	if err != nil {
		log.Fatalln("minio连接错误: ", err)
	} else {
		log.Printf("%#v\n", MinioClient)
		log.Println(MinioClient.EndpointURL())
		createBucket(config.VideoBucketName)
		createBucket(config.PictureBucketName)
	}

}
func createBucket(bucketName string) {

	err = MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "local"})
	s := "{\"Version\":\"2012-10-17\"," +
		"\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":" +
		"{\"AWS\":[\"*\"]},\"Action\":[\"s3:ListBucket\",\"s3:ListBucketMultipartUploads\"," +
		"\"s3:GetBucketLocation\"],\"Resource\":[\"arn:aws:s3:::" + bucketName +
		"\"]},{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":[\"*\"]},\"Action\":[\"s3:PutObject\",\"s3:AbortMultipartUpload\",\"s3:DeleteObject\",\"s3:GetObject\",\"s3:ListMultipartUploadParts\"],\"Resource\":[\"arn:aws:s3:::" +
		bucketName +
		"/*\"]}]}"
	if err != nil {
		// 查询桶是否存在，如果存在则打印已存在，否则打印错误
		exists, errBucketExists := MinioClient.BucketExists(ctx, bucketName)

		if errBucketExists == nil && exists {
			MinioClient.SetBucketPolicy(ctx, bucketName, s)
			log.Printf("已存在 %s\n", bucketName)
		} else {
			log.Fatalf("minio创建错误 %s\n", err)
		}
	} else {
		MinioClient.SetBucketPolicy(ctx, bucketName, s)
		log.Printf("成功创建 %s\n", bucketName)
	}
}

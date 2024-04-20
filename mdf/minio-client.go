package mdf

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)
     
func InitMinioClient() *minio.Client {
	endpoint := "127.0.0.1:80"
    accessKeyID := "lOGD8QU43Y3ElIB4"
    secretAccessKey := "YA9bL6RZTCDBSGp7Zh6RtmXFu477wjfb"
    useSSL := false

    minioClient, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
        Secure: useSSL,
    })

    if err != nil {
		log.Fatalln(err)
	}
   
    return minioClient
}


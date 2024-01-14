package pancham

import (
	"io"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type AwsPayload struct {
	AwsAccessKey string
	AwsSecretKey string
	AwsRegion    string
	Bucket       string
	Sess         *session.Session
}

func (awsProfile AwsPayload) SetupAws() (*session.Session, error) {
	log.Println("Setup Aws")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsProfile.AwsRegion),
		Credentials: credentials.NewStaticCredentials(awsProfile.AwsAccessKey, awsProfile.AwsSecretKey, ""),
	})
	if err != nil {
		log.Println("Error creating session ", err)
		return nil, err
	}

	return sess, nil
}

func (awsProfile AwsPayload) UploadFiletoS3(fileName string, srcFile io.ReadSeeker) (string, error) {
	sess := awsProfile.Sess
	svc := s3.New(sess)
	fileName = uuid.New().String() + "-" + fileName

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(awsProfile.Bucket),
		Key:         aws.String(fileName),
		Body:        srcFile,
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return fileName, nil
}

func (awsProfile AwsPayload) GetPreSignedUrl(fileName string) (string, error) {
	sess := awsProfile.Sess
	svc := s3.New(sess)
	resObj, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(awsProfile.Bucket),
		Key:    aws.String(fileName),
	})

	url, err := resObj.Presign(15 * time.Minute)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return url, nil
}

type CmdFileUrlDto struct {
	Key string `json:"key"`
	Url string `json:"url"`
}

func (awsProfile AwsPayload) GetListofKeys() ([]CmdFileUrlDto, error) {
	sess := awsProfile.Sess
	svc := s3.New(sess)

	result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(awsProfile.Bucket),
	})
	if err != nil {
		return nil, err
	}

	var keys []CmdFileUrlDto
	for _, obj := range result.Contents {
		if strings.Contains(*obj.Key, "cmd/") {
			url, _ := awsProfile.GetPreSignedUrl(*obj.Key)
			keys = append(keys, CmdFileUrlDto{
				Key: *obj.Key,
				Url: url,
			})
		}
	}

	return keys, nil
}

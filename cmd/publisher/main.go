package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	awsRegion := envOrDefault("GO_SQS_AWS_REGION", "ap-southeast-1")
	awsAccessKeyID := envOrDefault("GO_SQS_AWS_ACCESS_KEY_ID", "")
	awsSecret := envOrDefault("GO_SQS_AWS_SECRET", "")
	awsSessionToken := envOrDefault("GO_SQS_AWS_TOKEN", "")
	awsQueueURL := envOrDefault("GO_SQS_QUEUE_URL", "")

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecret, awsSessionToken),
		Region:      aws.String(awsRegion),
	}))

	var sqsService sqsiface.SQSAPI
	sqsService = sqs.New(sess)

	for {
		body := time.Now().Format(time.RFC3339)
		fmt.Println(body)
		out, err := sqsService.SendMessage(&sqs.SendMessageInput{
			QueueUrl:    aws.String(awsQueueURL),
			MessageBody: aws.String(body),
			MessageAttributes: map[string]*sqs.MessageAttributeValue{
				"string": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String("zz"),
				},
				"number": &sqs.MessageAttributeValue{
					DataType:    aws.String("Number"),
					StringValue: aws.String("1"),
				},
				"binary": &sqs.MessageAttributeValue{
					DataType:    aws.String("Binary"),
					BinaryValue: []byte(`{}`),
				},
			},
		})
		fmt.Println(out, err)
		time.Sleep(2 * time.Second)
	}
}

func envOrDefault(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

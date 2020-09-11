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

	session, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecret, awsSessionToken),
		Region:      aws.String(awsRegion),
	})
	if err != nil {
		log.Fatalln(err)
	}

	var sqsService sqsiface.SQSAPI
	sqsService = sqs.New(session)

	for {
		out, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:              aws.String(awsQueueURL),
			MessageAttributeNames: aws.StringSlice([]string{"string", "number", "binary"}),
			WaitTimeSeconds:       aws.Int64(10),
			VisibilityTimeout:     aws.Int64(0),
		})
		if err != nil {
			log.Fatalln(err)
		}

		time.Sleep(time.Second)
		for _, msg := range out.Messages {
			fmt.Printf("body: %s\n", *msg.Body)

			delOut, err := sqsService.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(awsQueueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				log.Fatalln(err)
			}
			_ = delOut
		}
	}
}

func envOrDefault(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

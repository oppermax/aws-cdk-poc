package main

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
)

var responseMsg = os.Getenv("HELLO_MESSAGE")

func Handler(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(jsonEntries)),
		Bucket: aws.String(),
		Key:    &fileKey,
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
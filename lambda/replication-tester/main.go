package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var responseMsg = os.Getenv("HELLO_MESSAGE")

func Handler(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {


	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
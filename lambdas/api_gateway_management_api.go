package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

func getAPIGatewayManagementAPIClient() *apigatewaymanagementapi.ApiGatewayManagementApi {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-3"),
	})
	if err != nil {
		log.Fatalln("Unable to create AWS session", err.Error())
	}

	return apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(os.Getenv("API_GATEWAY_WEBSOCKET_ENDPOINT")))
}

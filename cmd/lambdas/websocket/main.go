package main

import (
	"log"
	"net/http"

	"github.com/mvisonneau/approuvez/pkg/helpers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

func main() {
	lambda.Start(router)
}

func router(req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.RequestContext.RouteKey {
	case "$connect":
		return connect()
	case "$disconnect":
		return disconnect(req)
	case "get_connection_id":
		return handler(req, req.RequestContext.ConnectionID)
	case "keepalive":
		return handler(req, "")
	default:
		return handler(req, "unsupported request")
	}
}

func connect() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func disconnect(req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       req.RequestContext.RequestID,
	}, nil
}

func handler(req events.APIGatewayWebsocketProxyRequest, response string) (events.APIGatewayProxyResponse, error) {
	a := helpers.GetAPIGatewayManagementAPIClient()
	_, err := a.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(req.RequestContext.ConnectionID),
		Data:         []byte(response),
	})

	if err != nil {
		log.Println("ERROR", err.Error())
	}

	log.Println(req.RequestContext)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

package main

import (
	"log"
	"net/http"

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
		return getConnectionID(req)
	default:
		return unsupportedRequest(req)
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

func unsupportedRequest(req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	a := getAPIGatewayManagementAPIClient()
	_, err := a.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(req.RequestContext.ConnectionID),
		Data:         []byte("{\"error\":\"unsupported request\"}"),
	})

	if err != nil {
		log.Println("ERROR", err.Error())
	}

	log.Println(req.RequestContext)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func getConnectionID(req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	a := getAPIGatewayManagementAPIClient()
	_, err := a.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(req.RequestContext.ConnectionID),
		Data:         []byte(req.RequestContext.ConnectionID),
	})

	if err != nil {
		log.Println("ERROR", err.Error())
	}

	log.Println(req.RequestContext)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/slack-go/slack"
)

func main() {
	lambda.Start(slackCallback)
}

func slackCallback(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var p slack.InteractionCallback
	unescapedBody, _ := url.QueryUnescape(string(req.Body)[8:])

	if err := json.Unmarshal([]byte(unescapedBody), &p); err != nil {
		log.Println(err, unescapedBody)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("invalid JSON")
	}

	if len(p.ActionCallback.AttachmentActions) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("unable to fetch action name from payload")
	}

	a := getAPIGatewayManagementAPIClient()
	_, err := a.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(p.CallbackID),
		Data:         []byte(fmt.Sprintf("%s/%s", p.User.ID, p.ActionCallback.AttachmentActions[0].Name)),
	})

	if err != nil {
		log.Println("ERROR", err.Error())
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

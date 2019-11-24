package main

import (
	"bartenderAsFunction/dao"
	"bartenderAsFunction/model"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/satori/go.uuid"
)

var IotConnectionManager dao.IotConnectionInterface

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := []byte(request.Body)
	drunkClient := model.DrunkClient{}
	err := json.Unmarshal([]byte(body), &drunkClient)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}
	//assign an Id to the device when it does not have
	if drunkClient.IdClient == "" {
		uid := uuid.NewV4()
		drunkClient.IdClient = uid.String()
	}

	IotConnectionManager.RegisterDevice(&drunkClient)
	b, err := json.Marshal(drunkClient)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(b)}, nil
}

func main() {
	IotConnectionManager = dao.CreateIotConnection()
	lambda.Start(Handler)
}

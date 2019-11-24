package main

import (
	"bartenderAsFunction/dao"
	"bartenderAsFunction/model"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/satori/go.uuid"
)

var IotConnectionManager dao.IotConnectionInterface

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//TODO 1 get the body form the request
	var body []byte
	drunkClient := model.DrunkClient{}

	//TODO 2 Use the json unmarshall to parse the body request into the drunkClient variable
	err := json.Unmarshal([]byte(body), &drunkClient)

	//TODO 3 validate the error and retunr error 400. You can use events.APIGatewayProxyResponse
	/*if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}*/
	fmt.Println(err)
	if drunkClient.IdClient == "" {
		uid := uuid.NewV4()
		drunkClient.IdClient = uid.String()
	}
	//TODO 4 go to the IotConnectionManager in dao/iotDao.go and implement the method RegisterDevice
	IotConnectionManager.RegisterDevice(&drunkClient)
	//assign an Id to the device when it does not have
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

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type User struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	JoinedDate      string `json:"joineddate"`
	SamaritanPoints string `json:"samaritanpoints"`
	Location        string `json:"location"`
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return fetch(req)
	case "POST":
		return insert(req)
	default:
		return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed,
			Body: http.StatusText(http.StatusMethodNotAllowed)}, nil
	}
}

func fetch(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	issues, err := getItems()
	if err != nil {
		//See if we can pass err instead

		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadGateway,
			Body: err.Error()}, nil
	}
	issues_json, err := json.Marshal(issues)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError)}, nil
	}
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		"Access-Control-Allow-Methods": "OPTIONS,POST,GET"}
	return events.APIGatewayProxyResponse{
		Body:       string(issues_json),
		Headers:    headers,
		StatusCode: 201,
	}, nil
}

func insert(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.Headers["content-type"] != "application/json" && request.Headers["Content-Type"] != "application/json" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotAcceptable,
			Body: http.StatusText(http.StatusNotAcceptable)}, nil
	}
	user := new(User)
	err := json.Unmarshal([]byte(request.Body), user)
	user.ID = uuid.New().String()
	user.JoinedDate = time.Now().Local().String()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest,
			Body: http.StatusText(http.StatusBadRequest)}, nil
	}
	err = putItem(user)
	if err != nil {
		//See if we can pass err instead
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadGateway,
			Body: err.Error()}, nil
	}
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		"Access-Control-Allow-Methods": "OPTIONS,POST,GET"}
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Successfully added the user"),
		Headers:    headers,
		StatusCode: 201,
	}, nil
}

func main() {
	lambda.Start(router)
}

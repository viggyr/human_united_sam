package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

type Test struct {
	ID   string `json:"id"`
	Text string `json:"text"`
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
	fmt.Println(issues)
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

	// issue := Issue{
	// 	ID:       "1234",
	// 	Created:  "08/30/2020",
	// 	Title:    "Dummy Issue",
	// 	Body:     "I got a real problem here. Please help me",
	// 	Private:  1,
	// 	User:     "Viggy",
	// 	Location: "Bangalore",
	// }
	if request.Headers["content-type"] != "application/json" && request.Headers["Content-Type"] != "application/json" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotAcceptable,
			Body: http.StatusText(http.StatusNotAcceptable)}, nil
	}
	issue := new(Test)
	fmt.Println(request.Body)
	err := json.Unmarshal([]byte(request.Body), issue)
	issue.ID = uuid.New().String()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest,
			Body: http.StatusText(http.StatusBadRequest)}, nil
	}
	fmt.Println(issue)
	err = putItem(issue)
	if err != nil {
		//See if we can pass err instead

		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadGateway,
			Body: err.Error()}, nil
	}
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		"Access-Control-Allow-Methods": "OPTIONS,POST,GET"}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Successfully stored the entry"),
		Headers:    headers,
		StatusCode: 201,
	}, nil
}

func main() {
	env := os.Getenv("AWSENV")
	dbEndpoint := os.Getenv("DBENDPOINT")
	createDBConnection(env, dbEndpoint)
	lambda.Start(router)
}

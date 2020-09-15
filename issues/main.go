package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

type Issue struct {
	ID       string   `json:"id"`
	Created  string   `json:"created"`
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Private  int      `json:"private"`
	UserID   string   `json:"userid"`
	UserName string   `json:"username"`
	Location string   `json:"location"`
	Personal int      `json:"personal"`
	Helpers  []string `json:"helpers"`
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		"Access-Control-Allow-Methods": "OPTIONS,POST,GET"}
	switch req.HTTPMethod {
	case "GET":
		return fetch(req)
	case "POST":
		return insert(req)
	default:
		return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed,
			Headers: headers,
			Body:    http.StatusText(http.StatusMethodNotAllowed)}, nil
	}
}

func fetch(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		"Access-Control-Allow-Methods": "OPTIONS,POST,GET"}
	if issueID, ok := request.PathParameters["issueId"]; ok {
		fmt.Println(issueID)
		issue, err := getIssueById(issueID)
		issue_json, err := json.Marshal(issue)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    headers,
				Body:       http.StatusText(http.StatusInternalServerError)}, nil
		}
		return events.APIGatewayProxyResponse{
			Body:       string(issue_json),
			Headers:    headers,
			StatusCode: 201,
		}, nil
	} else {

		issues, err := getItems()

		if err != nil {
			//See if we can pass err instead
			fmt.Printf("Failed to fetch data %s", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadGateway,
				Headers: headers,
				Body:    err.Error()}, nil
		}
		fmt.Println(issues)
		issues_json, err := json.Marshal(issues)
		if err != nil {

			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    headers,
				Body:       http.StatusText(http.StatusInternalServerError)}, nil
		}
		fmt.Println("Success")
		return events.APIGatewayProxyResponse{
			Body:       string(issues_json),
			Headers:    headers,
			StatusCode: 201,
		}, nil
	}

}

func insert(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		"Access-Control-Allow-Methods": "OPTIONS,POST,GET"}
	if request.Headers["content-type"] != "application/json" && request.Headers["Content-Type"] != "application/json" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotAcceptable,
			Headers: headers,
			Body:    http.StatusText(http.StatusNotAcceptable)}, nil
	}
	issue := new(Issue)
	issue.Personal = 1
	fmt.Println(request.Body)
	err := json.Unmarshal([]byte(request.Body), issue)
	issue.ID = uuid.New().String()
	issue.Created = time.Now().Local().String()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest,
			Headers: headers,
			Body:    http.StatusText(http.StatusBadRequest)}, nil
	}
	fmt.Println(issue)
	err = putItem(issue)
	if err != nil {
		//See if we can pass err instead

		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadGateway,
			Headers: headers,
			Body:    err.Error()}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Successfully stored the entry"),
		Headers:    headers,
		StatusCode: 201,
	}, nil
}

func main() {
	//get parameters here from environment
	env := os.Getenv("AWSENV")
	dbEndpoint := os.Getenv("DBENDPOINT")
	createDBConnection(env, dbEndpoint)
	lambda.Start(router)
}

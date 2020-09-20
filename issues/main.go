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

type CommentsRequest struct {
	UserID   string `json:"userid"`
	UserName string `json:"username"`
	Comment  string `json:"comment"`
}

type HelpersRequest struct {
	UserID   string `json:"userid"`
	UserName string `json:"username"`
}

type Issue struct {
	ID        string            `json:"id"`
	Created   string            `json:"created"`
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	Private   int               `json:"private"`
	UserID    string            `json:"userid"`
	UserName  string            `json:"username"`
	Location  string            `json:"location"`
	Personal  int               `json:"personal"`
	Helpers   map[string]string `json:"helpers"`
	Comments  []CommentsRequest `json:comments`
	StatusMsg string            `json:"statusmsg"`
}

type StatusRequest struct {
	StatusMsg string `json:"statusmsg"`
}

func getHeaders() map[string]string {
	return map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		"Access-Control-Allow-Methods": "OPTIONS,POST,GET"}
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return fetch(req)
	case "POST":
		return insert(req)
	case "PUT":
		return update(req)
	case "OPTIONS":
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    getHeaders()}, nil
	default:
		return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed,
			Headers: getHeaders(),
			Body:    http.StatusText(http.StatusMethodNotAllowed)}, nil
	}
}
func fetch(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if issueID, ok := request.PathParameters["issueId"]; ok {
		issue, err := getIssueById(issueID)
		issue_json, err := json.Marshal(issue)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    getHeaders(),
				Body:       http.StatusText(http.StatusInternalServerError)}, nil
		}
		return events.APIGatewayProxyResponse{
			Body:       string(issue_json),
			Headers:    getHeaders(),
			StatusCode: 201,
		}, nil
	} else {

		issues, err := getItems()

		if err != nil {
			//See if we can pass err instead
			fmt.Printf("Failed to fetch data %s", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadGateway,
				Headers: getHeaders(),
				Body:    err.Error()}, nil
		}
		fmt.Println(issues)
		issues_json, err := json.Marshal(issues)
		if err != nil {

			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    getHeaders(),
				Body:       http.StatusText(http.StatusInternalServerError)}, nil
		}
		fmt.Println("Success")
		return events.APIGatewayProxyResponse{
			Body:       string(issues_json),
			Headers:    getHeaders(),
			StatusCode: 201,
		}, nil
	}

}

func insert(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.Headers["content-type"] != "application/json" && request.Headers["Content-Type"] != "application/json" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotAcceptable,
			Headers: getHeaders(),
			Body:    http.StatusText(http.StatusNotAcceptable)}, nil
	}
	issue := new(Issue)
	issue.Personal = 1
	issue.StatusMsg = "Need Help"
	err := json.Unmarshal([]byte(request.Body), issue)
	issue.ID = uuid.New().String()
	issue.Created = time.Now().Local().String()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest,
			Headers: getHeaders(),
			Body:    http.StatusText(http.StatusBadRequest)}, nil
	}
	err = putItem(issue)
	if err != nil {
		//See if we can pass err instead

		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadGateway,
			Headers: getHeaders(),
			Body:    err.Error()}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Successfully stored the entry"),
		Headers:    getHeaders(),
		StatusCode: 201,
	}, nil
}

func update(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	issueId := request.PathParameters["issueId"]
	field := request.PathParameters["field"]
	switch field {
	case "comment":
		commentReq := new(CommentsRequest)
		err := json.Unmarshal([]byte(request.Body), commentReq)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest,
				Headers: getHeaders(),
				Body:    http.StatusText(http.StatusBadRequest)}, nil
		}
		err = updateCommentsForIssue(issueId, commentReq)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    getHeaders(),
				Body:       err.Error()}, nil
		}

	case "help":
		helperReq := new(HelpersRequest)
		err := json.Unmarshal([]byte(request.Body), helperReq)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest,
				Headers: getHeaders(),
				Body:    http.StatusText(http.StatusBadRequest)}, nil
		}
		err = updateHelpersForIssue(issueId, helperReq)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    getHeaders(),
				Body:       err.Error()}, nil
		}
	case "status":
		statusReq := new(StatusRequest)
		err := json.Unmarshal([]byte(request.Body), statusReq)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest,
				Headers: getHeaders(),
				Body:    http.StatusText(http.StatusBadRequest)}, nil
		}
		err = updateStatusForIssue(issueId, statusReq)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    getHeaders(),
				Body:       "Failed to update status for issue"}, nil
		}
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    getHeaders(),
			Body:       "Invalid request parameters"}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Successfully updated the Issue"),
		Headers:    getHeaders(),
		StatusCode: 201,
	}, nil
}

//Add put for discussion and status
func main() {
	//get parameters here from environment
	env := os.Getenv("AWSENV")
	dbEndpoint := os.Getenv("DBENDPOINT")
	createDBConnection(env, dbEndpoint)
	lambda.Start(router)
}

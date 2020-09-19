package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type User struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Email             string   `json:"email"`
	ProfileImageUrl   string   `json:"profileimageurl"`
	JoinedDate        string   `json:"joineddate"`
	LastLogin         string   `json:"lastlogin"`
	SamaritanPoints   int      `json:"samaritanpoints"`
	UserIssues        []string `json:"userissues"`
	UserHelps         []string `json:"userhelps"`
	UserInterests     []string `json:userinterests`
	UsersCurrentHelps []string `json:usercurrenthelps`
}

type UserRequest struct {
	Action   string `json:"action"`
	Scenario string `json:"scenario"`
	IssueID  string `json:"issue_id"`
}

func getHeaders() map[string]string {
	return map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		"Access-Control-Allow-Methods": "OPTIONS,POST,GET"}
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userId := req.PathParameters["userId"]
	fmt.Printf("Path parameter user id :%s", userId)
	switch req.HTTPMethod {
	case "GET":
		return fetch(req, userId)
	case "PUT":
		return insert(req, userId)
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

func fetch(request events.APIGatewayProxyRequest, userId string) (events.APIGatewayProxyResponse, error) {
	userInfo, err := getUserById(userId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    getHeaders(),
			Body:       err.Error()}, nil
	}
	user_json, err := json.Marshal(userInfo)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    getHeaders(),
			Body:       http.StatusText(http.StatusInternalServerError)}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(user_json),
		Headers:    getHeaders(),
		StatusCode: 200,
	}, nil
}

func insert(request events.APIGatewayProxyRequest, userId string) (events.APIGatewayProxyResponse, error) {

	if request.Headers["content-type"] != "application/json" && request.Headers["Content-Type"] != "application/json" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotAcceptable,
			Headers: getHeaders(),
			Body:    http.StatusText(http.StatusNotAcceptable)}, nil
	}
	userRequest := new(UserRequest)
	err := json.Unmarshal([]byte(request.Body), userRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest,
			Headers: getHeaders(),
			Body:    http.StatusText(http.StatusBadRequest)}, nil
	}
	err = updateUser(userId, userRequest)
	if err != nil {

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    getHeaders(),
			Body:       err.Error()}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Successfully updated User"),
		Headers:    getHeaders(),
		StatusCode: 201,
	}, nil
}

func main() {
	env := os.Getenv("AWSENV")
	dbEndpoint := os.Getenv("DBENDPOINT")
	createDBConnection(env, dbEndpoint)
	lambda.Start(router)
}

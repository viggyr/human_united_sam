package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type User struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	ProfileImageUrl string   `json: "imageurl"`
	Location        string   `json:"location"`
	JoinedDate      string   `json:"joineddate"`
	LastLogin       string   `json: "lastlogin"`
	SamaritanPoints string   `json:"samaritanpoints"`
	UserIssues      []string `json:"userissues"`
	UserHelps       []string `json: "userhelps"`
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
	users, err := getUsers()
	if err != nil {
		//See if we can pass err instead

		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadGateway,
			Body: err.Error()}, nil
	}
	users_json, err := json.Marshal(users)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError)}, nil
	}
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		"Access-Control-Allow-Methods": "OPTIONS,POST,GET"}
	return events.APIGatewayProxyResponse{
		Body:       string(users_json),
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

	// Inserting a new user would first get all users , check if user exists.
	// if user exists , do partial update with the fields passed in request body.
	// If user does not exists, create new entry with default values where applicable.
	//users, err := getItems()

	err := json.Unmarshal([]byte(request.Body), user)
	user.ID = uuid.New().String()
	user.JoinedDate = time.Now().Local().String()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest,
			Body: http.StatusText(http.StatusBadRequest)}, nil
	}
	err = putUser(user)
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
	env := os.Getenv("AWSENV")
	dbEndpoint := os.Getenv("DBENDPOINT")
	createDBConnection(env, dbEndpoint)
	lambda.Start(router)
}

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
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	ProfileImageUrl string `json:"imageurl"`
	JoinedDate      string `json:"joineddate"`
	LastLogin       string `json:"lastlogin"`
	SamaritanPoints int    `json:"samaritanpoints"`
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
	default:
		return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed,
			Headers: getHeaders(),
			Body:    http.StatusText(http.StatusMethodNotAllowed)}, nil
	}
}

func fetch(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	users, err := getUsers()
	if err != nil {
		//See if we can pass err instead

		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadGateway,
			Headers: getHeaders(),
			Body:    err.Error()}, nil
	}
	users_json, err := json.Marshal(users)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
			Headers:    getHeaders()}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(users_json),
		Headers:    getHeaders(),
		StatusCode: 201,
	}, nil
}

func insert(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if request.Headers["content-type"] != "application/json" && request.Headers["Content-Type"] != "application/json" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotAcceptable,
			Headers:    getHeaders(),
			Body:       http.StatusText(http.StatusNotAcceptable)}, nil
	}
	user := new(User)
	err := json.Unmarshal([]byte(request.Body), user)
	currTime := time.Now().Local().String()
	existingUser, err := checkIfUserExists(user.Email)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    getHeaders(),
			Body:       err.Error()}, nil
	}
	if !existingUser {
		user.ID = uuid.New().String()
		user.JoinedDate = currTime
		user.LastLogin = currTime
		// default samaratian points - 10
		user.SamaritanPoints = 10
		err = putUser(user)
	} else {
		err = updateUserLastLogin(currTime, user.Email)
	}

	if err != nil {
		//See if we can pass err instead
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadGateway,
			Headers: getHeaders(),
			Body:    err.Error()}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Successfully added the user"),
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

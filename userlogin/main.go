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

type User struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	ProfileImageUrl string `json:"imageurl"`
	JoinedDate      string `json:"joineddate"`
	LastLogin       string `json:"lastlogin"`
	SamaritanPoints int    `json:"samaritanpoints"`
}

type LoginResponse struct {
	UserID          string
	SamaritanPoints int
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
	loginResponse := new(LoginResponse)
	err := json.Unmarshal([]byte(request.Body), user)
	currTime := time.Now().Local().String()
	existingUser, err := checkIfUserExists(user.Email)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    getHeaders(),
			Body:       fmt.Sprintf("Failed to check if user exists")}, nil
	}

	if existingUser != nil {
		err = updateUserLastLogin(currTime, existingUser.ID)
		if err != nil {
			//See if we can pass err instead
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    getHeaders(),
				Body:       err.Error()}, nil
		}
		loginResponse.UserID = existingUser.ID
		loginResponse.SamaritanPoints = existingUser.SamaritanPoints
		loginJson, err := json.Marshal(loginResponse)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       http.StatusText(http.StatusInternalServerError),
				Headers:    getHeaders()}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 201,
			Headers:    getHeaders(),
			Body:       string(loginJson)}, nil
	}

	user.ID = uuid.New().String()
	user.JoinedDate = currTime
	user.LastLogin = currTime
	// default samaratian points - 10
	user.SamaritanPoints = 10
	err = putUser(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
			Headers:    getHeaders()}, nil
	}
	loginResponse.UserID = user.ID
	loginResponse.SamaritanPoints = user.SamaritanPoints
	loginJson, err := json.Marshal(loginResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
			Headers:    getHeaders()}, nil
	}
	return events.APIGatewayProxyResponse{
		Body:       string(loginJson),
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

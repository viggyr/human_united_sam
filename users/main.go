package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

type Issue struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	StatusMsg string `json:"statusmsg"`
}

type User struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	ProfileImageUrl string   `json:"profileimageurl"`
	JoinedDate      string   `json:"joineddate"`
	LastLogin       string   `json:"lastlogin"`
	SamaritanPoints int      `json:"samaritanpoints"`
	UserIssues      []*Issue `json:"userissues"`
	UserHelps       []*Issue `json:"userhelps"`
	//UserInterests     []Issue `json:userinterests`
	//UsersCurrentHelps []Issue `json:usercurrenthelps`
}

type Post struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PostTime    string `json:"posttime"`
	UserId      string `json:"userid"`
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
	if strings.HasPrefix(req.Path, "/users") {
		userId := req.PathParameters["userId"]
		fmt.Printf("Path parameter user id :%s", userId)
		switch req.HTTPMethod {
		case "GET":
			return fetch(req, userId)
		case "PUT":
			return insert(req, userId)
		case "POST":
			return insertPost(req)
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
	if strings.HasPrefix(req.Path, "/userPosts") {
		userId := req.PathParameters["userId"]
		fmt.Printf("Path parameter user id :%s", userId)
		switch req.HTTPMethod {
		case "GET":
			return fetchPostsByUserId(req, userId)
		case "POST":
			return insertPost(req)
		default:
			return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed,
				Headers: getHeaders(),
				Body:    http.StatusText(http.StatusMethodNotAllowed)}, nil
		}
	}
	if strings.HasPrefix(req.Path, "/posts") {
		switch req.HTTPMethod {
		case "GET":
			return fetchAllPosts(req)
		default:
			return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed,
				Headers: getHeaders(),
				Body:    http.StatusText(http.StatusMethodNotAllowed)}, nil
		}
	}
	return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed,
		Headers: getHeaders(),
		Body:    http.StatusText(http.StatusMethodNotAllowed)}, nil
}

func fetch(request events.APIGatewayProxyRequest, userId string) (events.APIGatewayProxyResponse, error) {
	userInfo, err := getUserById(userId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    getHeaders(),
			Body:       err.Error()}, nil
	}

	// one filter call to get issues created by user
	userIssues, err := getIssuesCreatedByUser(userId)
	// one filter call to get issues helped by user
	helpedIssues, err := getIssuesHelpedByUser(userId, userInfo.Name)
	userInfo.UserIssues = userIssues
	userInfo.UserHelps = helpedIssues
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
	//err = updateUser(userId, userRequest)
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

func insertPost(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		"Access-Control-Allow-Methods": "OPTIONS,POST,GET"}
	if request.Headers["content-type"] != "application/json" && request.Headers["Content-Type"] != "application/json" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotAcceptable,
			Headers: headers,
			Body:    http.StatusText(http.StatusNotAcceptable)}, nil
	}
	post := new(Post)
	err := json.Unmarshal([]byte(request.Body), post)
	post.ID = uuid.New().String()
	post.PostTime = time.Now().Local().String()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest,
			Headers: headers,
			Body:    http.StatusText(http.StatusBadRequest)}, nil
	}
	err = addPost(post)
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

func fetchPostsByUserId(request events.APIGatewayProxyRequest, userId string) (events.APIGatewayProxyResponse, error) {
	userInfo, err := getPostsByUserId(userId)
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

func fetchAllPosts(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userInfo, err := getAllPosts()
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

func main() {
	env := os.Getenv("AWSENV")
	dbEndpoint := os.Getenv("DBENDPOINT")
	createDBConnection(env, dbEndpoint)
	lambda.Start(router)
}

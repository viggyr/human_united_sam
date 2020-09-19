package main

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var db *dynamodb.DynamoDB

//const usersTable = "huManUnited-UsersTable-16HJ59LOVEINZ"
var usersTable = "users"
var postsTable = "posts"

func createDBConnection(env string, endpoint string) {
	if env == "AWS_SAM_LOCAL" {
		sess, err := session.NewSession(&aws.Config{
			Region:   aws.String("ap-south-1"),
			Endpoint: aws.String(endpoint)})
		if err != nil {
			fmt.Println("Failed to create dynamodb session")

		}
		db = dynamodb.New(sess)
	} else {
		db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("ap-south-1"))
	}
}

func getUsers() ([]*User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(usersTable),
	}

	result, err := db.Scan(input)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, nil
	}
	users := make([]*User, 0)
	for _, i := range result.Items {
		user := new(User)

		err = dynamodbattribute.UnmarshalMap(i, &user)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func addPost(post *Post) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(postsTable),
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(post.ID),
			},
			"Title": {
				S: aws.String(post.Title),
			},
			"Description": {
				S: aws.String(post.Description),
			},
			"PostTime": {
				S: aws.String(post.PostTime),
			},
			"UserId": {
				S: aws.String(post.UserId),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}

func getPostsByUserId(userId string) ([]*Post, error) {
	filt := expression.Name("UserId").Equal(expression.Value(userId))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		fmt.Println("Failed to build filter by email expression ")
		fmt.Println(err.Error())
	}
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(postsTable),
	}
	result, err := db.Scan(input)
	fmt.Println("Kiran: 3 %s", result)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, nil
	}
	posts := make([]*Post, 0)
	for _, i := range result.Items {
		post := new(Post)
		err = dynamodbattribute.UnmarshalMap(i, &post)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func getAllPosts() ([]*Post, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(postsTable),
	}
	result, err := db.Scan(input)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, nil
	}
	posts := make([]*Post, 0)
	for _, i := range result.Items {
		post := new(Post)
		err = dynamodbattribute.UnmarshalMap(i, &post)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func deleteIssueForUser(userID string, issueID string) error {
	return nil
}
func deleteHelpForUser(userID string, issueID string) error {
	return nil
}
func deleteInterestForUser(userID string, issueID string) error {
	return nil
}

func deleteListItemForUser(userID string, issueID string, scenario string) error {
	switch scenario {
	case "issue":
		return deleteIssueForUser(userID, issueID)
	case "help":
		return deleteHelpForUser(userID, issueID)
	case "interest":
		return deleteInterestForUser(userID, issueID)
	default:
		return errors.New("Invalid Scenario")
	}
}

func addIssueForUser(userId string, issueId string) error {
	fmt.Printf("User %s has posted issue with issue ID %s", userId, issueId)
	issueIDList := []string{issueId}
	issueAVs, err := dynamodbattribute.MarshalList(issueIDList)
	if err != nil {
		fmt.Printf("Could not Marshal issue id list %s", err.Error())
		return err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":i": {
				L: issueAVs,
			},
			":empty_list": {
				L: []*dynamodb.AttributeValue{},
			},
		},
		TableName: aws.String(usersTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(userId),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET UserIssues = list_append(if_not_exists(UserIssues, :empty_list),:i)"),
	}
	_, err = db.UpdateItem(input)
	return err
}

func addHelpForUser(userId string, issueId string) error {
	fmt.Printf("User %s is providing help for issue ID %s", userId, issueId)
	issueIDList := []string{issueId}
	issueAVs, err := dynamodbattribute.MarshalList(issueIDList)
	if err != nil {
		fmt.Printf("Could not Marshal issue id list %s", err.Error())
		return err
	}
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":i": {
				L: issueAVs,
			},
			":empty_list": {
				L: []*dynamodb.AttributeValue{},
			},
			":s": {
				N: aws.String("5"),
			},
		},
		TableName: aws.String(usersTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(userId),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET UserHelps = list_append(if_not_exists(UserHelps, :empty_list),:i), SamaritanPoints = SamaritanPoints + :s"),
	}
	_, err = db.UpdateItem(input)
	return err
}

func addInterestForUser(userId string, issueId string) error {
	fmt.Printf("User %s interested in issue ID %s", userId, issueId)
	issueIDList := []string{issueId}
	issueAVs, err := dynamodbattribute.MarshalList(issueIDList)
	if err != nil {
		fmt.Printf("Could not Marshal issue id list %s", err.Error())
		return err
	}
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":i": {
				L: issueAVs,
			},
			":empty_list": {
				L: []*dynamodb.AttributeValue{},
			},
		},
		TableName: aws.String(usersTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(userId),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET UserInterests = list_append(if_not_exists(UserInterests, :empty_list),:i)"),
	}
	_, err = db.UpdateItem(input)
	return err
}

func addListItemForUser(userID string, issueID string, scenario string) error {
	switch scenario {
	case "issue":
		return addIssueForUser(userID, issueID)
	case "help":
		return addHelpForUser(userID, issueID)
	case "interest":
		return addInterestForUser(userID, issueID)
	default:
		return errors.New("Invalid Scenario passed!!")
	}
}

func updateUser(userId string, userReq *UserRequest) error {
	switch userReq.Action {
	case "DELETE":
		return deleteListItemForUser(userId, userReq.IssueID, userReq.Scenario)
	case "ADD":
		return addListItemForUser(userId, userReq.IssueID, userReq.Scenario)
	default:
		return nil
	}
}

func getUserById(userId string) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(userId),
			},
		},
		TableName: aws.String(usersTable),
	}
	result, err := db.GetItem(input)
	if err != nil {
		fmt.Printf("Failed to get Item from table %s for %s", usersTable, userId)
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	user := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

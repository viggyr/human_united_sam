package main

import (
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
var issuesTable = "issues"

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

func getIssuesCreatedByUser(userId string) ([]*Issue, error) {
	// get userid and filter by userid
	// return issue data with projection
	filt := expression.Name("UserID").Equal(expression.Value(userId))
	proj := expression.NamesList(expression.Name("Id"), expression.Name("Title"), expression.Name("StatusMsg"))
	expr, err := expression.NewBuilder().WithProjection(proj).WithFilter(filt).Build()
	if err != nil {
		fmt.Println("Failed to build filter by userId expression ")
		fmt.Println(err.Error())
	}
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(issuesTable),
	}
	result, err := db.Scan(input)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, nil
	}
	issues := make([]*Issue, 0)
	for _, i := range result.Items {
		issue := new(Issue)
		err = dynamodbattribute.UnmarshalMap(i, &issue)
		if err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}
	return issues, nil
}

func getIssuesHelpedByUser(userId string, userName string) ([]*Issue, error) {
	// check for userid and username in issues table
	//  return issue data with projection
	filt := expression.Name("Helpers." + userId).Contains(userName)
	proj := expression.NamesList(expression.Name("Id"), expression.Name("Title"), expression.Name("StatusMsg"))
	expr, err := expression.NewBuilder().WithProjection(proj).WithFilter(filt).Build()
	if err != nil {
		fmt.Println("Failed to build filter by Helpers expression ")
		fmt.Println(err.Error())
	}
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(issuesTable),
	}
	result, err := db.Scan(input)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, nil
	}
	issues := make([]*Issue, 0)
	for _, i := range result.Items {
		issue := new(Issue)
		err = dynamodbattribute.UnmarshalMap(i, &issue)
		if err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}
	return issues, nil
	return nil, nil
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

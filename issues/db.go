package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db *dynamodb.DynamoDB

var IssuesTable = os.Getenv("ISSUESTABLE")

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

func getItems() ([]*Issue, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(IssuesTable),
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

func putItem(issue *Issue) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(IssuesTable),
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(issue.ID),
			},
			"Title": {
				S: aws.String(issue.Title),
			},
			"Created": {
				S: aws.String(issue.Created),
			},
			"Body": {
				S: aws.String(issue.Body),
			},
			"Private": {
				N: aws.String(strconv.Itoa(issue.Private)),
			},
			"Location": {
				S: aws.String(issue.Location),
			},
			"UserID": {
				S: aws.String(issue.UserID),
			},
			"Personal": {
				N: aws.String(strconv.Itoa(issue.Personal)),
			},
			"Helpers": {
				SS: aws.StringSlice(issue.Helpers),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}

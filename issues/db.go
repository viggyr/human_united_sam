package main

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("ap-south-1"))

func getItems() ([]*Issue, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("huManUnited-IssuesTable-HCOFLC4PHPHQ"),
	}

	result, err := db.Scan(input)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, nil
	}
	fmt.Println(result)
	issues := make([]*Issue, 0)
	for _, i := range result.Items {
		issue := new(Issue)
		fmt.Println(i)
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
		TableName: aws.String("huManUnited-IssuesTable-HCOFLC4PHPHQ"),
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
			"User": {
				S: aws.String(issue.User),
			},
			"Personal": {
				N: aws.String(strconv.Itoa(issue.Personal)),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}

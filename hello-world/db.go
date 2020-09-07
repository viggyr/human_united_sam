package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db *dynamodb.DynamoDB

const testTable = "TestTable"

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

func getItems() ([]*Test, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("TestTable"),
	}

	result, err := db.Scan(input)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, nil
	}
	fmt.Println(result)
	issues := make([]*Test, 0)
	for _, i := range result.Items {
		issue := new(Test)
		fmt.Println(i)
		err = dynamodbattribute.UnmarshalMap(i, &issue)

		if err != nil {
			return nil, err
		}

		issues = append(issues, issue)
	}
	return issues, nil
}

func putItem(issue *Test) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("TestTable"),
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(issue.ID),
			},
			"Text": {
				S: aws.String(issue.Text),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}

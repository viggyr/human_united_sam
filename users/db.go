package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("ap-south-1"))

func getItems() ([]*User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("huManUnited-UsersTable-16HJ59LOVEINZ"),
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

func putItem(user *User) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("huManUnited-UsersTable-16HJ59LOVEINZ"),
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(user.ID),
			},
			"Name": {
				S: aws.String(user.Name),
			},
			"JoinedDate": {
				S: aws.String(user.JoinedDate),
			},
			"SamaritanPoints": {
				S: aws.String(user.SamaritanPoints),
			},
			"Location": {
				S: aws.String(user.Location),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}

package main

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var db *dynamodb.DynamoDB

//const usersTable = "huManUnited-UsersTable-16HJ59LOVEINZ"
var usersTable = "users"

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

func putUser(user *User) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(usersTable),
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(user.ID),
			},
			"Name": {
				S: aws.String(user.Name),
			},
			"Email": {
				S: aws.String(user.Email),
			},
			"JoinedDate": {
				S: aws.String(user.JoinedDate),
			},
			"SamaritanPoints": {
				N: aws.String(strconv.Itoa(user.SamaritanPoints)),
			},
			"ProfileImageUrl": {
				S: aws.String(user.ProfileImageUrl),
			},
			"LastLogin": {
				S: aws.String(user.LastLogin),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}

func updateUserLastLogin(loginTime string, userID string) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":l": {
				S: aws.String(loginTime),
			},
		},
		TableName: aws.String(usersTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(userID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set LastLogin = :l"),
	}

	_, err := db.UpdateItem(input)
	return err
}

func checkIfUserExists(usermail string) (*User, error) {
	filt := expression.Name("Email").Equal(expression.Value(usermail))
	proj := expression.NamesList(expression.Name("Id"), expression.Name("SamaritanPoints"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		fmt.Println("Failed to build filter by email expression")
		fmt.Println(err.Error())
	}
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(usersTable),
	}

	result, err := db.Scan(input)
	if err != nil {
		fmt.Printf("Failed to scan the table %s using filter expression", usersTable)
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, nil
	}
	user := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

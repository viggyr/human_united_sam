package main

import (
	"fmt"
	"os"
	"reflect"
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

func updateCommentsForIssue(issueId string, commentData *CommentsRequest) error {
	fmt.Printf("User %s is provided comment for issue ID %s", commentData.UserID, issueId)
	commentsList := []*CommentsRequest{commentData}
	commentAVs, err := dynamodbattribute.MarshalList(commentsList)
	if err != nil {
		fmt.Printf("Could not Marshal comments list %s", err.Error())
		return err
	}
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {
				L: commentAVs,
			},
			":empty_list": {
				L: []*dynamodb.AttributeValue{},
			},
		},
		TableName: aws.String(IssuesTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(issueId),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET Comments = list_append(if_not_exists(Comments, :empty_list),:c)"),
	}
	_, err = db.UpdateItem(input)
	return err
}

func updateStatusForIssue(issueId string, statusData *StatusRequest) error {
	fmt.Printf("Status changed to %s for issue ID %s", statusData.StatusMsg, issueId)
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				S: aws.String(statusData.StatusMsg),
			},
		},
		TableName: aws.String(IssuesTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(issueId),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set StatusMsg = :s"),
	}

	_, err := db.UpdateItem(input)
	return err

}

func updateHelpersForIssue(issueId string, helpersData *HelpersRequest) error {
	fmt.Printf("User %s is providing help for issue ID %s", helpersData.UserID, issueId)
	helperIDList := []string{helpersData.UserID}
	helperAVs, err := dynamodbattribute.MarshalList(helperIDList)
	if err != nil {
		fmt.Printf("Could not Marshal user ids list %s", err.Error())
		return err
	}
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":h": {
				L: helperAVs,
			},
			":empty_list": {
				L: []*dynamodb.AttributeValue{},
			},
		},
		TableName: aws.String(IssuesTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(issueId),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET Helpers = list_append(if_not_exists(Helpers, :empty_list),:h)"),
	}
	_, err = db.UpdateItem(input)
	return err
}

func getIssueById(issueID string) (*Issue, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(issueID),
			},
		},
		TableName: aws.String(IssuesTable),
	}
	result, err := db.GetItem(input)
	if err != nil {
		fmt.Printf("Failed to get Item from table %s for %s", IssuesTable, issueID)
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	issue := new(Issue)
	err = dynamodbattribute.UnmarshalMap(result.Item, &issue)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(*issue, Issue{}) {
		return nil, nil
	}
	return issue, nil
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
			"UserName": {
				S: aws.String(issue.UserName),
			},
			"Personal": {
				N: aws.String(strconv.Itoa(issue.Personal)),
			},
			"StatusMsg": {
				S: aws.String(issue.StatusMsg),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}

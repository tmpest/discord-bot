package discordbot

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/tmpest/discord-bot/pkg/server"
)

var discordTokenTableName string = "discord-token"

const tokenInfoKey string = "tokenInfo"

// SetTokenInfo stores token auth information in dynamo DB for a given accountID
func SetTokenInfo(accountID *string, tokenInformation []byte) error {
	dynamoDBClient, error := newDynamoClient()
	if error != nil {
		return error
	}

	accountIDAttributeValue := &dynamodb.AttributeValue{S: accountID}
	tokenInformationAttributeValue := &dynamodb.AttributeValue{B: tokenInformation}

	input := dynamodb.PutItemInput{
		Item:      map[string]*dynamodb.AttributeValue{"accountID": accountIDAttributeValue, tokenInfoKey: tokenInformationAttributeValue},
		TableName: &discordTokenTableName,
	}
	fmt.Printf("PutItemInput: \n%+v\n", input.GoString())

	request, _ := dynamoDBClient.PutItemRequest(&input)
	error = request.Send()

	if error != nil {
		fmt.Println("There was a problem writing to Dynamo DB", error)
		return error
	}
	return error
}

// GetTokenInfo fetches the token information for a given accountID
func GetTokenInfo(accountID *string) (*server.TokenInformation, error) {
	dynamoDBClient, error := newDynamoClient()
	if error != nil {
		return nil, error
	}

	accountIDAttributeValue := &dynamodb.AttributeValue{S: accountID}

	input := dynamodb.GetItemInput{
		Key:       map[string]*dynamodb.AttributeValue{"accountID": accountIDAttributeValue},
		TableName: &discordTokenTableName,
	}

	fmt.Printf("PutItemInput: \n%+v\n", input.GoString())

	request, response := dynamoDBClient.GetItemRequest(&input)

	error = request.Send()
	if error != nil {
		fmt.Println("There was a problem writing to Dynamo DB", error)
		return nil, error
	}

	tokenInfoAttributeValue, ok := response.Item[tokenInfoKey]
	if !ok {
		fmt.Println("There's no token information in the response from DynamoDB")
		return nil, nil // TODO
	}
	fmt.Printf("Got a value back from DynamoDB\n%+v\n", string(tokenInfoAttributeValue.B))

	return nil, nil
}

func newDynamoClient() (*dynamodb.DynamoDB, error) {
	session, error := session.NewSession()
	if error != nil {
		fmt.Printf("Unable to create new AWS Session\n%+v\n", error)
		return nil, error
	}
	return dynamodb.New(session, aws.NewConfig().WithRegion("us-west-2")), nil
}

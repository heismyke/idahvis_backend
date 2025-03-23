package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/heismyke/lambda/types"
)

const TABLE_NAME = "message"

type DynamoDBClient struct{
  databaseStore *dynamodb.DynamoDB
}

func NewDynamoDBClient() *DynamoDBClient{
  dbSession := session.Must(session.NewSession())
  db := dynamodb.New(dbSession) 
  return &DynamoDBClient{
    databaseStore: db,
  }
}


func (u DynamoDBClient) DoesMessageExists(email string) (bool, error) {
  item := &dynamodb.GetItemInput{
    TableName: aws.String(TABLE_NAME),
    Key : map[string]*dynamodb.AttributeValue{
      "email" : {
        S : aws.String(email),
      },
    },
  } 
  result, err := u.databaseStore.GetItem(item) 
  if err != nil{
    return false, err
  }

  if result.Item == nil {
    return false, nil
  }
  
  return true, nil
}


func (u DynamoDBClient) InsertMessage(event types.CreateMessage) error{
  //assemble items
  Item := &dynamodb.PutItemInput{
    TableName : aws.String(TABLE_NAME),
    Item : map[string]*dynamodb.AttributeValue{
      "name" : {
        S : aws.String(event.Name),
      },
      "email" : {
        S : aws.String(event.Email),
      },
      "phone" : {
        S : aws.String(event.Phone),
      },
      "subject" : {
        S : aws.String(event.Subject),
      },
      "message" : {
        S: aws.String(event.Message),
      },
    },
  }

  _, err := u.databaseStore.PutItem(Item)
  if err != nil{
    return fmt.Errorf("error inserting into database")
  }

  return nil 

}



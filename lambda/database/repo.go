package database

import "github.com/heismyke/lambda/types"



type MessageStore interface{
  DoesMessageExists(email string) (bool, error)
  InsertMessage(event types.CreateMessage) error
}

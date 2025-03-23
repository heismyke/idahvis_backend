package app

import (
	"github.com/heismyke/lambda/api"
	"github.com/heismyke/lambda/database"
)

type App struct{
  ApiHandler api.ApiHandler
}

func NewApp() App{
  db := database.NewDynamoDBClient()
  apiHandler := api.NewApiHandler(db) 
  return App{
    ApiHandler: apiHandler,
  }

}





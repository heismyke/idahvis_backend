package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/heismyke/lambda/app"
)

type myEvent struct{
  Name string `json:"name"`
  Email string `json:"email"`
  Phone string `json:"phone"`
  Subject string `json:"subject"`
  Message string `json:"message"`
}
func HandleRequest(event myEvent) (string,error) {
  if len(event.Phone) > 11 || len(event.Phone) < 11  {
    return "", fmt.Errorf("phone must be 11 digit")
  }
  if event.Name == "" || event.Email == "" ||   event.Subject == "" || event.Message == ""{
    return "", fmt.Errorf("fields are required")
  }

  return fmt.Sprintf("successfully called by %s", event.Name), nil
}

func main(){
  myApp := app.NewApp()
  lambda.Start(func(request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse, error){
    // Handle OPTIONS preflight requests
    if request.HTTPMethod == "OPTIONS" {
      return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Headers: map[string]string{
          "Access-Control-Allow-Origin": "https://www.idahvisng.com",
          "Access-Control-Allow-Credentials": "true",
          "Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
          "Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
        },
        Body: "",
      }, nil
    }

    switch request.Path {
    case "/message":
      return myApp.ApiHandler.CreateMessage(request)
    default: 
      return events.APIGatewayProxyResponse{
        Body : "Invalid Request",
        StatusCode : http.StatusBadRequest,
        Headers: map[string]string{
          "Content-Type": "application/json",
          "Access-Control-Allow-Origin": "https://www.idahvisng.com",
          "Access-Control-Allow-Credentials": "true",
          "Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
          "Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
        },
      }, nil
    }
  })
}

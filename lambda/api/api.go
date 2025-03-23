package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/heismyke/lambda/database"
	"github.com/heismyke/lambda/types"
)


type ApiHandler struct{
  dbStore database.MessageStore 
}

func NewApiHandler(dbStore database.MessageStore) ApiHandler{
  return ApiHandler{
    dbStore: dbStore,
  }
}


func (a ApiHandler) CreateMessage(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Define CORS headers
    corsHeaders := map[string]string{
        "Content-Type": "application/json",
        "Access-Control-Allow-Origin": "https://www.idahvisng.com",
        "Access-Control-Allow-Credentials": "true",
        "Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
        "Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
    }

    // Parse request body
    var message types.CreateMessage
    err := json.Unmarshal([]byte(request.Body), &message)
    if err != nil {
        errorResponse, _ := json.Marshal(map[string]string{"error": "invalid request body"})
        return events.APIGatewayProxyResponse{
            Body:       string(errorResponse),
            StatusCode: http.StatusBadRequest,
            Headers:    corsHeaders,
        }, nil
    }
    
    // Your validation logic...
    
    // Store in DynamoDB
    err = a.dbStore.InsertMessage(message)
    if err != nil {
        errorResponse, _ := json.Marshal(map[string]string{"error": "error inserting into database"})
        return events.APIGatewayProxyResponse{
            Body:       string(errorResponse),
            StatusCode: http.StatusInternalServerError,
            Headers:    corsHeaders,
        }, nil
    }
    
    // Send notification email
    
    emailErr := sendContactFormEmail(message)
    response := map[string]any{
    "message":   "Message created successfully",
    "emailSent": emailErr == nil,
  }



    if emailErr != nil {
    log.Printf("Error sending email: %v", emailErr)
    response["emailError"] = emailErr.Error()
    }


    successResponse, _ := json.Marshal(map[string]string{"message": "message created successfully"})
    return events.APIGatewayProxyResponse{
        Body:       string(successResponse),
        StatusCode: http.StatusCreated,
        Headers:    corsHeaders,
    }, nil
}

func sendContactFormEmail(message types.CreateMessage) error {
    // Create a new AWS session
    sess := session.Must(session.NewSession())
    
    // Create an SES client
    sesClient := ses.New(sess)
    
    // Email content
    emailSubject := "New Contact Form Submission from " + message.Name
    htmlBody := fmt.Sprintf(`
        <h1>New Contact Form Submission</h1>
        <p><strong>Name:</strong> %s</p>
        <p><strong>Email:</strong> %s</p>
        <p><strong>Phone:</strong> %s</p>
        <p><strong>Subject:</strong> %s</p>
        <p><strong>Message:</strong> %s</p>
    `, message.Name, message.Email, message.Phone, message.Subject, message.Message)
    
    textBody := fmt.Sprintf(
        "New message from website contact form:\n\n"+
        "Name: %s\n"+
        "Email: %s\n"+
        "Phone: %s\n"+
        "Subject: %s\n"+
        "Message: %s\n",
        message.Name, message.Email, message.Phone, message.Subject, message.Message,
    )
    
    // The recipient is your email address
    toAddress := "mickienorman5@gmail.com" // Replace with YOUR email
    
    // The sender should be a verified email address in SES
    fromAddress := "noreply@idahvisng.com" // Replace with your verified sender
    
    // Build the email
    input := &ses.SendEmailInput{
        Destination: &ses.Destination{
            ToAddresses: []*string{
                aws.String(toAddress),
            },
        },
        Message: &ses.Message{
            Body: &ses.Body{
                Html: &ses.Content{
                    Data: aws.String(htmlBody),
                },
                Text: &ses.Content{
                    Data: aws.String(textBody),
                },
            },
            Subject: &ses.Content{
                Data: aws.String(emailSubject),
            },
        },
        Source: aws.String(fromAddress),
        ReplyToAddresses: []*string{
            aws.String(message.Email), // Set reply-to as the form submitter's email
        },
    }
    
    // Send the email
  _, err := sesClient.SendEmail(input)
    return err
}

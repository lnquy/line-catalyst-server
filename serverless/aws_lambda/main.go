package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Wake(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(os.Getenv("CATALYST_URL"))
	if err != nil {
		log.Printf("Failed to make request to Catalyst: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Catalyst response with not OK status: %s", resp.Status)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusServiceUnavailable}, nil
	}

	log.Println("Catalyst is active")
	return events.APIGatewayProxyResponse{
		Body:       "Catalyst is active",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(Wake)
}

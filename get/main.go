package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sethcenterbar/serverless-movies/data"
)

// Global connection for multiple lambda invocations
var ddb = data.ConnectToDynamoDB()

// Handler runs on each lambda invocation
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]

	movie := data.GetMovieByID(id, ddb)
	fmt.Println(movie)
	return events.APIGatewayProxyResponse{
			Body:       movie.MovieToJSON(),
			StatusCode: 200,
		},
		nil
}

func main() {
	lambda.Start(Handler)
}

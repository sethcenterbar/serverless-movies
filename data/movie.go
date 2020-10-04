package data

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/mitchellh/mapstructure"
	"github.com/oklog/ulid/v2"
)

// ddbMovie wraps the Movie object with DynamoDB specific information such as keys
type ddbMovie struct {
	PK    string
	Movie *Movie
}

// Movie is
type Movie struct {
	Title       string   `json:"title"`
	ReleaseDate string   `json:"releaseDate"`
	Genres      []string `json:"genres"`
	Cast        []string `json:"cast"`
}

// ConnectToDynamoDB creates a ddb connection
func ConnectToDynamoDB() (ddb *dynamodb.DynamoDB) {
	return dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})))
}

// GetULID returns a ulid seeded from time.Now()
func GetULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	u := ulid.MustNew(ulid.Timestamp(t), entropy)
	return u.String()
}

// PutMovie does stuff
func (m *ddbMovie) CreateMovie(ddb *dynamodb.DynamoDB) {
	val, err := dynamodbattribute.MarshalMap(m.Movie)
	if err != nil {
		panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
	}

	_, err = ddb.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(m.PK),
			},
			"details": {
				M: val,
			},
		},
		TableName: aws.String("go-movies"),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to put Record to DynamoDB, %v", err))
	}
}

// DeleteMovie deletes a given movie
func (m *ddbMovie) DeleteMovie(ddb *dynamodb.DynamoDB) {
	_, err := ddb.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(m.PK),
			},
		},
		TableName: aws.String("go-movies"),
	})

	if err != nil {
		fmt.Println(err)
	}
}

// GetMovieByID gets a movie from the database
func GetMovieByID(id string, ddb *dynamodb.DynamoDB) Movie {
	result, err := ddb.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(id),
			},
		},
		TableName: aws.String("go-movies"),
	})
	if err != nil {
		fmt.Println(err)
	}

	movie := &Movie{}
	var rv map[string]interface{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &rv)
	if err != nil {
		fmt.Println(err)
	}
	mapstructure.Decode(rv["details"], movie)

	return *movie
}

// MovieToJSON Takes a movie and returns it as JSON
func (m *Movie) MovieToJSON() string {
	bytes, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	return string(bytes)
}

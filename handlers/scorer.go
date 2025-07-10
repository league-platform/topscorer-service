package handlers

import (
    "fmt"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/gofiber/fiber/v2"
)

type Scorer struct {
    Player   string    `json:"player"`
    Team     string    `json:"team"`
    Goals    int       `json:"goals"`
    Updated  time.Time `json:"updated"`
}

func CreateScorer(c *fiber.Ctx) error {
    var scorer Scorer
    if err := c.BodyParser(&scorer); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }
    scorer.Updated = time.Now()

    sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
    svc := dynamodb.New(sess)

    av, _ := dynamodbattribute.MarshalMap(scorer)
    input := &dynamodb.PutItemInput{
        TableName: aws.String("TopScorers"),
        Item:      av,
    }

    if _, err := svc.PutItem(input); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save scorer"})
    }

    fmt.Printf("EVENT: scorer.updated -> %+v\n", scorer)
    return c.Status(fiber.StatusCreated).JSON(scorer)
}

func GetScorers(c *fiber.Ctx) error {
    sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
    svc := dynamodb.New(sess)

    input := &dynamodb.ScanInput{
        TableName: aws.String("TopScorers"),
    }

    result, err := svc.Scan(input)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve scorers"})
    }

    var scorers []Scorer
    _ = dynamodbattribute.UnmarshalListOfMaps(result.Items, &scorers)
    return c.JSON(scorers)
}

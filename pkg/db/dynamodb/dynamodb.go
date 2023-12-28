package dynamodb

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bazuker/backend-bootstrap/pkg/db"
)

// DB represents a configurable DynamoDB instance.
// Note: it is important that fields in the table match the JSON representation of the fields from the schema.
type DB struct {
	cfg Config
	db  *dynamodb.DynamoDB
}

type Config struct {
	AWSSession     *session.Session
	UsersTableName string
}

func New(cfg Config) *DB {
	if cfg.AWSSession == nil {
		cfg.AWSSession = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}
	return &DB{
		cfg: cfg,
		db:  dynamodb.New(cfg.AWSSession),
	}
}

func (d DB) CreateUser(user *db.User) error {
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(d.cfg.UsersTableName),
	}

	_, err = d.db.PutItem(input)
	return err
}

func (d DB) UpdateUser(user *db.User) error {
	//TODO implement me
	panic("implement me")
}

func (d DB) UpdateUserPhotoURL(userID, photoURL string) error {
	_, err := d.db.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				S: aws.String(photoURL),
			},
		},
		TableName: aws.String(d.cfg.UsersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(userID),
			},
		},
		UpdateExpression: aws.String("set photoURL = :r"),
		ReturnValues:     aws.String("NONE"),
	})
	return err
}

func (d DB) GetUserByID(id string) (db.User, error) {
	if id == "" {
		return db.User{}, errors.New("missing ID")
	}

	result, err := d.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.cfg.UsersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return db.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	if result == nil || len(result.Item) == 0 {
		return db.User{}, db.ErrNotFound
	}

	var u db.User
	err = dynamodbattribute.UnmarshalMap(result.Item, &u)
	if err != nil {
		return db.User{}, fmt.Errorf("failed to unmarshal user: %w", err)
	}
	return u, nil
}

func (d DB) GetUserByEmail(email string) (db.User, error) {
	if email == "" {
		return db.User{}, errors.New("missing email")
	}

	result, err := d.db.Query(&dynamodb.QueryInput{
		TableName: aws.String(d.cfg.UsersTableName),
		IndexName: aws.String("email-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"email": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(email),
					},
				},
			},
		},
	})
	if err != nil {
		return db.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}
	if result == nil || len(result.Items) == 0 || len(result.Items[0]) == 0 {
		return db.User{}, db.ErrNotFound
	}

	var u db.User
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &u)
	if err != nil {
		return db.User{}, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return u, nil
}

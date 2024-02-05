package usecase

import (
	"context"
	"database-remote-commander/internal/model"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os/exec"
	"strings"
)

type RemoteQueryUseCase struct {
	DB     *mongo.Database
	Log    *logrus.Logger
	Config *viper.Viper
}

func NewRemoteQueryUseCase(db *mongo.Database, log *logrus.Logger, config *viper.Viper) *RemoteQueryUseCase {
	return &RemoteQueryUseCase{
		DB:     db,
		Log:    log,
		Config: config,
	}
}

func (c *RemoteQueryUseCase) ExecQuery(ctx context.Context, request *model.QueryRequest) (result string, error error) {
	if len(request.Query) == 0 {
		return "", errors.New("query is empty")
	}

	querySplit := strings.Split(request.Query, ".")
	if querySplit[0] != "db" {
		return "", errors.New("first query is not db")
	}

	collection := querySplit[1]

	err := c.isCollectionExist(collection)
	if err != nil {
		return "", errors.New(fmt.Sprintf("%v", err))
	}

	result, err = c.execMongoShell(request.Query)
	if err != nil {
		return "", errors.New(fmt.Sprintf("%v", err))
	}

	return result, nil
}

func (c *RemoteQueryUseCase) isCollectionExist(collection string) error {
	filter := bson.D{{"name", collection}}
	collectionCursor, err := c.DB.ListCollections(context.TODO(), filter)

	if err != nil {
		return errors.New(fmt.Sprintf("%v", err))
	}

	isCollectionExist := false
	for collectionCursor.Next(context.Background()) {
		var result bson.M
		err := collectionCursor.Decode(&result)
		if err != nil {
			return errors.New(fmt.Sprintf("%v", err))
		}

		if result["name"] == collection {
			isCollectionExist = true
			break
		}
	}

	if !isCollectionExist {
		return errors.New(fmt.Sprintf("Collection %s is not exist", collection))
	}

	return nil
}

func (c *RemoteQueryUseCase) execMongoShell(query string) (result string, error error) {
	username := c.Config.GetString("database.username")
	password := c.Config.GetString("database.password")
	host := c.Config.GetString("database.host")
	port := c.Config.GetInt("database.port")
	database := c.Config.GetString("database.name")

	cmd := exec.Command("mongosh",
		"--quiet",
		fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", username, password, host, port, database),
		"--authenticationDatabase", "admin",
		"--eval", fmt.Sprintf("'%s'", query),
	)

	execResult, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return string(execResult), nil
}

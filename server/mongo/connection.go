package mongo

import (
	"context"
	"fmt"
	"github.com/hellphone/gomud/domain/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: store credentials as env vars
func ConnectToDB(cfg *models.Config) (context.Context, *mongo.Client, error) {
	uri := fmt.Sprintf(
		"mongodb+srv://%s:%s@cluster0.41k9y.mongodb.net/myFirstDatabase?retryWrites=true&w=majority",
		cfg.Database.Username,
		cfg.Database.Password,
	)
	clientOptions := options.Client().ApplyURI(uri)
	// TODO: figure out how to deal with the context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, err
	}

	return ctx, client, nil
}

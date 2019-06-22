package mongoclient

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/remussino/reddit-consumer-go/data-retriever/httpclient"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveSubmissionsToMongo(submissions *[]httpclient.Submission, wg *sync.WaitGroup) {

	if len(*submissions) == 0 {
		log.Println("Error. No submissions to save")
		return
	}

	log.Printf("Saving %d submissions to Mongo\n", len(*submissions))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println(err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	submissionCollection := client.Database("reddit").Collection("submissions")

	documents := *submissions
	var toInsert []interface{}
	for _, d := range documents {
		toInsert = append(toInsert, d)
	}

	submissionCollection.InsertMany(ctx, toInsert)
	log.Println("Saved to Mongo")

	wg.Done()
}

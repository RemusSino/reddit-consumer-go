package mongoclient

import (
	"log"
	"sync"
	"time"

	"github.com/remussino/reddit-consumer-go/data-retriever/httpclient"
)

func SaveSubmissionsToMongo(submissions *[]httpclient.Submission, wg *sync.WaitGroup) {
	// var s []httpclient.Submission
	log.Printf("Saving %d submissions to Mongo\n", len(*submissions))
	time.Sleep(time.Millisecond * 2000)
	//todo
	wg.Done()
}

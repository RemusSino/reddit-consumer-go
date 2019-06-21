package main

import (
	"log"
	"sync"

	"github.com/remussino/reddit-consumer-go/data-retriever/httpclient"
	"github.com/remussino/reddit-consumer-go/data-retriever/mongoclient"
)

var wg sync.WaitGroup

func retrieveAndSaveSubmission(subreddit string) {
	log.Println("Retrieving submissions for subreddit", subreddit)
	submissions := httpclient.RetrieveSubmissions(subreddit, "")
	log.Printf("Found %d submissions for subreddit %s\n", len(submissions), subreddit)
	wg.Add(1)
	go mongoclient.SaveSubmissionsToMongo(&submissions, &wg)
	wg.Done()
}

func main() {
	subreddits := []string{"Python", "Java"}
	for _, v := range subreddits {
		wg.Add(1)
		go retrieveAndSaveSubmission(v)
	}
	wg.Wait()
}

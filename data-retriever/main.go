package main

import (
	"log"
	"sync"
	"time"

	"github.com/remussino/reddit-consumer-go/data-retriever/httpclient"
	"github.com/remussino/reddit-consumer-go/data-retriever/mongoclient"
)

var wg sync.WaitGroup
var lastSubmissionIdMap map[string]string = make(map[string]string)

func retrieveAndSaveSubmission(subreddit string) {
	log.Println("Retrieving submissions for subreddit", subreddit)
	log.Println("Last submission id for subreddit", subreddit, "is", lastSubmissionIdMap[subreddit])

	submissions := httpclient.RetrieveSubmissions(subreddit, lastSubmissionIdMap[subreddit])

	log.Printf("Found %d submissions for subreddit %s\n", len(submissions), subreddit)

	if len(submissions) > 0 {
		lastSubmissionIdMap[subreddit] = "t3_" + submissions[0].RedditId
		wg.Add(1)
		go mongoclient.SaveSubmissionsToMongo(&submissions, &wg)
	}

	wg.Done()
}

func main() {
	subreddits := []string{"Python", "Java"}

	for {
		log.Println("Start retrieving data at ", time.Now())
		for _, v := range subreddits {
			wg.Add(1)
			go retrieveAndSaveSubmission(v)
		}
		log.Println("Waiting for WG")
		wg.Wait()
		log.Println("Sleeping 10 sec before next http call")
		time.Sleep(10 * time.Second)
	}
}

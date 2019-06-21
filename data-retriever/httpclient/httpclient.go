package httpclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Submission struct {
	Subreddit string
	RedditId  string
	Title     string
	Timestamp float64
}

func RetrieveSubmissions(subreddit string, lastSubmissionFullname string) []Submission {
	urlTemplate := "https://www.reddit.com/r/%s/new/.json?limit=100&before=%s"
	url := fmt.Sprintf(urlTemplate, subreddit, lastSubmissionFullname)
	log.Println("Preparing to retrieve submissions from url", url)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-agent", `Chrome"`)

	resp, err := client.Do(req)

	log.Println("Http response status code", resp.Status)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	return unmarshalJSON(responseData)
}

func unmarshalJSON(jsonData []byte) []Submission {
	submissionSlice := make([]Submission, 0, 100)

	var v interface{}
	json.Unmarshal(jsonData, &v)
	data := v.(map[string]interface{})
	dataChildren := data["data"].(map[string]interface{})
	submissions := dataChildren["children"].([]interface{})

	for _, s := range submissions {
		child := s.(map[string]interface{})
		submission := child["data"].(map[string]interface{})
		submissionSlice = append(submissionSlice, extractSubmission(&submission))
	}
	return submissionSlice
}

func extractSubmission(s *map[string]interface{}) Submission {
	submissionMap := *s

	subreddit := submissionMap["subreddit"].(string)
	redditId := submissionMap["id"].(string)
	title := submissionMap["title"].(string)
	timestamp := submissionMap["created_utc"].(float64)

	submission := Submission{
		Subreddit: subreddit,
		RedditId:  redditId,
		Title:     title,
		Timestamp: timestamp}

	return submission
}

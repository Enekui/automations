package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var topic = "arn:aws:sns:eu-west-2:075107581003:tower-monitoring"

func main() {
	lambda.Start(Handler)
}

func Handler() {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	var wg sync.WaitGroup
	links := os.Args[1:]

	wg.Add(len(links))

	for _, l := range links {
		go func(link string, wg *sync.WaitGroup) {
			_, err := GetLink(link)
			if err != nil {
				PublishSNSNotification(topic, fmt.Sprint(err), sess)
			}
			defer wg.Done()
		}(l, &wg)
	}

	wg.Wait()
}

func GetLink(l string) (*http.Response, error) {

	resp, err := http.Get(l)
	if err != nil {
		log.Println("DOWN:", l, "ERROR:", err)
		return resp, err
	} else if resp.StatusCode != 200 {
		log.Println("DOWN:", l, "ERROR:", resp.StatusCode)
		return resp, errors.New("Link DOWN: " + l + " Status code " + resp.Status)
	}

	log.Println("UP:", l, resp.Status)
	return resp, nil
}

func PublishSNSNotification(t, m string, sess *session.Session) (*sns.PublishOutput, error) {
	var result *sns.PublishOutput

	svc := sns.New(sess)
	input := &sns.PublishInput{
		Message:  aws.String(m),
		TopicArn: aws.String(t),
	}

	result, err := svc.Publish(input)
	if err != nil {
		return result, err
	}

	return result, err
}

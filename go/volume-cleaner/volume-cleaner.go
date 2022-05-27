package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// function to retrive the volumes
// using the "available" "status" as filter
func GetVolumes(sess *session.Session) (*ec2.DescribeVolumesOutput, error) {

	var (
		result *ec2.DescribeVolumesOutput
		err    error
	)

	svc := ec2.New(sess)
	input := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("status"),
				Values: []*string{
					aws.String("available"),
				},
			},
		},
	}

	result, err = svc.DescribeVolumes(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return result, errors.New(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return result, errors.New(err.Error())
		}
	}

	return result, nil
}

// function to delete the volume, this function takes a string as input
// the string will be the volume ID
func DeleteVolume(v string, sess *session.Session) (*ec2.DeleteVolumeOutput, error) {

	var (
		result *ec2.DeleteVolumeOutput
		err    error
	)

	svc := ec2.New(sess)
	input := &ec2.DeleteVolumeInput{
		VolumeId: aws.String(v),
	}

	result, err = svc.DeleteVolume(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return result, errors.New(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return result, errors.New(err.Error())
		}
	}

	return result, nil
}

// The handler function is the funtion used by lambda to call the main code
// All the code logic is inside of this function
func Handler() {

	var (
		sess            *session.Session
		err             error
		count           int
		volumes         *ec2.DescribeVolumesOutput
		region, profile string
	)

	// Craeting a session using the Must function
	// this ensure the session will be working before creating it
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
		Config: aws.Config{
			Region: aws.String(region),
		},
	}))

	// Get the volumes
	volumes, err = GetVolumes(sess)
	switch {
	case err != nil:
		fmt.Fprintln(os.Stderr, err)
		return
	case len(volumes.Volumes) == 0:
		fmt.Println("0 orphan volumes found!")
		return
	}

	// loop through the volumes to delete them
	for _, v := range volumes.Volumes {
		fmt.Println("Deleting orphan volume: " + *v.VolumeId)
		_, err := DeleteVolume(*v.VolumeId, sess)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting orphan volume: %v\n%v ", *v.VolumeId, err)
			return
		}
		count++
	}

	fmt.Printf("%v volemes has been deleted!\n", count)
}

func main() {
	// invoking habdle function from lambda
	lambda.Start(Handler)
}

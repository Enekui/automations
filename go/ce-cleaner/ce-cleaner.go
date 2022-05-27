package main

import (
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/aws/aws-sdk-go/service/iam"
)

func main() {

	var (
		roles []string
		// launchTemplates []string
	)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "development",
	}))

	jq, err := GetJobQueue("", sess)
	if err != nil {
		log.Println(err)
		return
	}

	for _, i := range jq.JobQueues {
		result, err := GetJobQueue(*i.JobQueueName, sess)
		if err != nil {
			log.Println(err)
			continue
		}
		if *result.JobQueues[0].JobQueueName == "INVALID" {
			log.Println("Job Queue", *i.JobQueueName, "is in invalid state.")
			continue
		}
		log.Println("Disabling JobQueue:", *i.JobQueueName)
		_, err = DisableJobQueue(*i.JobQueueName, sess)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Waiting for Job Queue", *i.JobQueueName, "to be disabled")
		for status := false; !status; {
			result, err = GetJobQueue(*i.JobQueueName, sess)
			if *result.JobQueues[0].Status == "VALID" {
				status = true
			}
		}
	}

	for _, i := range jq.JobQueues {
		log.Println("Deleting JobQueue:", *i.JobQueueName)
		_, err := DeleteJobQueue(*i.JobQueueName, sess)
		if err != nil {
			log.Println(err)
			return
		}

	}

	ce, err := GetComputeEnvironment("", sess)
	if err != nil {
		log.Println(err)
	}

	for _, i := range ce.ComputeEnvironments {
		result, err := GetComputeEnvironment(*i.ComputeEnvironmentName, sess)
		if err != nil {
			log.Println(err)
			continue
		}
		if *result.ComputeEnvironments[0].Status == "INVALID" {
			log.Println("ComputeEnvironment:", *i.ComputeEnvironmentName, "is in invalid state")
			continue
		}
		log.Println("Disabling ComputeEnvironment:", *i.ComputeEnvironmentName)
		_, err = DisableComputeEnvironment(*i.ComputeEnvironmentName, sess)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Waiting for ComputeEnvironment:", *i.ComputeEnvironmentName, "to be disabled...")
		for status := false; !status; {
			result, err = GetComputeEnvironment(*i.ComputeEnvironmentName, sess)
			if err != nil {
				log.Println(err)
				continue
			}
			if *result.ComputeEnvironments[0].Status == "VALID" {
				status = true
			}
		}
	}

	for _, i := range ce.ComputeEnvironments {
		result, err := GetComputeEnvironment(*i.ComputeEnvironmentName, sess)
		if *result.ComputeEnvironments[0].Status == "INVALID" {
			continue
		}
		log.Println("Deleting ComputeEnvironment:", *i.ComputeEnvironmentName)
		_, err = DeleteComputeEnvironment(*i.ComputeEnvironmentName, sess)
		if err != nil {
			log.Println(err)
			continue
		}

		roles = append(roles, strings.Split(*i.ServiceRole, "/")[1])
	}

	for _, r := range roles {
		log.Println("Deleting service role:", r)
		_, err := DeleteRole(r, sess)
		if err != nil {
			log.Println(err)
			continue
		}
	}

}

func GetComputeEnvironment(ce string, sess *session.Session) (*batch.DescribeComputeEnvironmentsOutput, error) {

	var result *batch.DescribeComputeEnvironmentsOutput
	var input *batch.DescribeComputeEnvironmentsInput

	if ce != "" {
		input = &batch.DescribeComputeEnvironmentsInput{
			ComputeEnvironments: []*string{
				aws.String(ce),
			},
		}
	}

	svc := batch.New(sess)

	result, err := svc.DescribeComputeEnvironments(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case batch.ErrCodeClientException:

				return result, errors.New(batch.ErrCodeClientException + aerr.Error())
			case batch.ErrCodeServerException:
				return result, errors.New(batch.ErrCodeServerException + aerr.Error())
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

func GetJobQueue(jq string, sess *session.Session) (*batch.DescribeJobQueuesOutput, error) {

	var result *batch.DescribeJobQueuesOutput
	var input *batch.DescribeJobQueuesInput

	svc := batch.New(sess)
	if jq != "" {
		input = &batch.DescribeJobQueuesInput{
			JobQueues: []*string{
				aws.String(jq),
			},
		}
	}
	result, err := svc.DescribeJobQueues(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case batch.ErrCodeClientException:

				return result, errors.New(batch.ErrCodeClientException + aerr.Error())
			case batch.ErrCodeServerException:

				return result, errors.New(batch.ErrCodeServerException + aerr.Error())
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

func DisableJobQueue(jq string, sess *session.Session) (*batch.UpdateJobQueueOutput, error) {

	var result *batch.UpdateJobQueueOutput

	svc := batch.New(sess)
	input := &batch.UpdateJobQueueInput{
		JobQueue: aws.String(jq),
		State:    aws.String("DISABLED"),
	}

	result, err := svc.UpdateJobQueue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case batch.ErrCodeClientException:
				return result, errors.New(batch.ErrCodeClientException + aerr.Error())
			case batch.ErrCodeServerException:
				return result, errors.New(batch.ErrCodeServerException + aerr.Error())
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

func DeleteJobQueue(jq string, sess *session.Session) (*batch.DeleteJobQueueOutput, error) {

	var result *batch.DeleteJobQueueOutput

	svc := batch.New(sess)
	input := &batch.DeleteJobQueueInput{
		JobQueue: aws.String(jq),
	}

	result, err := svc.DeleteJobQueue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case batch.ErrCodeClientException:
				return result, errors.New(batch.ErrCodeClientException + aerr.Error())
			case batch.ErrCodeServerException:
				return result, errors.New(batch.ErrCodeServerException + aerr.Error())
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

func DisableComputeEnvironment(ce string, sess *session.Session) (*batch.UpdateComputeEnvironmentOutput, error) {

	var result *batch.UpdateComputeEnvironmentOutput

	svc := batch.New(sess)
	input := &batch.UpdateComputeEnvironmentInput{
		ComputeEnvironment: aws.String(ce),
		State:              aws.String("DISABLED"),
	}

	result, err := svc.UpdateComputeEnvironment(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case batch.ErrCodeClientException:
				return result, errors.New(batch.ErrCodeClientException + aerr.Error())
			case batch.ErrCodeServerException:
				return result, errors.New(batch.ErrCodeServerException + aerr.Error())
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

func DeleteComputeEnvironment(ce string, sess *session.Session) (*batch.DeleteComputeEnvironmentOutput, error) {

	var result *batch.DeleteComputeEnvironmentOutput

	svc := batch.New(sess)
	input := &batch.DeleteComputeEnvironmentInput{
		ComputeEnvironment: aws.String(ce),
	}

	result, err := svc.DeleteComputeEnvironment(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case batch.ErrCodeClientException:
				return result, errors.New(batch.ErrCodeClientException + aerr.Error())
			case batch.ErrCodeServerException:
				return result, errors.New(batch.ErrCodeServerException + aerr.Error())
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

func GetJobDefinition(sess *session.Session) (*batch.DescribeJobDefinitionsOutput, error) {

	var result *batch.DescribeJobDefinitionsOutput

	svc := batch.New(sess)
	// input := &batch.DescribeJobDefinitionsInput{
	// 	Status: aws.String("ACTIVE"),
	// }

	result, err := svc.DescribeJobDefinitions(nil)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case batch.ErrCodeClientException:
				return result, errors.New(batch.ErrCodeClientException + aerr.Error())
			case batch.ErrCodeServerException:
				return result, errors.New(batch.ErrCodeServerException + aerr.Error())
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

func GetRole(r string, sess *session.Session) (*iam.GetRoleOutput, error) {

	var result *iam.GetRoleOutput

	svc := iam.New(sess)

	input := &iam.GetRoleInput{
		RoleName: aws.String(r),
	}

	result, err := svc.GetRole(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return result, errors.New(iam.ErrCodeNoSuchEntityException + aerr.Error())
			case iam.ErrCodeServiceFailureException:
				return result, errors.New(iam.ErrCodeServiceFailureException + aerr.Error())
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

func DeleteRole(r string, sess *session.Session) (*iam.DeleteRoleOutput, error) {

	var result *iam.DeleteRoleOutput

	svc := iam.New(sess)
	input := &iam.DeleteRoleInput{
		RoleName: aws.String(r),
	}

	result, err := svc.DeleteRole(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return result, errors.New(iam.ErrCodeNoSuchEntityException + aerr.Error())
			case iam.ErrCodeDeleteConflictException:
				return result, errors.New(iam.ErrCodeDeleteConflictException + aerr.Error())
			case iam.ErrCodeLimitExceededException:
				return result, errors.New(iam.ErrCodeLimitExceededException + aerr.Error())
			case iam.ErrCodeUnmodifiableEntityException:
				return result, errors.New(iam.ErrCodeUnmodifiableEntityException + aerr.Error())
			case iam.ErrCodeConcurrentModificationException:
				return result, errors.New(iam.ErrCodeConcurrentModificationException + aerr.Error())
			case iam.ErrCodeServiceFailureException:
				return result, errors.New(iam.ErrCodeServiceFailureException + aerr.Error())
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

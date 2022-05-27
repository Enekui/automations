package main

import (
	"errors"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/rds"
)

func GetCluster(c string, sess *session.Session) (*rds.DescribeDBClustersOutput, error) {

	var result *rds.DescribeDBClustersOutput

	svc := rds.New(sess)
	input := &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(c),
	}

	result, err := svc.DescribeDBClusters(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBClusterNotFoundFault:
				return result, errors.New(rds.ErrCodeDBClusterNotFoundFault + aerr.Error())
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

func GetKMSKeyAlias(sess *session.Session) (*kms.ListAliasesOutput, error) {
	var result *kms.ListAliasesOutput

	svc := kms.New(sess)
	input := &kms.ListAliasesInput{}

	result, err := svc.ListAliases(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case kms.ErrCodeDependencyTimeoutException:
				return result, errors.New(kms.ErrCodeDependencyTimeoutException + aerr.Error())
			case kms.ErrCodeInvalidMarkerException:
				return result, errors.New(kms.ErrCodeInvalidMarkerException + aerr.Error())
			case kms.ErrCodeInternalException:
				return result, errors.New(kms.ErrCodeInternalException + aerr.Error())
			case kms.ErrCodeInvalidArnException:
				return result, errors.New(kms.ErrCodeInvalidArnException + aerr.Error())
			case kms.ErrCodeNotFoundException:
				return result, errors.New(kms.ErrCodeNotFoundException + aerr.Error())
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

func ListKMSKeys(k string, sess *session.Session) (*kms.ListKeysOutput, error) {
	var result *kms.ListKeysOutput

	svc := kms.New(sess)
	input := &kms.ListKeysInput{}

	result, err := svc.ListKeys(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case kms.ErrCodeDependencyTimeoutException:
				return result, errors.New(kms.ErrCodeDependencyTimeoutException + aerr.Error())
			case kms.ErrCodeInternalException:
				return result, errors.New(kms.ErrCodeInternalException + aerr.Error())
			case kms.ErrCodeInvalidMarkerException:
				return result, errors.New(kms.ErrCodeInvalidMarkerException + aerr.Error())
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

func GetKMSKey(k string, sess *session.Session) (*kms.DescribeKeyOutput, error) {

	var result *kms.DescribeKeyOutput

	svc := kms.New(sess)
	input := &kms.DescribeKeyInput{
		KeyId: aws.String("alias/" + k),
	}

	result, err := svc.DescribeKey(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case kms.ErrCodeNotFoundException:
				return result, errors.New(kms.ErrCodeNotFoundException + aerr.Error())
			case kms.ErrCodeInvalidArnException:
				return result, errors.New(kms.ErrCodeInvalidArnException + aerr.Error())
			case kms.ErrCodeDependencyTimeoutException:
				return result, errors.New(kms.ErrCodeDependencyTimeoutException + aerr.Error())
			case kms.ErrCodeInternalException:
				return result, errors.New(kms.ErrCodeInternalException + aerr.Error())
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

func GetClusterSnapshot(s, t string, sess *session.Session) (*rds.DescribeDBClusterSnapshotsOutput, error) {
	var (
		result *rds.DescribeDBClusterSnapshotsOutput
		input  *rds.DescribeDBClusterSnapshotsInput
	)

	svc := rds.New(sess)
	if t == "" {
		input = &rds.DescribeDBClusterSnapshotsInput{
			DBClusterSnapshotIdentifier: aws.String(s),
			SnapshotType:                aws.String("manual"),
			DBClusterIdentifier:         aws.String(SourceClusterName),
		}
	} else if t == "shared" {
		input = &rds.DescribeDBClusterSnapshotsInput{
			DBClusterSnapshotIdentifier: aws.String(s),
			SnapshotType:                aws.String(t), // Use "shared" instead
			DBClusterIdentifier:         aws.String(SourceClusterName),
		}
	}

	result, err := svc.DescribeDBClusterSnapshots(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBClusterSnapshotNotFoundFault:
				return result, errors.New(rds.ErrCodeDBClusterSnapshotNotFoundFault + aerr.Error())
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

func CreateClusterSnapshot(c, s string, sess *session.Session) (*rds.CreateDBClusterSnapshotOutput, error) {
	var result *rds.CreateDBClusterSnapshotOutput

	svc := rds.New(sess)
	input := &rds.CreateDBClusterSnapshotInput{
		DBClusterIdentifier:         aws.String(c),
		DBClusterSnapshotIdentifier: aws.String(s),
	}

	result, err := svc.CreateDBClusterSnapshot(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBClusterSnapshotAlreadyExistsFault:
				return result, errors.New(rds.ErrCodeDBClusterSnapshotAlreadyExistsFault + aerr.Error())
			case rds.ErrCodeInvalidDBClusterStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBClusterStateFault + aerr.Error())
			case rds.ErrCodeDBClusterNotFoundFault:
				return result, errors.New(rds.ErrCodeDBClusterNotFoundFault + aerr.Error())
			case rds.ErrCodeSnapshotQuotaExceededFault:
				return result, errors.New(rds.ErrCodeSnapshotQuotaExceededFault + aerr.Error())
			case rds.ErrCodeInvalidDBClusterSnapshotStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBClusterSnapshotStateFault + aerr.Error())
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

func CopyClusterSnapshot(s, t, k string, sess *session.Session) (*rds.CopyDBClusterSnapshotOutput, error) {
	var result *rds.CopyDBClusterSnapshotOutput

	svc := rds.New(sess)
	input := &rds.CopyDBClusterSnapshotInput{
		SourceDBClusterSnapshotIdentifier: aws.String(s),
		TargetDBClusterSnapshotIdentifier: aws.String(t),
		KmsKeyId:                          aws.String("alias/" + k),
		SourceRegion:                      aws.String(SourceProfileRegion),
		DestinationRegion:                 aws.String(DestinationProfileRegion),
	}

	input.SetDestinationRegion(DestinationProfileRegion)

	result, err := svc.CopyDBClusterSnapshot(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBClusterSnapshotAlreadyExistsFault:
				return result, errors.New(rds.ErrCodeDBClusterSnapshotAlreadyExistsFault + aerr.Error())
			case rds.ErrCodeDBClusterSnapshotNotFoundFault:
				return result, errors.New(rds.ErrCodeDBClusterSnapshotNotFoundFault + aerr.Error())
			case rds.ErrCodeInvalidDBClusterStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBClusterStateFault + aerr.Error())
			case rds.ErrCodeInvalidDBClusterSnapshotStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBClusterSnapshotStateFault + aerr.Error())
			case rds.ErrCodeSnapshotQuotaExceededFault:
				return result, errors.New(rds.ErrCodeSnapshotQuotaExceededFault + aerr.Error())
			case rds.ErrCodeKMSKeyNotAccessibleFault:
				return result, errors.New(rds.ErrCodeKMSKeyNotAccessibleFault + aerr.Error())
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

func ShareClusterSnapshot(s, id string, sess *session.Session) (*rds.ModifyDBClusterSnapshotAttributeOutput, error) {
	svc := rds.New(sess)

	input := &rds.ModifyDBClusterSnapshotAttributeInput{
		AttributeName:               aws.String("restore"),
		DBClusterSnapshotIdentifier: aws.String(s),
		ValuesToAdd: []*string{
			aws.String(id),
		},
	}

	result, err := svc.ModifyDBClusterSnapshotAttribute(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBClusterSnapshotNotFoundFault:
				return result, errors.New(rds.ErrCodeDBClusterSnapshotNotFoundFault + aerr.Error())
			case rds.ErrCodeInvalidDBClusterSnapshotStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBClusterSnapshotStateFault + aerr.Error())
			case rds.ErrCodeSharedSnapshotQuotaExceededFault:
				return result, errors.New(rds.ErrCodeSharedSnapshotQuotaExceededFault + aerr.Error())
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

func CreateClusterFromSnapshot(sess *session.Session) (*rds.RestoreDBClusterFromSnapshotOutput, error) {
	var result *rds.RestoreDBClusterFromSnapshotOutput

	svc := rds.New(sess)
	input := &rds.RestoreDBClusterFromSnapshotInput{
		DBClusterIdentifier: aws.String(DestinationClusterName),
		Engine:              aws.String(DestinationClusterEngine),
		EngineVersion:       aws.String(DestinationClusterEngineVersion),
		EngineMode:          aws.String(DestinationClusterEngineMode),
		DBSubnetGroupName:   aws.String(DestinationClusterSubnetGroup),
		DeletionProtection:  aws.Bool(true),
		KmsKeyId:            aws.String("alias/" + DestinationKMSKeyAlias),
		VpcSecurityGroupIds: []*string{
			aws.String(DestinationClusterSecurityGroup),
		},
		SnapshotIdentifier: aws.String(MigrationSnapshotARN),
	}

	result, err := svc.RestoreDBClusterFromSnapshot(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBClusterAlreadyExistsFault:
				return result, errors.New(rds.ErrCodeDBClusterAlreadyExistsFault + aerr.Error())
			case rds.ErrCodeDBClusterQuotaExceededFault:
				return result, errors.New(rds.ErrCodeDBClusterQuotaExceededFault + aerr.Error())
			case rds.ErrCodeStorageQuotaExceededFault:
				return result, errors.New(rds.ErrCodeStorageQuotaExceededFault + aerr.Error())
			case rds.ErrCodeDBSubnetGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeDBSubnetGroupNotFoundFault + aerr.Error())
			case rds.ErrCodeDBSnapshotNotFoundFault:
				return result, errors.New(rds.ErrCodeDBSnapshotNotFoundFault + aerr.Error())
			case rds.ErrCodeDBClusterSnapshotNotFoundFault:
				return result, errors.New(rds.ErrCodeDBClusterSnapshotNotFoundFault + aerr.Error())
			case rds.ErrCodeInsufficientDBClusterCapacityFault:
				return result, errors.New(rds.ErrCodeInsufficientDBClusterCapacityFault + aerr.Error())
			case rds.ErrCodeInsufficientStorageClusterCapacityFault:
				return result, errors.New(rds.ErrCodeInsufficientStorageClusterCapacityFault + aerr.Error())
			case rds.ErrCodeInvalidDBSnapshotStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBSnapshotStateFault + aerr.Error())
			case rds.ErrCodeInvalidDBClusterSnapshotStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBClusterSnapshotStateFault + aerr.Error())
			case rds.ErrCodeInvalidVPCNetworkStateFault:
				return result, errors.New(rds.ErrCodeInvalidVPCNetworkStateFault + aerr.Error())
			case rds.ErrCodeInvalidRestoreFault:
				return result, errors.New(rds.ErrCodeInvalidRestoreFault + aerr.Error())
			case rds.ErrCodeInvalidSubnet:
				return result, errors.New(rds.ErrCodeInvalidSubnet + aerr.Error())
			case rds.ErrCodeOptionGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeOptionGroupNotFoundFault + aerr.Error())
			case rds.ErrCodeKMSKeyNotAccessibleFault:
				return result, errors.New(rds.ErrCodeKMSKeyNotAccessibleFault + aerr.Error())
			case rds.ErrCodeDomainNotFoundFault:
				return result, errors.New(rds.ErrCodeDomainNotFoundFault + aerr.Error())
			case rds.ErrCodeDBClusterParameterGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeDBClusterParameterGroupNotFoundFault + aerr.Error())
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

func SetCluster(c string, sess *session.Session) (*rds.ModifyDBClusterOutput, error) {
	var result *rds.ModifyDBClusterOutput

	svc := rds.New(sess)
	input := &rds.ModifyDBClusterInput{
		ApplyImmediately:           aws.Bool(true),
		DBClusterIdentifier:        aws.String(c),
		PreferredBackupWindow:      aws.String("22:00-06:00"),
		PreferredMaintenanceWindow: aws.String("Tue:05:00-Tue:05:30"),
		DeletionProtection:         aws.Bool(true),
		BackupRetentionPeriod:      aws.Int64(1),
	}

	result, err := svc.ModifyDBCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBClusterNotFoundFault:
				return result, errors.New(rds.ErrCodeDBClusterNotFoundFault + aerr.Error())
			case rds.ErrCodeInvalidDBClusterStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBClusterStateFault + aerr.Error())
			case rds.ErrCodeStorageQuotaExceededFault:
				return result, errors.New(rds.ErrCodeStorageQuotaExceededFault + aerr.Error())
			case rds.ErrCodeDBSubnetGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeDBSubnetGroupNotFoundFault + aerr.Error())
			case rds.ErrCodeInvalidVPCNetworkStateFault:
				return result, errors.New(rds.ErrCodeInvalidVPCNetworkStateFault + aerr.Error())
			case rds.ErrCodeInvalidDBSubnetGroupStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBSubnetGroupStateFault + aerr.Error())
			case rds.ErrCodeInvalidSubnet:
				return result, errors.New(rds.ErrCodeInvalidSubnet + aerr.Error())
			case rds.ErrCodeDBClusterParameterGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeDBClusterParameterGroupNotFoundFault + aerr.Error())
			case rds.ErrCodeInvalidDBSecurityGroupStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBSecurityGroupStateFault + aerr.Error())
			case rds.ErrCodeInvalidDBInstanceStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBInstanceStateFault + aerr.Error())
			case rds.ErrCodeDBClusterAlreadyExistsFault:
				return result, errors.New(rds.ErrCodeDBClusterAlreadyExistsFault + aerr.Error())
			case rds.ErrCodeDomainNotFoundFault:
				return result, errors.New(rds.ErrCodeDomainNotFoundFault + aerr.Error())
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

func RemoveClusterSnapshot(s string, sess *session.Session) (*rds.DeleteDBClusterSnapshotOutput, error) {
	var result *rds.DeleteDBClusterSnapshotOutput

	svc := rds.New(sess)
	input := &rds.DeleteDBClusterSnapshotInput{
		DBClusterSnapshotIdentifier: aws.String(s),
	}

	result, err := svc.DeleteDBClusterSnapshot(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeInvalidDBClusterSnapshotStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBClusterSnapshotStateFault + aerr.Error())
			case rds.ErrCodeDBClusterSnapshotNotFoundFault:
				return result, errors.New(rds.ErrCodeDBClusterSnapshotNotFoundFault + aerr.Error())
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

func CreateClusterInstance(n, t string, sess *session.Session) (*rds.CreateDBInstanceOutput, error) {
	var result *rds.CreateDBInstanceOutput

	svc := rds.New(sess)
	input := &rds.CreateDBInstanceInput{
		DBClusterIdentifier:  aws.String(DestinationClusterName),
		DBInstanceClass:      aws.String(t),
		DBInstanceIdentifier: aws.String(n),
		Engine:               aws.String(DestinationClusterEngine),
	}

	result, err := svc.CreateDBInstance(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBInstanceAlreadyExistsFault:
				return result, errors.New(rds.ErrCodeDBInstanceAlreadyExistsFault + aerr.Error())
			case rds.ErrCodeInsufficientDBInstanceCapacityFault:
				return result, errors.New(rds.ErrCodeInsufficientDBInstanceCapacityFault + aerr.Error())
			case rds.ErrCodeDBParameterGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeDBParameterGroupNotFoundFault + aerr.Error())
			case rds.ErrCodeDBSecurityGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeDBSecurityGroupNotFoundFault + aerr.Error())
			case rds.ErrCodeInstanceQuotaExceededFault:
				return result, errors.New(rds.ErrCodeInstanceQuotaExceededFault + aerr.Error())
			case rds.ErrCodeStorageQuotaExceededFault:
				return result, errors.New(rds.ErrCodeStorageQuotaExceededFault + aerr.Error())
			case rds.ErrCodeDBSubnetGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeDBSubnetGroupNotFoundFault + aerr.Error())
			case rds.ErrCodeDBSubnetGroupDoesNotCoverEnoughAZs:
				return result, errors.New(rds.ErrCodeDBSubnetGroupDoesNotCoverEnoughAZs + aerr.Error())
			case rds.ErrCodeInvalidDBClusterStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBClusterStateFault + aerr.Error())
			case rds.ErrCodeInvalidSubnet:
				return result, errors.New(rds.ErrCodeInvalidSubnet + aerr.Error())
			case rds.ErrCodeInvalidVPCNetworkStateFault:
				return result, errors.New(rds.ErrCodeInvalidVPCNetworkStateFault + aerr.Error())
			case rds.ErrCodeProvisionedIopsNotAvailableInAZFault:
				return result, errors.New(rds.ErrCodeProvisionedIopsNotAvailableInAZFault + aerr.Error())
			case rds.ErrCodeOptionGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeOptionGroupNotFoundFault + aerr.Error())
			case rds.ErrCodeDBClusterNotFoundFault:
				return result, errors.New(rds.ErrCodeDBClusterNotFoundFault + aerr.Error())
			case rds.ErrCodeStorageTypeNotSupportedFault:
				return result, errors.New(rds.ErrCodeStorageTypeNotSupportedFault + aerr.Error())
			case rds.ErrCodeAuthorizationNotFoundFault:
				return result, errors.New(rds.ErrCodeAuthorizationNotFoundFault + aerr.Error())
			case rds.ErrCodeKMSKeyNotAccessibleFault:
				return result, errors.New(rds.ErrCodeKMSKeyNotAccessibleFault + aerr.Error())
			case rds.ErrCodeDomainNotFoundFault:
				return result, errors.New(rds.ErrCodeDomainNotFoundFault + aerr.Error())
			case rds.ErrCodeBackupPolicyNotFoundFault:
				return result, errors.New(rds.ErrCodeBackupPolicyNotFoundFault + aerr.Error())
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

func SetClusterInstance(i string, sess *session.Session) (*rds.ModifyDBInstanceOutput, error) {
	var result *rds.ModifyDBInstanceOutput

	svc := rds.New(sess)
	input := &rds.ModifyDBInstanceInput{
		ApplyImmediately:           aws.Bool(true),
		BackupRetentionPeriod:      aws.Int64(10),
		DBInstanceIdentifier:       aws.String(i),
		PreferredBackupWindow:      aws.String("04:00-04:30"),
		PreferredMaintenanceWindow: aws.String("Tue:05:00-Tue:05:30"),
	}

	result, err := svc.ModifyDBInstance(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeInvalidDBInstanceStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBInstanceStateFault + aerr.Error())
			case rds.ErrCodeInvalidDBSecurityGroupStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBSecurityGroupStateFault + aerr.Error())
			case rds.ErrCodeDBInstanceAlreadyExistsFault:
				return result, errors.New(rds.ErrCodeDBInstanceAlreadyExistsFault + aerr.Error())
			case rds.ErrCodeDBInstanceNotFoundFault:
				return result, errors.New(rds.ErrCodeDBInstanceNotFoundFault + aerr.Error())
			case rds.ErrCodeDBSecurityGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeDBSecurityGroupNotFoundFault + aerr.Error())
			case rds.ErrCodeDBParameterGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeDBParameterGroupNotFoundFault + aerr.Error())
			case rds.ErrCodeInsufficientDBInstanceCapacityFault:
				return result, errors.New(rds.ErrCodeInsufficientDBInstanceCapacityFault + aerr.Error())
			case rds.ErrCodeStorageQuotaExceededFault:
				return result, errors.New(rds.ErrCodeStorageQuotaExceededFault + aerr.Error())
			case rds.ErrCodeInvalidVPCNetworkStateFault:
				return result, errors.New(rds.ErrCodeInvalidVPCNetworkStateFault + aerr.Error())
			case rds.ErrCodeProvisionedIopsNotAvailableInAZFault:
				return result, errors.New(rds.ErrCodeProvisionedIopsNotAvailableInAZFault + aerr.Error())
			case rds.ErrCodeOptionGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeOptionGroupNotFoundFault + aerr.Error())
			case rds.ErrCodeDBUpgradeDependencyFailureFault:
				return result, errors.New(rds.ErrCodeDBUpgradeDependencyFailureFault + aerr.Error())
			case rds.ErrCodeStorageTypeNotSupportedFault:
				return result, errors.New(rds.ErrCodeStorageTypeNotSupportedFault + aerr.Error())
			case rds.ErrCodeAuthorizationNotFoundFault:
				return result, errors.New(rds.ErrCodeAuthorizationNotFoundFault + aerr.Error())
			case rds.ErrCodeCertificateNotFoundFault:
				return result, errors.New(rds.ErrCodeCertificateNotFoundFault + aerr.Error())
			case rds.ErrCodeDomainNotFoundFault:
				return result, errors.New(rds.ErrCodeDomainNotFoundFault + aerr.Error())
			case rds.ErrCodeBackupPolicyNotFoundFault:
				return result, errors.New(rds.ErrCodeBackupPolicyNotFoundFault + aerr.Error())
			case rds.ErrCodeKMSKeyNotAccessibleFault:
				return result, errors.New(rds.ErrCodeKMSKeyNotAccessibleFault + aerr.Error())
			case rds.ErrCodeInvalidDBClusterStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBClusterStateFault + aerr.Error())
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

func CreateClusterInstanceReadReplica(sess *session.Session) (*rds.CreateDBInstanceReadReplicaOutput, error) {
	var result *rds.CreateDBInstanceReadReplicaOutput

	svc := rds.New(sess)
	input := &rds.CreateDBInstanceReadReplicaInput{
		CopyTagsToSnapshot:         aws.Bool(true),
		DBInstanceClass:            aws.String(DestinationWriterInstanceType),
		DBInstanceIdentifier:       aws.String(DestinationClusterReaderInstanceName),
		PubliclyAccessible:         aws.Bool(true),
		SourceDBInstanceIdentifier: aws.String(DestinationClusterWriterInstanceName),
	}

	result, err := svc.CreateDBInstanceReadReplica(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBInstanceAlreadyExistsFault:
				return result, errors.New(rds.ErrCodeDBInstanceAlreadyExistsFault + aerr.Error())
			case rds.ErrCodeInsufficientDBInstanceCapacityFault:
				return result, errors.New(rds.ErrCodeInsufficientDBInstanceCapacityFault + aerr.Error())
			case rds.ErrCodeDBParameterGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeDBParameterGroupNotFoundFault + aerr.Error())
			case rds.ErrCodeDBSecurityGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeDBSecurityGroupNotFoundFault + aerr.Error())
			case rds.ErrCodeInstanceQuotaExceededFault:
				return result, errors.New(rds.ErrCodeInstanceQuotaExceededFault + aerr.Error())
			case rds.ErrCodeStorageQuotaExceededFault:
				return result, errors.New(rds.ErrCodeStorageQuotaExceededFault + aerr.Error())
			case rds.ErrCodeDBInstanceNotFoundFault:
				return result, errors.New(rds.ErrCodeDBInstanceNotFoundFault + aerr.Error())
			case rds.ErrCodeInvalidDBInstanceStateFault:
				return result, errors.New(rds.ErrCodeInvalidDBInstanceStateFault + aerr.Error())
			case rds.ErrCodeDBSubnetGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeDBSubnetGroupNotFoundFault + aerr.Error())
			case rds.ErrCodeDBSubnetGroupDoesNotCoverEnoughAZs:
				return result, errors.New(rds.ErrCodeDBSubnetGroupDoesNotCoverEnoughAZs + aerr.Error())
			case rds.ErrCodeInvalidSubnet:
				return result, errors.New(rds.ErrCodeInvalidSubnet + aerr.Error())
			case rds.ErrCodeInvalidVPCNetworkStateFault:
				return result, errors.New(rds.ErrCodeInvalidVPCNetworkStateFault + aerr.Error())
			case rds.ErrCodeProvisionedIopsNotAvailableInAZFault:
				return result, errors.New(rds.ErrCodeProvisionedIopsNotAvailableInAZFault + aerr.Error())
			case rds.ErrCodeOptionGroupNotFoundFault:
				return result, errors.New(rds.ErrCodeOptionGroupNotFoundFault + aerr.Error())
			case rds.ErrCodeDBSubnetGroupNotAllowedFault:
				return result, errors.New(rds.ErrCodeDBSubnetGroupNotAllowedFault + aerr.Error())
			case rds.ErrCodeInvalidDBSubnetGroupFault:
				return result, errors.New(rds.ErrCodeInvalidDBSubnetGroupFault + aerr.Error())
			case rds.ErrCodeStorageTypeNotSupportedFault:
				return result, errors.New(rds.ErrCodeStorageTypeNotSupportedFault + aerr.Error())
			case rds.ErrCodeKMSKeyNotAccessibleFault:
				return result, errors.New(rds.ErrCodeKMSKeyNotAccessibleFault + aerr.Error())
			case rds.ErrCodeDomainNotFoundFault:
				return result, errors.New(rds.ErrCodeDomainNotFoundFault + aerr.Error())
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

func GetClusterInstance(n string, sess *session.Session) (*rds.DescribeDBInstancesOutput, error) {

	var result *rds.DescribeDBInstancesOutput

	svc := rds.New(sess)
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(n),
	}

	result, err := svc.DescribeDBInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBInstanceNotFoundFault:
				return result, errors.New(rds.ErrCodeDBInstanceNotFoundFault + aerr.Error())
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

func Log(m string) {
	log.Println(m)
}

var (
	SourceClusterName                    string
	DestinationClusterName               string
	MigrationKeyAlias                    string
	MigrationSnapshotARN                 string
	SourceProfile                        string
	SourceProfileRegion                  string = "eu-west-2"
	DestinationProfile                   string
	DestinationProfileRegion             string = "eu-west-2"
	DestinationKMSKeyAlias               string
	DestinationClusterWriterInstanceName string = "writer"
	DestinationClusterReaderInstanceName string = "reader"
	DestinationWriterInstanceType        string = "db.r5.2xlarge"
	DestinationReaderInstanceType        string = "db.r5.xlarge"
	DestinationClusterEngine             string
	DestinationClusterEngineVersion      string
	DestinationClusterEngineMode         string = "serverless"
	DestinationClusterSubnetGroup        string
	DestinationAccountID                 string
	DestinationClusterSecurityGroup      string
	ClusterAdministratorUserName         string = "admin"
)

var (
	ClusterSnapshotName     = "migrationsnapshot-" + time.Now().Format("11063912340")
	ClusterSnapshotCopyName = "migrationsnapshotshared-" + time.Now().Format("11063912340")
)

func main() {

	start := time.Now()
	var msg string
	var wg sync.WaitGroup

	flag.StringVar(&SourceClusterName, "SourceClusterName", SourceClusterName, "Specify the name of the cluster to migrate.")
	flag.StringVar(&DestinationClusterName, "DestinationClusterName", SourceClusterName, "Enter the name that will belong to the cluster created in the destination account.")
	flag.StringVar(&MigrationKeyAlias, "MigrationKeyAlias", MigrationKeyAlias, "The name of the key used to share the snapshot with the destination account.")
	flag.StringVar(&SourceProfile, "SourceProfile", SourceProfile, "The name of the profile with access to the source db cluster account")
	flag.StringVar(&DestinationProfile, "DestinationProfile", DestinationProfile, "The name of the profile with access to the destination db cluster account")
	flag.StringVar(&DestinationKMSKeyAlias, "DestinationKMSKeyAlias", DestinationKMSKeyAlias, "The alias of the key that will be used to encrypt the db cluster in the destination account")
	flag.StringVar(&DestinationClusterWriterInstanceName, "DestinationClusterWriterInstanceName", DestinationClusterWriterInstanceName, "The name of the  writer instnace that will be part of the migrated cluster in the destination account")
	flag.StringVar(&DestinationClusterReaderInstanceName, "DestinationClusterReaderInstanceName", DestinationClusterReaderInstanceName, "The name of the reader instnace that will be part of the migrated cluster in the destination account")
	flag.StringVar(&DestinationWriterInstanceType, "DestinationWriterInstanceType", DestinationWriterInstanceType, "The instance type of the db cluster writer instance in the destination account")
	flag.StringVar(&DestinationReaderInstanceType, "DestinationReaderInstanceType", DestinationReaderInstanceType, "The instance type of the db cluster reader instances in the destination account")
	flag.StringVar(&DestinationAccountID, "DestinationAccountID", DestinationAccountID, "The ID of the account where the db will be migrated")
	flag.StringVar(&DestinationClusterEngine, "DestinationClusterEngine", DestinationClusterEngine, "The destination cluster engine version")
	flag.StringVar(&DestinationClusterEngineMode, "DestinationClusterEngineMode", DestinationClusterEngineMode, "The destination cluster engine mode")
	flag.StringVar(&DestinationClusterEngineVersion, "DestinationClusterEngineVersion", DestinationClusterEngineVersion, "The destination cluster engine version")
	flag.StringVar(&DestinationClusterSubnetGroup, "DestinationClusterSubnetGroup", DestinationClusterSubnetGroup, "The VPC rds subnets group where the cluster should be placed")
	flag.StringVar(&DestinationClusterSecurityGroup, "DestinationClusterSecurityGroup", DestinationClusterSecurityGroup, "The security group to be assosiated with the destination cluster")
	flag.StringVar(&ClusterAdministratorUserName, "ClusterAdministratorUserName", ClusterAdministratorUserName, "The admin user name of the db cluster that will be migrated")
	flag.StringVar(&SourceProfileRegion, "SourceProfileRegion", SourceProfileRegion, "Specify the region where the db is located.")
	flag.StringVar(&DestinationProfileRegion, "DestinationProfileRegion", DestinationProfileRegion, "Specify the region where the db is going to be migrated.")

	flag.Parse()

	SourceSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           SourceProfile,
		Config: aws.Config{
			Region: aws.String(SourceProfileRegion),
		},
	}))

	DestinationSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           DestinationProfile,
		Config: aws.Config{
			Region: aws.String(DestinationProfileRegion),
		},
	}))

	// Create cluster snapshot from source cluster
	Log("Creating db cluster snapshot: " + ClusterSnapshotName)
	_, err := CreateClusterSnapshot(SourceClusterName, ClusterSnapshotName, SourceSession)
	if err != nil {
		log.Fatal(err)
	}
	Log("Wait until Snapshot is completed...")
	for status := false; !status; {
		time.Sleep(1 * time.Minute)
		result, err := GetClusterSnapshot(ClusterSnapshotName, "", SourceSession)
		if err != nil {
			log.Fatal(err)
		}
		if *result.DBClusterSnapshots[0].Status == "available" {
			status = true
		}
	}
	Log("Cluster snapshot successfully created")

	Log("Copying snapshot with new KMS key: " + MigrationKeyAlias)
	_, err = CopyClusterSnapshot(ClusterSnapshotName, ClusterSnapshotCopyName, MigrationKeyAlias, SourceSession)
	if err != nil {
		log.Fatal(err)
	}

	Log("Wait until Snapshot is completed...")
	for status := false; !status; {
		time.Sleep(1 * time.Minute)
		result, err := GetClusterSnapshot(ClusterSnapshotCopyName, "", SourceSession)
		if err != nil {
			log.Fatal(err)
		}
		if *result.DBClusterSnapshots[0].Status == "available" {
			status = true
		}
	}

	Log("Cluster snapshot copy successfully created")

	_, err = RemoveClusterSnapshot(ClusterSnapshotName, SourceSession)
	if err != nil {
		log.Fatal(err)
	}

	Log("Sharing snapshot with destination account: " + DestinationAccountID)
	_, err = ShareClusterSnapshot(ClusterSnapshotCopyName, DestinationAccountID, SourceSession)
	if err != nil {
		log.Fatal(err)
	}

	// Get shared snapshot

	s, err := GetClusterSnapshot(ClusterSnapshotCopyName, "", SourceSession)
	if err != nil {
		log.Fatal(err)
	}
	MigrationSnapshotARN = *s.DBClusterSnapshots[0].DBClusterSnapshotArn

	Log("Creating cluster " + DestinationClusterName + " in destination account " + DestinationAccountID)
	_, err = CreateClusterFromSnapshot(DestinationSession)
	if err != nil {
		log.Fatal(err)
	}

	Log("Wait untill cluster is ready...")
	for status := false; !status; {
		time.Sleep(1 * time.Minute)
		result, err := GetCluster(DestinationClusterName, DestinationSession)
		if err != nil {
			log.Fatal(err)
		}

		if *result.DBClusters[0].Status == "available" {
			status = true
		}
	}

	Log("Cluster " + DestinationClusterName + " successfully created")

	_, err = RemoveClusterSnapshot(ClusterSnapshotCopyName, SourceSession)
	if err != nil {
		log.Fatal(err)
	}

	msg = "ready"

	if msg == "ready" && DestinationClusterEngineMode != "serverless" {
		wg.Add(2)
		go func() {

			Log("Creating Writer instances: " + DestinationClusterWriterInstanceName)
			_, err := CreateClusterInstance(DestinationClusterWriterInstanceName, DestinationWriterInstanceType, DestinationSession)
			if err != nil {
				log.Fatal(err)
			}

			Log("Wait for Writer instance to be ready...")
			for status := false; !status; {
				time.Sleep(1 * time.Minute)
				result, err := GetClusterInstance(DestinationClusterWriterInstanceName, DestinationSession)
				if err != nil {
					log.Fatal(err)
				}
				if *result.DBInstances[0].DBInstanceStatus == "available" {
					status = true
				}
			}

			Log("Writer instances successfully created")
			defer wg.Done()
		}()

		go func() {
			time.Sleep(5 * time.Millisecond)
			Log("Creating Reader instance: " + DestinationClusterReaderInstanceName)
			_, err := CreateClusterInstance(DestinationClusterReaderInstanceName, DestinationReaderInstanceType, DestinationSession)
			if err != nil {
				log.Fatal(err)
			}
			Log("Wait for Reader instance to be ready...")
			for status := false; !status; {
				time.Sleep(1 * time.Minute)
				result, err := GetClusterInstance(DestinationClusterReaderInstanceName, DestinationSession)
				if err != nil {
					log.Fatal(err)
				}
				if *result.DBInstances[0].DBInstanceStatus == "available" {
					status = true
				}
			}

			Log("Reader instances successfully created")

			defer wg.Done()
		}()

	}

	wg.Wait()
	Log("Migration Completed")
	Log("Total migration time: " + time.Since(start).String())
}

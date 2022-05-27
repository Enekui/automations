# Seqera Labs AWS Database Cross-Account Migration Tool
### This tool will take care of migrating a database that has been created using the default AWS kms/rds key. 
### Once the database is created with a default AWS managed kms key, the snapshots created from this DB can't be shared with another account. That happens because the AWS managed KMS keys policies can't be updated.

## This script will perform the follwoing steps:
1. Create an temporary spapshot of the DB that we want to migrate.
2. Restore the previously created snapshot into a new temporary db using a customer managed KMS key that has been previously shared with the destination account.
3. Once that db is created with the new customer managed kms key, another snapshot is created from this db and shared with the destination account.
4. After the snapshot has been properly shared with the destination account, this snapshot is restored in the destination account using a kms key specified in the script parameters that belongs to the destination account.
5. Once the db is ready in the destination account, the script will create a db cluster instance.

### Note: The script will take care of cleaning up all the temporary created resources.

## The script can be invoked with the following parameters:
### --destination-account-id string
        Enter the destination account ID where the DB is going to be migrated
### --destination-cluster string
        Enter the name of the DB cluster in the destination account (default "migrated-cluster")
### --destination-instance-name string
        Enter the name of the DB cluster instance in the migrated account (default "migrated-cluster-instance")
### --destination-kms-key-id string
        Enter the KMS key ID to encrypt the DB cluster in the destination account
### --destination-profile string
        Enter the profile name for the destination account
### --source-cluster string
        Enter the name of the source db cluster to be migrated
### --source-profile string
        Enter the profile name to connect to the source DB account (default "default")
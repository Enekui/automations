#!/bin/bash

# Production Gitea
go run db-migration-v2.go --DestinationAccountID=<account ID>  --DestinationClusterName="gitea" --DestinationProfile="production" \
 --DestinationClusterSecurityGroup=<SG ID> --ClusterAdministratorUserName="admin" \
 --SourceClusterName="gitea" --DestinationClusterEngine="aurora" --MigrationKeyAlias="rds/migration" --SourceProfile="default" \
 --DestinationClusterEngineVersion="5.6.mysql_aurora.1.23.4" --DestinationClusterSubnetGroup="rds_subnet_group" --DestinationKMSKeyAlias="rds" \
 --SourceProfileRegion="eu-west-2" --DestinationProfileRegion="eu-west-2"

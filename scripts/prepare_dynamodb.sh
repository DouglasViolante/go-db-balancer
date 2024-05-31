#!/bin/bash

LSTACK_URL="http://localhost:4566"

echo "Setting AWS Environment Config"
aws configure set default.region 		us-east-1
aws configure set aws_access_key_id 	1234
aws configure set aws_secret_access_key 4321

echo "Checking for Existing Tables..."
for tableName in $(aws --endpoint $LSTACK_URL dynamodb list-tables --output text --query 'TableNames[]'); do 
	echo "Now Deleting $tableName"
    aws --endpoint $LSTACK_URL dynamodb delete-table --table-name $tableName
	echo "$tableName Purged!"
	
	echo "Waiting $tableName Purge to be Consolidated!"
	aws --endpoint $LSTACK_URL dynamodb wait table-not-exists --table-name $tableName
	echo "Done!"
done

tableNames=("db_A" "db_B")

for tableName in "${tableNames[@]}"; do
	echo "Creating $tableName"
    aws --endpoint $LSTACK_URL dynamodb create-table \
        --table-name $tableName \
        --attribute-definitions \
				AttributeName=ID,AttributeType=N \
        --key-schema \
				AttributeName=ID,KeyType=HASH \
        --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
		--tags Key=DbType,Value=DynamoDB
	echo "$tableName has been Created!"
	
	echo "Waiting $tableName to be Consolidated!"
	aws --endpoint $LSTACK_URL dynamodb wait table-exists --table-name $tableName
	echo "Done!"
done

echo "Loading 30 Data Items into ${tableNames[0]}"
for i in $(seq 1 30); do
	randomID=$(shuf -i 1-1000000 -n 1) 
	randomPW=$(shuf -i 0-5000000 -n 1) 
	randomCT=$(shuf -i 1-9000000 -n 1)
	
	echo "Inserting $randomID ..."
	aws --endpoint $LSTACK_URL dynamodb put-item \
		--table-name ${tableNames[0]} \
		--item '{"ID": 		{"N": "'$randomID'"},
				 "Power": 	{"N": "'$randomPW'"},
				 "Cost": 	{"N": "'$randomCT'"}}'
	echo "$randomID Sucessfully Inserted!"
done
echo "Data Loading on ${tableNames[0]} Completed!"


echo "Loading 15 Data Items into ${tableNames[1]}"
for i in $(seq 1 15); do
	randomID=$(shuf -i 1-1000000 -n 1) 
	randomPW=$(shuf -i 0-5000000 -n 1) 
	randomCT=$(shuf -i 1-9000000 -n 1)
	
	echo "Inserting $randomID ..."
	aws --endpoint $LSTACK_URL dynamodb put-item \
		--table-name ${tableNames[1]} \
		--item '{"ID": 		{"N": "'$randomID'"},
				 "Power": 	{"N": "'$randomPW'"},
				 "Cost": 	{"N": "'$randomCT'"}}'
	echo "$randomID Sucessfully Inserted!"
done
echo "Data Loading on ${tableNames[1]} Completed!"



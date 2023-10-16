#!/bin/bash
echo "########### SCRIPT 01 ###########"
echo "########### Creating table transfer-engine ###########"
aws dynamodb --endpoint-url=http://localhost:4566 create-table \
    --table-name users \
    --stream-specification StreamEnabled=true,StreamViewType=NEW_IMAGE \
    --attribute-definitions \
        AttributeName=email,AttributeType=S \
        AttributeName=user_id,AttributeType=S \
    --key-schema \
        AttributeName=email,KeyType=HASH \
    --provisioned-throughput \
    	      ReadCapacityUnits=100,WriteCapacityUnits=100 \
    --global-secondary-indexes \
            "[
                {
                    \"IndexName\": \"user_id_index\",
                    \"KeySchema\": [{\"AttributeName\":\"user_id\",\"KeyType\":\"HASH\"}],
                    \"Projection\":{
                        \"ProjectionType\":\"ALL\"
                    },
                    \"ProvisionedThroughput\": {
                        \"ReadCapacityUnits\": 10,
                        \"WriteCapacityUnits\": 5
                    }
                }
            ]"
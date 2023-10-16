#!/bin/bash
echo "########### SCRIPT 00 ###########"
echo pwd
echo "########### Setting aws_access_key_id profile ###########"
aws --endpoint-url=http://localhost:4566 configure set aws_access_key_id test --profile=default
echo "########### Setting aws_secret_access_key ###########"
aws --endpoint-url=http://localhost:4566 configure set aws_secret_access_key test --profile=default
echo "########### Setting region ###########"
aws --endpoint-url=http://localhost:4566 configure set region us-east-1 --profile=default
#!/bin/bash

P_USER=$1

if [[ -z "${P_USER}" ]]; then
  echo "you should set an user"
  exit 1
fi

aws s3api put-object --bucket handsonbartender --key $P_USER/swagger.yaml --body swagger.yaml --profile workshop

aws cloudformation package --template-file sam.yml --s3-bucket handsonbartender --output-template-file packaged.yml --profile workshop

## see https://docs.aws.amazon.com/cli/latest/reference/cloudformation/package.html 

echo "deploying...."

aws cloudformation deploy --template-file packaged.yml --stack-name $P_USER-bartender-sam-deploy --parameter-overrides User=$P_USER --capabilities CAPABILITY_IAM --profile workshop

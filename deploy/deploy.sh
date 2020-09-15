#!/bin/bash

# Optionally, set default values
# removeServer="default value for 127.0.0.1"

. deploy/deploy.config

#go build -o ../build/officetime cmd/officetime/main.go
env GOOS=linux GOARCH=amd64 go build -o ../build/o_api -v main.go
ssh root@$removeServer 'systemctl stop officetime-api.service && rm /var/officetime/o_api'
scp ../build/o_api root@$removeServer:/var/officetime/
ssh root@$removeServer 'systemctl start officetime-api.service'

exit 0
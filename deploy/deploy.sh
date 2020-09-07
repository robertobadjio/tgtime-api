#!/bin/bash

# Optionally, set default values
# removeServer="default value for 127.0.0.1"

. deploy/deploy.config

#go build -o ../build/officetime cmd/officetime/main.go
env GOOS=linux GOARCH=amd64 go build -o ../build/officetime_api -v main.go
ssh root@$removeServer 'pkill officetime_api; rm /var/officetime/officetime_api'
scp ../build/officetime_api root@$removeServer:/var/officetime/
ssh root@$removeServer 'cd /var/officetime && ./officetime_api -config="./config"'

exit 0
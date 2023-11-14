#!/bin/bash

########################################################

## Shell Script to Build Docker Image Insightful

########################################################

DATE=`date +%Y.%m.%d.%H.%M.%S`
git checkout develop
git pull origin develop
docker build -t insightful:$DATE -f /src/build/Dockerfile .
docker run -it --network=host -p 8899:8899 insightful:$DATE
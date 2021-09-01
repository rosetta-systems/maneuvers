#!/usr/bin/env bash

container=`docker ps --format {{.Names}}`
containerid=`docker container inspect $container -f {{.Image}} | cut -d ":" -f 2`
matchid=${containerid:0:12}

docker images --format {{.ID}} | while IFS= read line; do
if [[ $line == $matchid ]]; then
  


#WAIT, JUST GET IMAGE ID FROM WEBHOOK

#!/bin/bash

app=mesh
ssh -t ajmoon "
  if [[ ! -d '/opt/$app' ]]; then
    mkdir /opt/$app
    cd /opt/$app
    git clone https://github.com/alex-moon/$app .
  fi
  cd /opt/$app
  git reset --hard HEAD
  git pull origin \$(git rev-parse --abbrev-ref HEAD)
  docker compose build
  docker stack deploy -c docker-compose.yml $app
  docker run --rm -i \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -u root ubirak/docker-php:latest \
    stack:converge $app
  docker system prune -f
"

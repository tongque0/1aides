#!/bin/bash

# 检查并创建 Docker 网络，如果网络已存在则不创建
docker network ls | grep aides-network > /dev/null || docker network create aides-network

# 删除已存在的 mongo 容器，如果存在的话
docker ps -a | grep -q mongo && docker rm -f mongo

# 启动 mongo 服务
docker run -d --name mongo \
  -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=aides \
  -e MONGO_INITDB_ROOT_PASSWORD=dGhpcyBpcyBhaWRlcw== \
  --restart always \
  --network aides-network \
  --volume 1aides_mongo:/data/db \
  hub.atomgit.com/amd64/mongo:latest

# 删除已存在的 1aides 容器，如果存在的话
docker ps -a | grep -q 1aides && docker rm -f 1aides

# 启动 1aides 服务
docker run -d --name 1aides \
  -p 8999:8999 \
  -e TEMPLATE_PATH=/app/frontend/templates/* \
  -e STATIC_PATH=/app/frontend/static \
  -e MONGO_USER=aides \
  -e MONGO_PASSWORD=dGhpcyBpcyBhaWRlcw== \
  -e MONGO_HOST=mongo \
  --restart always \
  --network aides-network \
  --volume 1aides_logs:/logs \
  serverless-100026543835-docker.pkg.coding.net/1aides/1aides/1aides:latest

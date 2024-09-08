#!/bin/bash

# 设置脚本在遇到错误时停止执行
set -e

# 检查并创建 Docker 网络，如果网络已存在则不创建
if ! docker network ls | grep -q aides-network; then
  echo "Creating Docker network: aides-network"
  docker network create aides-network
else
  echo "Docker network 'aides-network' already exists."
fi

# 删除已存在的 mongo 容器，如果存在的话
if docker ps -a | grep -q mongo; then
  echo "Removing existing mongo container..."
  docker rm -f mongo
fi

echo "Starting mongo service..."
# 启动 mongo 服务
docker run -d --name mongo \
  -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=aides \
  -e MONGO_INITDB_ROOT_PASSWORD=dGhpcyBpcyBhaWRlcw== \
  --restart always \
  --network aides-network \
  --volume 1aides_mongo:/data/db \
  hub.atomgit.com/amd64/mongo:latest

# 删除已存在的 1aides 容器和镜像，如果存在的话
if docker ps -a | grep -q 1aides; then
  echo "Removing existing 1aides container..."
  docker rm -f 1aides
fi

if docker images | grep -q 'serverless-100026543835-docker.pkg.coding.net/1aides/1aides/1aides'; then
  echo "Removing old 1aides image..."
  docker rmi -f serverless-100026543835-docker.pkg.coding.net/1aides/1aides/1aides:latest
fi

echo "Pulling the latest 1aides image..."
# 拉取最新的 1aides 镜像
docker pull serverless-100026543835-docker.pkg.coding.net/1aides/1aides/1aides:latest

echo "Starting 1aides service..."
# 启动 1aides 服务
docker run -d --name 1aides \
  -p 8999:8999 \
  -e TEMPLATE_PATH=/app/frontend/templates/* \
  -e STATIC_PATH=/app/frontend/static \
  -e MONGO_USER=aides \
  -e MONGO_PASSWORD=dGhpcyBpcyBhaWRlcw== \
  -e MONGO_HOST=mongo \
  -e GIN_MODE=release \
  --restart always \
  --network aides-network \
  --volume 1aides_logs:/logs \
  serverless-100026543835-docker.pkg.coding.net/1aides/1aides/1aides:latest

echo "Services are up and running!"

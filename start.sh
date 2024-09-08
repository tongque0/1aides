#!/bin/bash

# 设置脚本在遇到错误时停止执行
set -e

# 显示启动信息
echo "   ______    _     __"
echo "  <  /   |  (_)___/ /__  _____"
echo "  / / /| | / / __  / _ \/ ___/"
echo " / / ___ |/ / /_/ /  __(__  )"
echo "/_/_/  |_/_/\__,_/\___/____/"
echo "正在初始化 1Aides 服务..."
echo ""

# 检查并创建 Docker 网络，如果网络已存在则不创建
if ! docker network ls | grep -q aides-network; then
  echo -e "\e[1;32m正在创建 Docker 网络：aides-network\e[0m"
  docker network create aides-network
else
  echo -e "\e[1;33mDocker 网络 'aides-network' 已经存在。\e[0m"
fi

# 检查是否存在 aides-mongo 容器，如果不存在则创建
if ! docker ps -a | grep -q aides-mongo; then
  echo -e "\e[1;32m正在启动 aides-mongo 服务...\e[0m"
  docker run -d --name aides-mongo \
    -p 27017:27017 \
    -e MONGO_INITDB_ROOT_USERNAME=aides \
    -e MONGO_INITDB_ROOT_PASSWORD=dGhpcyBpcyBhaWRlcw== \
    --restart always \
    --network aides-network \
    --volume 1aides_mongo:/data/db \
    hub.atomgit.com/amd64/mongo:latest
else
  echo -e "\e[1;33m已存在的 aides-mongo 容器正在运行。\e[0m"
fi

# 删除已存在的 1aides 容器和镜像，如果存在的话
if docker ps -a | grep -q 1aides; then
  echo -e "\e[1;31m正在移除已存在的 1aides 容器...\e[0m"
  docker rm -f 1aides
fi

if docker images | grep -q 'serverless-100026543835-docker.pkg.coding.net/1aides/1aides/1aides'; then
  echo -e "\e[1;31m正在移除旧的 1aides 镜像...\e[0m"
  docker rmi -f serverless-100026543835-docker.pkg.coding.net/1aides/1aides/1aides:latest
fi

echo -e "\e[1;32m正在拉取最新的 1aides 镜像...\e[0m"
# 拉取最新的 1aides 镜像
docker pull serverless-100026543835-docker.pkg.coding.net/1aides/1aides/1aides:latest

echo -e "\e[1;32m正在启动 1aides 服务...\e[0m"
# 启动 1aides 服务
docker run -d --name 1aides \
  -p 8999:8999 \
  -e TEMPLATE_PATH=/app/frontend/templates/* \
  -e STATIC_PATH=/app/frontend/static \
  -e MONGO_USER=aides \
  -e MONGO_PASSWORD=dGhpcyBpcyBhaWRlcw== \
  -e MONGO_HOST=aides-mongo \
  -e GIN_MODE=release \
  --restart always \
  --network aides-network \
  --volume 1aides_logs:/logs \
  serverless-100026543835-docker.pkg.coding.net/1aides/1aides/1aides:latest

echo -e "\e[1;32m服务已启动并运行！\e[0m"
echo -e "\e[1;35m1Aides 服务成功启动！\e[0m"
echo ""

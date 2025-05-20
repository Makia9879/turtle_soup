#!/bin/bash

# 设置变量
APP_NAME="turtle-soup"
VERSION="1.0.0"

# 构建Docker镜像
echo "Building Docker image for ${APP_NAME}..."
docker build -t ${APP_NAME}:${VERSION} .

# 标记镜像
echo "Tagging image..."
docker tag ${APP_NAME}:${VERSION} ${APP_NAME}:${VERSION}

# 推送镜像到仓库
echo "Pushing image to registry..."
docker push ${APP_NAME}:${VERSION}

echo "Build and push completed successfully!"
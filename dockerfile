# 使用官方Go 1.23镜像作为构建环境
FROM golang:1.23 AS builder
ENV GOPROXY="https://goproxy.cn|https://mirrors.tencentyun.com/go/|https://mirrors.aliyun.com/goproxy/|direct"
# 设置工作目录
WORKDIR /app

# 复制go mod和sum文件
COPY go.mod go.sum ./

# 下载所有依赖
RUN go mod download

# 复制源代码到容器中
COPY . .

# 构建应用程序，确保切换到包含main.go的目录，这里假设它在cmd目录下
RUN cd cmd && CGO_ENABLED=0 GOOS=linux go build -o 1aides .

# 获取CA证书（使用Alpine作为基础）
FROM alpine:latest AS certs
RUN apk --update add ca-certificates

# 使用scratch作为最终基础镜像，创建一个更小的镜像
FROM scratch

# 从Alpine中复制CA证书到scratch镜像
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# 从builder阶段复制可执行文件到根目录并重新命名为1aides
COPY --from=builder /app/cmd/1aides /app/1aides
# 确保路径正确，复制frontend目录
COPY --from=builder /app/frontend/ /app/frontend/
# 设置环境变量以确保应用使用正确的CA证书路径
ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt
ENV SSL_CERT_DIR=/etc/ssl/certs

# 运行应用
CMD ["/app/1aides"]

# 使用官方 Golang 镜像作为基础镜像
FROM golang:1.20-alpine

# 设置工作目录
WORKDIR /app

# 将 Go Modules 文件复制到容器中
COPY go.mod go.sum ./

# 拉取依赖包
RUN go mod tidy

# 将整个项目复制到容器中的工作目录
COPY . .

# 编译 Go 项目
RUN go build -o server ./cmd/server

# 暴露应用运行的端口
EXPOSE 8080

# 启动应用
CMD ["./server"]

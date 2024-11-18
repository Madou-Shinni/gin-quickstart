# 基础镜像
FROM golang:1.23 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

COPY . .

RUN go mod tidy

RUN GOOS=linux GOARCH=amd64 go build -o ./bin/server ./cmd/main.go

RUN chmod +x ./bin/server

FROM alpine:latest

# 设置时区为上海
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

COPY --from=builder /app/bin /app
COPY --from=builder /app/configs /app/configs

WORKDIR /app

CMD ["./server","-c","./configs"]
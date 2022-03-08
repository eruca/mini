FROM golang:alpine AS builder
WORKDIR /goapp
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod tidy
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
# RUN apk --no-cache add gcc g++ make git
# WORKDIR /go/src/app
# COPY . .
# RUN go env -w GOPROXY=https://goproxy.cn,direct && go get ./...
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/goapp .

FROM alpine
# RUN apk --no-cache add ca-certificates
WORKDIR /goapp
COPY --from=builder /goapp/bin /goapp
EXPOSE 8080
CMD [ "./goapp" ]
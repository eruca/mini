go mod tidy
GOOS=linux go build -ldflags="-s -w" -o ./bin/goapp .

docker build -t eruca/goapp .
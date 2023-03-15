# Start from golang base image
FROM golang:1.18

RUN apt-get update -y
WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/article-management-sys cmd/article-management-sys/main.go
EXPOSE 8080
CMD [ "/go/bin/article-management-sys" ]

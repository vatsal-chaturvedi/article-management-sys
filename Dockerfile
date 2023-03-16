# Start from golang base image
FROM golang:1.18 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/article-management-sys cmd/article-management-sys/main.go

FROM scratch

EXPOSE 8080
# Copy the Pre-built binary file
COPY --from=builder /go/bin/article-management-sys main
COPY --from=builder /app/configs configs
CMD [ "/main" ]

FROM golang:1.24-alpine

RUN apk add --no-cache protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

COPY notification-service/go.mod notification-service/go.sum ./
COPY /proto gen
RUN #protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative .\gen\notification.proto

RUN go mod download

COPY notification-service .

RUN go build -o notification-service ./cmd/main.go
EXPOSE 8080
CMD ["./notification-service"]
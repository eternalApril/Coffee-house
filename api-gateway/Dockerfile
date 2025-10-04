FROM golang:1.24-alpine

RUN apk add --no-cache protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

COPY api-gateway/go.mod api-gateway/go.sum ./
COPY /proto gen
RUN #protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative .\gen\order.proto

RUN go mod download

COPY api-gateway .


RUN go build -o api-gateway ./cmd/main.go
EXPOSE 8080
CMD ["./api-gateway"]
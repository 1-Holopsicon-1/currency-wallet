FROM golang:1.23.2-alpine
LABEL authors="Holopsicon"
RUN apk add --no-cache protoc
RUN apk update && apk add --no-cache make protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
ENV PATH="/root/go/bin:$PATH"
WORKDIR /app/protos
COPY protos .
RUN  protoc -I proto proto/exchange/exchange.proto --go_out=./gen/ --go_opt=paths=source_relative --go-grpc_out=./gen/ --go-grpc_opt=paths=source_relative
RUN go mod download
WORKDIR /app/gw-rates
COPY gw-exchanger/go.mod .
RUN go mod download
COPY gw-exchanger .
EXPOSE 50051
CMD ["go", "run", "cmd/main.go", "-docker", "-migrate", "-startGrpc"]
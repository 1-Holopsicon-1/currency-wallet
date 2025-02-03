FROM golang:1.23.2-alpine
LABEL authors="Holopsicon"
RUN apk add --no-cache protoc
RUN apk update && apk add --no-cache make protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
ENV PATH="/go/bin:$PATH"
RUN export PATH
WORKDIR /app/protos
COPY protos .
RUN  protoc -I proto proto/exchange/exchange.proto --go_out=./gen/ --go_opt=paths=source_relative --go-grpc_out=./gen/ --go-grpc_opt=paths=source_relative
RUN go mod download
WORKDIR /app/gw-wallet
COPY gw-currency-wallet/go.mod .
RUN go mod download
COPY gw-currency-wallet .
RUN go install github.com/swaggo/swag/cmd/swag@latest
ENV PATH="/go/bin:$PATH"
RUN export PATH
RUN swag init -g ../../../cmd/main.go --dir ./internal/app/handler --parseDependency --output ./docs

EXPOSE 5000
CMD ["go", "run", "cmd/main.go", "-docker", "-migrate", "-start"]
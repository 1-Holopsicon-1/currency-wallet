package protoExchange

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	proto "proto-exchange/gen/exchange"
	"time"
)

func Connect() (proto.ExchangeServiceClient, error) {
	e := godotenv.Load(".env")
	if e != nil {
		panic("No env")
	}
	host := os.Getenv("proto_host")
	port := os.Getenv("proto_port")
	serverAddress := fmt.Sprintf("%s:%s", host, port)
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	//Пытался сделать так, но у меня никогда не мог достучаться сервер, так и не понял почему
	//conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("could not connect to server: %v", err)
		return nil, err
	}

	client := proto.NewExchangeServiceClient(conn)

	return client, nil
}

func GetAllExchange(client proto.ExchangeServiceClient) (*proto.ExchangeRatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := client.GetExchangeRates(ctx, &proto.Empty{})
	if err != nil {
		log.Printf("could not get exchange rates: %v\n", err)
		return nil, err
	}
	return response, nil
}

func GetExchangeRatesForCurrency(client proto.ExchangeServiceClient, from_currency, to_currency string) (*proto.ExchangeRateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := client.GetExchangeRateForCurrency(ctx, &proto.CurrencyRequest{
		FromCurrency: from_currency,
		ToCurrency:   to_currency,
	})
	if err != nil {
		log.Printf("could not get exchange rates: %v\n", err)
		return nil, err
	}
	return response, nil
}

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cicdteam/go-deribit"
	"github.com/cicdteam/go-deribit/models"
	"github.com/phoebetronic/orderbook-deribit/pkg/orderbook"
)

func main() {
	fmt.Println()
	fmt.Println("|===========================================|")
	fmt.Println("|      wss://www.deribit.com/ws/api/v2      |")
	fmt.Println("|===========================================|")
	fmt.Println()

	var err error

	var erc chan error
	var clo chan bool
	{
		erc = make(chan error)
		clo = make(chan bool)
	}

	var cli *deribit.Exchange
	{
		cli, err = deribit.NewExchange(false, erc, clo)
		if err != nil {
			log.Fatalf("Error creating connection: %s", err)
		}

		err = cli.Connect()
		if err != nil {
			log.Fatalf("Error connecting to exchange: %s", err)
		}

		defer cli.Close()
	}

	go func() {
		log.Fatalf("RPC error: %s", <-erc)
		clo <- true
	}()

	var obk *orderbook.Orderbook
	{
		obk = orderbook.New()
	}

	var msg chan *models.BookNotification
	{
		msg, err = cli.SubscribeBookGroup("ETH-PERPETUAL", "none", "10", "100ms")
		if err != nil {
			log.Fatalf("Error subscribing to the book: %s", err)
		}
	}

	for x := range msg {
		{
			err = obk.Middleware(x)
			if err != nil {
				panic(err)
			}
		}

		var byt []byte
		{
			byt, err = json.Marshal(obk)
			if err != nil {
				panic(err)
			}
		}

		{
			fmt.Printf("%s\n", byt)
		}
	}
}

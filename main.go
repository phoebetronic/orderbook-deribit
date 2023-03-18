package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cicdteam/go-deribit"
	"github.com/cicdteam/go-deribit/models"
	"github.com/phoebetronic/orderbook-deribit/pkg/orderbook"
)

const (
	brk = "\n"
)

func main() {
	fmt.Println()
	fmt.Println("|===========================================|")
	fmt.Println("|      wss://www.deribit.com/ws/api/v2      |")
	fmt.Println("|===========================================|")
	fmt.Println()

	var err error

	var pat *string
	{
		pat = flag.String("pat", "-", "file path for the JSON encoded orderbook stream log, - for stdout")
		flag.Parse()
	}

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

		if *pat == "-" {
			fmt.Printf("%s\n", byt)
		} else {
			var fil *os.File
			{
				fil, err = os.OpenFile(*pat, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					panic(err)
				}
			}

			_, err := fil.Write(append(byt, brk...))
			if err != nil {
				panic(err)
			}

			err = fil.Close()
			if err != nil {
				panic(err)
			}
		}
	}
}

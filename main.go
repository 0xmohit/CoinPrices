package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"
)

func init() {
	log.SetFlags(0)
}

func main() {
	currencyOptions := fmt.Sprintf("Currency for price and market cap data.\n   \t%v", listCurrencies())
	currency := flag.String("currency", "USD", currencyOptions)
	limit := flag.Uint("limit", 10, "Default number of currencies.")
	flag.Parse()

	if !isValidCurrency(*currency, Currencies) {
		log.Fatalf("Unknown currency: %v\n    %v", *currency, listCurrencies())
	}

	options := Options{
		Limit:    *limit,
		Currency: *currency,
	}

	for {
		out, err := fetchPrices(options)
		if err != nil {
			log.Fatal(err)
		}
		v := allCoinData{Currency: strings.ToLower(*currency)}
		if err := v.UnmarshalJson(out); err != nil {
			log.Fatal(err)
		}
		printData(v)
		time.Sleep(5 * time.Minute)
	}
}

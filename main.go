// CoinPrices is a utility to fetch cryptocurrency prices.
//
// Usage:
//
//     CoinPrices
//     CoinPrices [ -currency name ] [ -limit number ]
//
// By default, the prices are displayed in USD and the number of
// cryptocurrencies is limited to 10.
//
// Valid currencies are: USD, AUD, BRL, BGN, CAD, CHF, CLP, CNY, CZK, DKK,
// EUR, GBP, HKD, HUF, IDR, ILS, INR, JPY, KRW, MXN, MYR, NOK, NZD, PHP, PKR,
// PLN, RUB, SEK, SGD, THB, TRY, TWD, ZAR.
//
// The data is obtained from [CoinMarketCap
// API](https://coinmarketcap.com/api/) and the quotes are refreshed every
// 5 minutes.
//
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

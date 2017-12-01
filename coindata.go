package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

var Currencies = []string{
	"USD", "AUD", "BRL", "BGN", "CAD", "CHF", "CLP", "CNY",
	"CZK", "DKK", "EUR", "GBP", "HKD", "HUF", "IDR", "ILS",
	"INR", "JPY", "KRW", "MXN", "MYR", "NOK", "NZD", "PHP",
	"PKR", "PLN", "RUB", "SEK", "SGD", "THB", "TRY", "TWD",
	"ZAR",
}

const apiUrl string = "https://api.coinmarketcap.com/v1/ticker/"

type coinData struct {
	Symbol     string `json:"symbol"`
	Change_1h  string `json:"percent_change_1h"`
	Change_24h string `json:"percent_change_24h"`
	Change_7d  string `json:"percent_change_7d"`
	Price      string
	MarketCap  string
}
type allCoinData struct {
	coin     []coinData
	Currency string
}

type Options struct {
	Limit    uint
	Currency string
}

func (d *allCoinData) UnmarshalJson(data []byte) error {
	var coinDataMap []map[string]interface{}
	if d == nil {
		return errors.New("UnmarshalJson on nil pointer")
	}
	if err := json.Unmarshal(data, &coinDataMap); err != nil {
		return err
	}
	for _, coin := range coinDataMap {
		var c coinData
		fmt.Printf("%+v\n", coin)
		c.Symbol = coin["symbol"].(string)
		c.Change_1h = getValue(coin["percent_change_1h"])
		c.Change_24h = getValue(coin["percent_change_24h"])
		c.Change_7d = getValue(coin["percent_change_7d"])
		c.Price = getValue(coin["price_"+d.Currency])
		c.MarketCap = getValue(coin["market_cap_"+d.Currency])
		(*d).coin = append((*d).coin, c)
	}
	return nil
}

func getValue(i interface{}) string {
	if v, ok := i.(string); ok {
		return v
	}
	return ""
}

// fetchPrices obtains the price data using the CoinMarketCap API.
func fetchPrices(o Options) ([]byte, error) {
	url := fmt.Sprintf("%s?limit=%v&convert=%s", apiUrl, o.Limit, o.Currency)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error fetching data from %v: %v", apiUrl, err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response from %v: %v", apiUrl, err)
	}
	return b, nil
}

func printData(d allCoinData) {
	headers := []string{
		bold("  Coin"),
		bold(fmt.Sprintf("   Price (%s)", strings.ToUpper(d.Currency))),
		bold(" Change (1H)"),
		bold("Change (24H)"),
		bold(fmt.Sprintf("Market Cap (%s)", strings.ToUpper(d.Currency))),
	}
	fmt.Printf("\033[H\033[2J")
	fmt.Printf("Price data as on %v\n", time.Now().Format("2006/01/02 03:04 PM"))
	fmt.Println(strings.Repeat("-", 76))
	fmt.Printf("| %s |\n", strings.Join(headers, " | "))
	fmt.Println(strings.Repeat("-", 76))
	for _, coin := range d.coin {
		records := []string{
			fmt.Sprintf("%6s", coin.Symbol),
			fmt.Sprintf("%14s", coin.Price),
			colorize(coin.Change_1h),
			colorize(coin.Change_24h),
			fmt.Sprintf("%16s", human(coin.MarketCap)),
		}
		fmt.Printf("| %s |\n", strings.Join(records, " | "))
		fmt.Println(strings.Repeat("-", 76))
	}
}

// isValidCurrency determines if the specified currency is valid.
// It ignores the case while performing the comparison.
func isValidCurrency(currency string, allCurrencies []string) bool {
	for _, c := range allCurrencies {
		if c == strings.ToUpper(currency) {
			return true
		}
	}
	return false
}

// listCurrencies returns the list of valid currencies for coversion.
func listCurrencies() string {
	options := strings.Join(Currencies, ", ")
	return fmt.Sprintf("Valid currencies are: %v", options)
}

// colorize returns a string with ANSI escape sequences.  If the input begins
// with a negative sign, the string is formatted in red else green.
func colorize(s string) string {
	if s == "" {
		return strings.Repeat(" ", 12)
	}
	color := "1;32m"
	if s[0] == '-' {
		color = "1;31m"
	}
	return fmt.Sprintf("\033[%s%12s\033[0m", color, s)
}

func bold(s string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", s)
}

// human returns a string form of the given number in base 10 with commas
// after every three orders of magnitude.
func human(s string) string {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return ""
	}
	return humanize.Commaf(v)
}

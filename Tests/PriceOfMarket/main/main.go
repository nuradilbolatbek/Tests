package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type coin struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"current_price"`
}

type crypto struct {
	Data       map[string]coin
	prevUpdate time.Time
	mutex      sync.RWMutex
}

type Client struct {
	Url        string
	Crypto     *crypto
	Updatetime time.Duration
}

func NewCryptoClient(url string, updatetime time.Duration) *Client {
	return &Client{
		Url: url,
		Crypto: &crypto{
			Data:       make(map[string]coin),
			prevUpdate: time.Now().Add(-updatetime),
		},
		Updatetime: updatetime,
	}
}

func (c *Client) FetchData() error {
	c.Crypto.mutex.Lock()
	defer c.Crypto.mutex.Unlock()
	//fmt.Printf("Time since last update: %v\n", time.Since(c.Crypto.prevUpdate))

	if time.Since(c.Crypto.prevUpdate) < c.Updatetime {
		//log.Println("currency doesn't need update")
		return nil
	}

	response, err := http.Get(c.Url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var coins []coin
	if err = json.Unmarshal(body, &coins); err != nil {
		return err
	}

	for _, coin := range coins {
		c.Crypto.Data[coin.ID] = coin
	}
	c.Crypto.prevUpdate = time.Now()
	//log.Printf("updated at %v", c.Crypto.prevUpdate)

	return nil
}

func (c *Client) GetCoinPrice(id string) (float64, bool) {
	c.Crypto.mutex.RLock()
	defer c.Crypto.mutex.RUnlock()

	coin, exists := c.Crypto.Data[id]
	return coin.Price, exists
}

func main() {
	client := NewCryptoClient("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1", 10*time.Minute)

	if err := client.FetchData(); err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("All cryptocurrency prices:")
	for id, coin := range client.Crypto.Data {
		fmt.Printf("%s (%s): $%.2f\n", coin.Name, id, coin.Price)
	}

	for {
		if err := client.FetchData(); err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Println("Enter the ID of the cryptocurrency to get its current price or enter 'all' to get all cryptocurrency (or type 'exit' to quit):")
		var coinID string
		if _, err := fmt.Scanln(&coinID); err != nil {
			log.Fatalf("Error reading input: %v", err)
		}
		if coinID == "exit" {
			break
		}
		if coinID == "all" {
			fmt.Println("All cryptocurrency prices:")
			for id, coin := range client.Crypto.Data {
				fmt.Printf("%s (%s): $%.2f\n", coin.Name, id, coin.Price)
			}
			continue
		}

		price, found := client.GetCoinPrice(coinID)
		if found {
			fmt.Printf("The current price of %s is $%.2f\n", coinID, price)
		} else {
			fmt.Println("not found, try again")
		}
	}
}

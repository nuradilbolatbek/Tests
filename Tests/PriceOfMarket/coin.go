package PriceOfMarket

import (
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

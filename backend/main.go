package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/cors"
)

type LTPResponse struct {
	Pair   string  `json:"pair"`
	Amount float64 `json:"amount"`
}

type APIResponse struct {
	LTP []LTPResponse `json:"ltp"`
}

type cacheItem struct {
	value      float64
	expiration time.Time
}

var (
	cache     = sync.Map{}
	cacheTTL  = 60 * time.Second // Cache duration: 60 seconds
	cacheLock = sync.Mutex{}
)

func fetchLTP(pair string) (float64, error) {
	// Check cache
	if cached, ok := getCachedLTP(pair); ok {
		return cached, nil
	}

	// Fetch from Kraken API
	url := fmt.Sprintf("https://api.kraken.com/0/public/Ticker?pair=%s", pair)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	errorData, ok := result["error"].([]interface{})
	if ok && len(errorData) > 0 {
		return 0, fmt.Errorf("error fetching LTP: %v", errorData)
	}

	resultData := result["result"].(map[string]interface{})
	for _, v := range resultData {
		tickerData := v.(map[string]interface{})
		lastTradeClosed := tickerData["c"].([]interface{})
		amount := lastTradeClosed[0].(string)
		parsedAmount, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			return 0, err
		}

		// Cache the result
		setCachedLTP(pair, parsedAmount)

		return parsedAmount, nil
	}

	return 0, fmt.Errorf("could not find LTP")
}

func getCachedLTP(pair string) (float64, bool) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	item, ok := cache.Load(pair)
	if !ok {
		return 0, false
	}

	cached := item.(cacheItem)
	if time.Now().After(cached.expiration) {
		// Cache expired
		cache.Delete(pair)
		return 0, false
	}

	return cached.value, true
}

func setCachedLTP(pair string, value float64) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	cache.Store(pair, cacheItem{
		value:      value,
		expiration: time.Now().Add(cacheTTL),
	})
}

func ltpHandler(w http.ResponseWriter, r *http.Request) {
	pairsParam := r.URL.Query().Get("pairs")
	var pairs []string
	var wg sync.WaitGroup

	if pairsParam == "" {
		pairs = []string{"BTCUSD", "BTCCHF", "BTCEUR"}
	} else {
		pairs = strings.Split(pairsParam, ",")
	}

	ltpData := []LTPResponse{}
	for _, pair := range pairs {
		wg.Add(1)
		var amount float64
		var err error
		pair := pair // new var per iteration
		go func() {
			defer wg.Done()
			amount, err = fetchLTP(pair)
			if err != nil {
				log.Printf("Error fetching LTP for %s: %v, %v", pair, err, pairs)
			} else {
				ltpData = append(ltpData, LTPResponse{
					Pair:   pair,
					Amount: amount,
				})
			}
		}()
	}
	wg.Wait()

	response := APIResponse{LTP: ltpData}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/api/v1/ltp", ltpHandler)

	handler := cors.Default().Handler(r)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

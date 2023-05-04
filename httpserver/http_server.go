package httpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"yuval/home-exercise/cyolo/wordcounter"
)

type WordCountable interface {
	CountWord(word string)
	TopFiveFrequentWords() map[string]int
	LowestFrequency() int
	MedianFrequency() int
}

type HttpServer struct {
	counter WordCountable
	mutex   *sync.Mutex
}

type StatsResponse struct {
	Top5            map[string]int `json:"Top5"`
	MinFrequency    int            `json:"Least"`
	MedianFrequency int            `json:"Median"`
}

func New() *HttpServer {
	return &HttpServer{
		counter: wordcounter.New(),
		mutex:   &sync.Mutex{}}
}

func (hs *HttpServer) homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Yuval's home exercise for Cyolo!")
}

func (hs *HttpServer) wordsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	text := string(body)
	words := strings.Split(text, ",")

	hs.mutex.Lock()
	for _, word := range words {
		hs.counter.CountWord(word)
	}
	hs.mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request processed successfully"))
}

func (hs *HttpServer) statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	hs.mutex.Lock()
	top5 := hs.counter.TopFiveFrequentWords()
	min := hs.counter.LowestFrequency()
	med := hs.counter.MedianFrequency()
	hs.mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StatsResponse{
		Top5:            top5,
		MinFrequency:    min,
		MedianFrequency: med,
	})
}

func (hs *HttpServer) Start() {
	http.HandleFunc("/", hs.homePage)
	http.HandleFunc("/words", hs.wordsHandler)
	http.HandleFunc("/stats", hs.statsHandler)
	http.ListenAndServe(":8081", nil)
}

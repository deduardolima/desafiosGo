package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Response struct {
	USDBRL Quote `json:"USDBRL"`
}

type Quote struct {
    Bid string `json:"bid"`
}


func main() {
	db, err := sql.Open("sqlite3", "./../db/quotes.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		getQuote(w, r, db)
	})
	http.ListenAndServe(":8080", mux)
}

func getQuote(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	start := time.Now()


	quote, err := fetchQuote(r.Context())
	if err != nil {
		log.Println("Erro ao obter a cotação:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := saveQuote(db, quote.Bid); err != nil {
		log.Println("Erro ao salvar a cotação:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quote)

	fmt.Printf("%s para processar a solicitação.\n", time.Since(start))
}

func fetchQuote(ctx context.Context) (*Quote, error) {
	apiCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(apiCtx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	var resp Response
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp.USDBRL, nil
}

func saveQuote(db *sql.DB, bidStr string) error {
	bid, err := strconv.ParseFloat(bidStr, 64)
	if err != nil {
		return err
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = db.ExecContext(dbCtx, "INSERT INTO quotes (bid) VALUES (?)", bid)
	return err
}

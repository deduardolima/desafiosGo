package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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
	defer r.Body.Close()

	apiCtx, apiCancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer apiCancel()

	req, err := http.NewRequestWithContext(apiCtx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer res.Body.Close()

	responseData, err := getPrice(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bid, err := strconv.ParseFloat(responseData.USDBRL.Bid, 64)
	if err != nil {
		http.Error(w, "Erro ao parsear bid", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData.USDBRL)

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer dbCancel()

	_, err = db.ExecContext(dbCtx, "INSERT INTO quotes (bid) VALUES (?)", bid)
	if err != nil {
		log.Printf("Falha ao inserir no Banco de dados: %v\n", err)
		return
	}

	fmt.Printf(" %s para processar a solicitação.\n", time.Since(start))
}

func getPrice(body io.ReadCloser) (*Response, error) {
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var resp Response
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

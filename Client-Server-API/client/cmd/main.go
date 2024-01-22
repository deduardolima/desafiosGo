package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil{
		panic(err)
	}

    client := &http.Client{}
    resp, err := client.Do(req)
  	if err != nil{
		panic(err)
	}
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		panic(err)
	}

    err = ioutil.WriteFile("cotacao.txt", body, 0644)
	if err != nil{
		panic(err)
	}

    fmt.Println("Cotação salva em cotacao.txt")
}

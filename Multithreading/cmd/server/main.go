package main

import (
	"fmt"
	"time"

	"github.com/desafiosgo/multithreading/cmd/external"
)

func main() {
	cep := "80060010"
	api1 := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	api2 := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	for {
		c := make(chan external.ApiResponse)
		go external.FetchAPI(api1, "BrasilAPI", c)
		go external.FetchAPI(api2, "ViaCEP", c)

		select {
		case res := <-c:
			fmt.Printf("Resposta mais rápida de %s: %s (demorou %s)\n", res.Api, res.Data, res.Duration)
		case <-time.After(1 * time.Second):
			fmt.Println("error: request timeout")
			return
		}
	}
}

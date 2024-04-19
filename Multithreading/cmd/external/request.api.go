package external

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type ApiResponse struct {
	Data     string
	Api      string
	Duration time.Duration
}

func FetchAPI(url string, apiName string, c chan ApiResponse) {
	startTime := time.Now()

	client := http.Client{
		Timeout: 1 * time.Second,
	}
	res, err := client.Get(url)
	duration := time.Since(startTime) // Calcula a duração da requisição

	if err != nil {
		c <- ApiResponse{Data: fmt.Sprintf("error: %v", err), Api: apiName, Duration: duration}
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		c <- ApiResponse{Data: fmt.Sprintf("error: %v", err), Api: apiName, Duration: duration}
		return
	}
	c <- ApiResponse{Data: string(body), Api: apiName, Duration: duration}
}

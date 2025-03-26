package config

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func PingBaseAPI() {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(APIBASEURL)
	if err != nil {
		log.Fatalf("base API is not reachable")
	}

	defer resp.Body.Close()

	if resp.StatusCode == 500 {
		log.Fatalf("base API is not responsible at the moment")
	}

	fmt.Println("✅ API is available!")
}

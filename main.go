package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	statusCode int
	err        error
}

func worker(url string, requests int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < requests; i++ {
		resp, err := http.Get(url)

		if err != nil {
			results <- Result{statusCode: 0, err: err}
			continue
		}
		results <- Result{statusCode: resp.StatusCode}
		resp.Body.Close()
	}
}

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 1, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" {
		fmt.Println("É necessário fornecer uma URL.")
		return
	}

	fmt.Printf("Iniciando teste de carga em %s com %d requisições e %d concorrentes.\n", *url, *requests, *concurrency)

	results := make(chan Result, *requests)
	var wg sync.WaitGroup

	requestsPerWorker := *requests / *concurrency
	excessRequests := *requests % *concurrency

	startTime := time.Now()
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)

		r := requestsPerWorker
		if excessRequests > 0 {
			r++
			excessRequests--
		}

		go worker(*url, r, results, &wg)
	}

	wg.Wait()
	close(results)
	totalDuration := time.Since(startTime)

	statusCount := make(map[int]int)
	errs := make(map[string]struct{})

	var totalRequests int

	for result := range results {
		statusCount[result.statusCode]++
		totalRequests++

		if result.statusCode == 0 {
			errs[result.err.Error()] = struct{}{}
		}
	}

	fmt.Println("\nRelatório do Teste:")
	fmt.Printf("Tempo total: %v\n", totalDuration)
	fmt.Printf("Total de requests: %d\n", totalRequests)
	fmt.Printf("Status 200: %d\n", statusCount[200])
	for status, count := range statusCount {
		if status != 200 {
			fmt.Printf("Status %d: %d\n", status, count)
		}
	}

	if len(errs) > 0 {
		fmt.Println("Lista de erros:")
		for k := range errs {
			fmt.Printf("    - %s\n", k)
		}
	}

}

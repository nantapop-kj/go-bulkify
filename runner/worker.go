package runner

import (
	"log"
	"net/http"
	"sync"

	"github.com/nantapop-kj/go-bulkify/config"
)

type result struct {
	workerID int
	index    int
	label    string
	err      error
}

func worker(id int, jobs <-chan int, results chan<- result, wg *sync.WaitGroup, cfg config.Config) {
	defer wg.Done()

	client := &http.Client{Timeout: cfg.RequestTimeout}

	for index := range jobs {
		payload, label := cfg.BuildPayload(index)
		err := ClientRequest(client, cfg, payload)
		results <- result{
			workerID: id,
			index:    index,
			label:    label,
			err:      err,
		}
	}
}

func Run(cfg config.Config) (successCount, failCount int) {
	jobs := make(chan int, cfg.TotalRecords)
	results := make(chan result, cfg.TotalRecords)

	var wg sync.WaitGroup

	for w := 1; w <= cfg.WorkerCount; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg, cfg)
	}

	for i := 1; i <= cfg.TotalRecords; i++ {
		jobs <- i
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		if r.err != nil {
			log.Printf("❌ [Worker %d | #%d] FAILED  — %v", r.workerID, r.index, r.err)
			failCount++
		} else {
			log.Printf("✅ [Worker %d | #%d] SUCCESS — %s", r.workerID, r.index, r.label)
			successCount++
		}
	}

	return
}

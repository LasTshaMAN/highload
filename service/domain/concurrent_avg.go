package domain

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/LasTshaMAN/Go-Execute/jobs"
)

func NewConcurrentAvg(baseURL string, httpClient *http.Client) *ConcurrentAvg {
	return &ConcurrentAvg{
		baseURL:    baseURL,
		httpClient: httpClient,
		executor:   jobs.NewExecutor(32 * 1024, 1),
	}
}

type ConcurrentAvg struct {
	baseURL    string
	httpClient *http.Client
	executor   *jobs.Executor
}

func (a *ConcurrentAvg) Value() (int, error) {
	const fastN = 5
	fastValues, err := a.sendRequests("/api/fast", fastN)
	if err != nil {
		return 0, err
	}

	const slowN = 5
	slowValues, err := a.sendRequests("/api/slow", slowN)
	if err != nil {
		return 0, err
	}

	const randomN = 20
	randomValues, err := a.sendRequests("/api/random", randomN)
	if err != nil {
		return 0, err
	}

	return (sum(fastValues) + sum(slowValues) + sum(randomValues)) / (fastN + slowN + randomN), nil
}

func (a *ConcurrentAvg) sendRequests(url string, n int) (result chan int, err error) {
	result = make(chan int, n)

	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)

		a.executor.Enqueue(func() {
			v, rErr := a.sendRequest(url)
			if rErr != nil {
				err = rErr
				return
			}
			result <- v

			wg.Done()
		})
	}
	a.executor.Enqueue(func() {
		wg.Wait()

		close(result)
	})

	return
}

func (a *ConcurrentAvg) sendRequest(url string) (int, error) {
	resp, err := a.httpClient.Get(a.baseURL + url)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("got `%d` status code for `%s`", resp.StatusCode, url)
	}
	answer, err := parse(resp.Body)
	if err != nil {
		return 0, err
	}
	return answer.Value, nil
}

func sum(values chan int) (result int) {
	for v := range values {
		result += v
	}
	return
}

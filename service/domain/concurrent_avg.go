package domain

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/LasTshaMAN/Go-Execute/jobs"
)

func NewConcurrentAvg(baseURL string, httpClient *http.Client, workersN int) *ConcurrentAvg {
	return &ConcurrentAvg{
		baseURL:    baseURL,
		httpClient: httpClient,
		executor:   jobs.NewExecutor(workersN, 1),
	}
}

type ConcurrentAvg struct {
	baseURL    string
	httpClient *http.Client
	executor   *jobs.Executor
}

func (a *ConcurrentAvg) Value() (int, error) {
	const fastN = 5
	fastValues, fastErrs := a.sendRequests("/api/fast", fastN)

	const slowN = 5
	slowValues, slowErrs := a.sendRequests("/api/slow", slowN)

	const randomN = 20
	randomValues, randomErrs := a.sendRequests("/api/random", randomN)

	if err := hasErrors(fastErrs); err != nil {
		return 0, err
	}
	if err := hasErrors(slowErrs); err != nil {
		return 0, err
	}
	if err := hasErrors(randomErrs); err != nil {
		return 0, err
	}
	return (sum(fastValues) + sum(slowValues) + sum(randomValues)) / (fastN + slowN + randomN), nil
}

func (a *ConcurrentAvg) sendRequests(url string, n int) (result chan int, errs chan error) {
	result = make(chan int, n)
	errs = make(chan error, n)

	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)

		a.executor.Enqueue(func() {
			defer wg.Done()

			v, err := a.sendRequest(url)
			if err != nil {
				errs <- err
				return
			}
			result <- v
		})
	}
	a.executor.Enqueue(func() {
		wg.Wait()

		close(result)
		close(errs)
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

func hasErrors(errs chan error) error {
	var errMsgs []string
	for err := range errs {
		if err != nil {
			errMsgs = append(errMsgs, err.Error())
		}
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "; "))
	}
	return nil
}

func sum(values chan int) (result int) {
	for v := range values {
		result += v
	}
	return
}

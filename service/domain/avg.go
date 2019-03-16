//go:generate mockgen -source avg.go -destination avg.mock.gen.go -package domain

package domain

import (
	"encoding/json"
	"fmt"
	"highload/service/api/responses"
	"io"
	"net/http"
)

type Avg interface {
	Value() int
}

func NewAvg(baseURL string, httpClient *http.Client) *AvgImpl {
	return &AvgImpl{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

type AvgImpl struct {
	baseURL    string
	httpClient *http.Client
}

func (a *AvgImpl) Value() (int, error) {
	var sum int

	const fastN = 5
	v, err := a.sendRequests("/api/fast", fastN)
	if err != nil {
		return 0, err
	}
	sum += v

	const slowN = 5
	v, err = a.sendRequests("/api/slow", slowN)
	if err != nil {
		return 0, err
	}
	sum += v

	const randomN = 20
	v, err = a.sendRequests("/api/random", randomN)
	if err != nil {
		return 0, err
	}
	sum += v

	return sum / (fastN + slowN + randomN), nil
}

func (a *AvgImpl) sendRequests(url string, n int) (int, error) {
	var total int

	for i := 0; i < n; i++ {
		v, err := a.sendRequest(url)
		if err != nil {
			return 0, err
		}
		total += v
	}

	return total, nil
}

func (a *AvgImpl) sendRequest(url string) (int, error) {
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

func parse(respBody io.Reader) (responses.Answer, error) {
	var result responses.Answer
	if err := json.NewDecoder(respBody).Decode(&result); err != nil {
		return responses.Answer{}, err
	}
	return result, nil
}

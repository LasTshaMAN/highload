package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/tsenart/vegeta/lib/plot"

	"github.com/tsenart/vegeta/lib"
)

var fastTargeter = vegeta.NewStaticTargeter(vegeta.Target{
	Method: "GET",
	URL:    "https://127.0.0.1:8002/api/fast",
})

var endpointTargeter = vegeta.NewStaticTargeter(vegeta.Target{
	Method: "GET",
	URL:    "https://127.0.0.1:8001/api/endpoint",
})

//var targeter = fastTargeter
var targeter = endpointTargeter

func main() {
	const assetsDir = "./vegeta_client/assets"
	if err := os.MkdirAll(assetsDir, os.ModePerm); err != nil {
		panic(err)
	}
	var dataFile = path.Join(assetsDir, "attack_out.bin")
	var plotFile = path.Join(assetsDir, "plot.html")

	metrics, err := runLoadTest(dataFile)
	if err != nil {
		panic(err)
	}

	if err := savePlot(dataFile, plotFile); err != nil {
		panic(err)
	}

	fmt.Printf("attack rate: %f\n", metrics.Rate)
	fmt.Printf("total requests done: %d\n", metrics.Requests)
	fmt.Printf("success rate: %f\n", metrics.Success)
	fmt.Printf("errors: %v\n", metrics.Errors)
	fmt.Println()
	fmt.Println("Latencies:")
	fmt.Printf("-- mean: %s\n", metrics.Latencies.Mean)
	fmt.Printf("-- 50th percentile: %s\n", metrics.Latencies.P50)
	fmt.Printf("-- 95th percentile: %s\n", metrics.Latencies.P95)
	fmt.Printf("-- 99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf("-- max: %s\n", metrics.Latencies.Max)
}

func runLoadTest(dataFile string) (result vegeta.Metrics, err error) {
	out, err := os.Create(dataFile)
	if err != nil {
		return
	}
	defer func() {
		if cErr := out.Close(); cErr != nil {
			err = fmt.Errorf(strings.Join([]string{err.Error(), cErr.Error()}, "; "))
		}
	}()
	enc := vegeta.NewEncoder(out)

	attacker := vegeta.NewAttacker(vegeta.HTTP2(true))
	results := attacker.Attack(
		targeter,
		vegeta.Rate{
			Freq: 500,
			Per:  time.Second,
		},
		120*time.Second,
		"",
	)
	for res := range results {
		result.Add(res)
		if eErr := enc.Encode(res); eErr != nil {
			err = fmt.Errorf("error encoding result: %s", eErr)
			return
		}
	}
	result.Close()

	return
}

func savePlot(dataFile string, plotFile string) error {
	in, err := os.Open(dataFile)
	if err != nil {
		return err
	}
	dec := vegeta.DecoderFor(in)
	if dec == nil {
		return fmt.Errorf("decode: can't detect encoding of %s", dataFile)
	}

	p := plot.New(
		plot.Title(dataFile),
		plot.Label(plot.ErrorLabeler),
	)
	for {
		var r vegeta.Result
		if err := dec.Decode(&r); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if err = p.Add(&r); err != nil {
			return err
		}
	}
	p.Close()

	out, err := os.Create(plotFile)
	if err != nil {
		return err
	}
	defer func() {
		if cErr := out.Close(); cErr != nil {
			err = fmt.Errorf(strings.Join([]string{err.Error(), cErr.Error()}, "; "))
		}
	}()
	if _, err := p.WriteTo(out); err != nil {
		return err
	}

	return nil
}

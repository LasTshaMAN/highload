# highload

This project explores how micro-services written in Go behave under load.

## Architecture

[service](https://github.com/LasTshaMAN/highload/tree/master/service) folder contains the target service of interest.

[mocked_service](https://github.com/LasTshaMAN/highload/tree/master/mocked_service) is a dependent service target service relies on.

## Tools

[Prometheus](https://github.com/prometheus/prometheus) and [Grafana](https://github.com/grafana/grafana) are used for analysing service behavior.
  
[Vegeta](https://github.com/tsenart/vegeta) is used for loading service with requests.

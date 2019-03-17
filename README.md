# highload

This project explores how micro-services written in Go behave under load.

## Architecture

[service](https://github.com/LasTshaMAN/highload/tree/master/service) folder contains the target service of interest.

[mocked_service](https://github.com/LasTshaMAN/highload/tree/master/mocked_service) is a dependent service target service relies on.

## Tools

[Prometheus](https://github.com/prometheus/prometheus) is used to gather vital metrics ([Grafana](https://github.com/grafana/grafana) could be used to visualize these metrics in eye-pleasing way).
  
[Vegeta](https://github.com/tsenart/vegeta) is used for loading service with requests.

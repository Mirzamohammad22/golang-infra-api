# golang-infra-api

## Table of Content

- [Prerequisites](#Prerequisites)
- [Usage](#usage)
- [Assumptions made for application](#assumptions)
- [Design Question](#design-question)

## Prerequisites

- [go](https://go.dev/doc/install)

## Usage

```
cd <GOPATH>/src/github.com/
git clone git@github.com:Mirzamohammad22/golang-infra-api.git
cd golang-infra-api
go get -d ./...
go run main.go
```

## Assumptions

For the simplicity of the task and to avoid any misconceptions, the following assumptions have been made during development -

- The alert configs did not have unique ids, hence alertName was assumed to be unique.
- To keep the generated alert ids the same as the provided input alert ids (len=24), a hash of the alertName is used. Ideally, UUID could have been used for unique ids.
- The summary of the actions taken on the alerts was meant only for logging, not as part of the structure containing the alert modifications.

## Design Question

Question:
Provided that the external API may be unreliable and the alert configurations can be changed by sneaky SREs manually, how would you design such an application to ensure there's no accumulated drift between the real alerts and the desired alerts configuration?

Answer:
a) Cron Job: We can set up the application to run as a scheduled cron-job at an agreed upon frequency (eg: every 5 minutes) to request the external API for the alerts and compare them with the desired configuration. If there are any deviations, the application will send a POST request to the API to update the alerts to the desired configuration. This cron job could be implemented in multiple ways such as Github workflow, Kubernetes cron job or AWS Lambda to name a few.

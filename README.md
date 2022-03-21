# golang-infra-api

## Table of Content
- [Usage](#usage)
- [Assumptions made for application](#assumptions)
- [Design Question](#design-question)

### USAGE
```
cd <GOPATH>/src/github.com/
git clone git@github.com:Mirzamohammad22/golang-infra-api.git
cd Mirzamohammad22/golang-infra-api
go get -d ./...
go run main.go
```

### ASSUMPTIONS
For the simplicity of the task and to avoid any misconceptions around the task. Assumptions made during the solving the task has been documented

* Config had no unique id, hence alertName was assumed to be unique
* AlertID should match the length of input alert id (len=24), hence a hash of the alertName is used. Else ideally would use UUID for this generation
* Summary of the alerts had to be provided seperately instead of consolidating it into a bigger response json.


### DESIGN QUESTION

Question:
Providing that the external API may be unreliable and the alert configurations can be changed by sneaky SREs manually, how would you design such an application to ensure there’s no accumulated drift between the real alerts and the desired alerts configuration?

Answer:
    a) Cron Jobs: We can run the application as a scheduled cron-job(eg: every 30mins) to request the External Api for its alerts and compare it with the configuration.If there are any deviations, send a POST request to the API for correction. This cron-job could be via numerous ways eg( github workflows, kubernetes cronjobs, aws Lambda)

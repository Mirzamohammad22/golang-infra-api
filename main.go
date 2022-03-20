package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
	"gopkg.in/yaml.v2"
)

type Notification struct {
	NotificationType    string `json:"notificationType" yaml:"notificationType"`
	NotificationChannel string `json:"notificationChannel" yaml:"notificationChannel" yaml:"notificationChannel" yaml:"notificationChannel"`
	DelayMin            int    `json:"delayMin" yaml:"delayMin"`
	IntervalMin         int    `json:"intervalMin" yaml:"intervalMin"`
}

type MetricThreshold struct {
	MetricName string `json:"metricName" yaml:"metricName"`
	Operator   string `json:"operator" yaml:"operator"`
	Threshold  int    `json:"threshold" yaml:"threshold"`
	Units      string `json:"units" yaml:"units"`
}

type Alert struct {
	AlertName       string          `json:"alertName" yaml:"alertName"`
	Enabled         bool            `json:"enabled" yaml:"enabled"`
	MetricThreshold MetricThreshold `json:"metricThreshold" yaml:"metricThreshold"`
	Notifications   []Notification  `json:"notifications" yaml:"notifications"`
}

type CurrentAlert struct {
	Id      string `json:"id" yaml:"id"`
	Created time.Time `json:"created" yaml:"created"`
	Updated time.Time `json:"updated" yaml:"updated"`
	Alert
}

type Alerts struct {
	Alerts []Alert `yaml:"alerts"`
}

type CurrentAlerts struct {
	Alerts []CurrentAlert `json:"results"`
}

func getCurrentAlerts() map[string]CurrentAlert {
	var alertsArray CurrentAlerts
	jsonFile, err := ioutil.ReadFile("./input.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(jsonFile, &alertsArray)
	if err != nil {
		panic(err)
	}
	alerts := make(map[string]CurrentAlert)

	for _, alert := range alertsArray.Alerts {
		alerts[alert.AlertName] = alert
	}

	return alerts
}

func getAlertConfig() map[string]Alert {
	var alertsArray Alerts
	yamlFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &alertsArray)
	if err != nil {
		panic(err)
	}
	alerts := make(map[string]Alert)

	for _, alert := range alertsArray.Alerts {
		alerts[alert.AlertName] = alert
	}
	return alerts
}
func compareCurrentAndConfig(current map[string]CurrentAlert, config map[string]Alert){
	var alertsToBeAdded []Alert
	var alertsToBeUpdated []Alert
	for key, configAlert := range config {
		if currentAlert, found := current[key]; found {
			fmt.Printf("%+v",currentAlert.Alert)
		}
	}

}
func main() {
	currentAlerts := getCurrentAlerts()
	configuredAlerts := getAlertConfig()
	compareCurrentAndConfig(currentAlerts,configuredAlerts)

}

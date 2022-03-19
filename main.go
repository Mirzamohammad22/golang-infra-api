package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Created string `json:"created" yaml:"created"`
	Updated string `json:"updated" yaml:"updated"`
	Alert
}

type Alerts struct {
	Alerts []Alert `yaml:"alerts"`
}

type CurrentAlerts struct {
	Alerts []CurrentAlert `json:"results"`
}

func getCurrentAlerts() CurrentAlerts {
	var alerts CurrentAlerts
	jsonFile, err := ioutil.ReadFile("./input.json")
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(jsonFile, &alerts)
	return alerts
}

func getAlertConfig() Alerts {
	var alerts Alerts
	yamlFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		fmt.Println(err)
	}
	yaml.Unmarshal(yamlFile, &alerts)
	return alerts

}

func main() {
	currentAlerts := getCurrentAlerts()
	fmt.Printf("%+v", currentAlerts.Alerts[0])
	configuredAlerts := getAlertConfig()
	fmt.Printf("%+v", configuredAlerts.Alerts[0])
}

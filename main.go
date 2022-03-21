package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"time"

	"crypto/sha1"

	log "github.com/Sirupsen/logrus"
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

type AlertConfig struct {
	AlertName       string          `json:"alertName" yaml:"alertName"`
	Enabled         bool            `json:"enabled" yaml:"enabled"`
	MetricThreshold MetricThreshold `json:"metricThreshold" yaml:"metricThreshold"`
	Notifications   []Notification  `json:"notifications" yaml:"notifications"`
}

type Alert struct {
	Id      string    `json:"id" yaml:"id"`
	Created time.Time `json:"Created" yaml:"Created"`
	Updated time.Time `json:"updated" yaml:"updated"`
	AlertConfig
}

type AlertConfigs struct {
	Alerts []AlertConfig `yaml:"alerts"`
}

type Alerts struct {
	Alerts []Alert `json:"results"`
}

type SetAlertApi struct {
	AlertID string      `json:"alertID"`
	Action  string      `json:"action"`
	Body    AlertConfig `json:"body"`
}

func PrettyStructJSON(data interface{}) string {

	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatalf("Error: Failed to Marshal data. %v", err)
	}

	return string(val)
}

func getAlerts() map[string]Alert {
	var alertsArray Alerts

	jsonFile, err := ioutil.ReadFile("./input.json")
	if err != nil {
		log.Fatalf("Error: Failed to read input file. %v", err)
	}
	err = json.Unmarshal(jsonFile, &alertsArray)
	if err != nil {
		log.Fatalf("Error: Failed to unmarshal input json. %v", err)
	}
	alerts := make(map[string]Alert)

	for _, alert := range alertsArray.Alerts {
		alerts[alert.AlertName] = alert
	}
	return alerts
}

func getConfig() map[string]AlertConfig {
	var alertsArray AlertConfigs
	yamlFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("Error: Failed to read config file. %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &alertsArray)
	if err != nil {
		log.Fatalf("error:  Failed to unmarshal config yaml. %v", err)
	}
	alerts := make(map[string]AlertConfig)

	for _, alert := range alertsArray.Alerts {
		alerts[alert.AlertName] = alert
	}
	return alerts
}

func updateAlert(alertToUpdateId string, config AlertConfig) SetAlertApi {
	response := SetAlertApi{AlertID: alertToUpdateId, Action: "update", Body: config}
	return response
}

func createAlert(alertToCreate AlertConfig) SetAlertApi {
	h := sha1.New()
	h.Write([]byte(alertToCreate.AlertName))
	idFull := h.Sum(nil)
	id := string(idFull[:24])
	response := SetAlertApi{AlertID: id, Action: "create", Body: alertToCreate}
	return response
}
func deleteAlert(alertToDelete Alert) SetAlertApi {
	response := SetAlertApi{AlertID: alertToDelete.Id, Action: "delete", Body: alertToDelete.AlertConfig}
	return response
}

func CreateApiResponse(alertsToCreate []AlertConfig, alertsToUpdate []Alert, alertsToDelete []Alert, config map[string]AlertConfig) []SetAlertApi {
	var apiResponse []SetAlertApi

	for _, alert := range alertsToCreate {
		apiResponse = append(apiResponse, createAlert(alert))
	}

	for _, alert := range alertsToUpdate {
		configAlert := config[alert.AlertName]
		apiResponse = append(apiResponse, updateAlert(alert.Id, configAlert))
	}

	for _, alert := range alertsToDelete {
		apiResponse = append(apiResponse, deleteAlert(alert))
	}
	return apiResponse
}

func compareAlertsWithConfig(alerts map[string]Alert, alertConfigs map[string]AlertConfig) ([]AlertConfig, []Alert, []Alert) {
	var alertsToBeCreated []AlertConfig
	var alertsToBeUpdated []Alert
	var alertsToBeDeleted []Alert
	for key, alertConfig := range alertConfigs {
		if alert, found := alerts[key]; found {
			if !reflect.DeepEqual(alert.AlertConfig, alertConfig) {
				alertsToBeUpdated = append(alertsToBeUpdated, alert)
			}
			delete(alerts, key)
		} else {
			alertsToBeCreated = append(alertsToBeCreated, alertConfig)
		}
	}
	for _, alert := range alerts {
		alertsToBeDeleted = append(alertsToBeDeleted, alert)
	}
	return alertsToBeCreated, alertsToBeUpdated, alertsToBeDeleted
}
func setAlerts(alertsToCreate []AlertConfig, alertsToUpdate []Alert, alertsToDelete []Alert, alertConfigs map[string]AlertConfig) {
	apiResponse := CreateApiResponse(alertsToCreate, alertsToUpdate, alertsToDelete, alertConfigs)
	summary := fmt.Sprintf("SUMMARY:\n CREATED:%d \n UPDATED:%d \n DELETED:%d \n", len(alertsToCreate), len(alertsToUpdate), len(alertsToDelete))
	apiResponseJSON := PrettyStructJSON(apiResponse)
	log.Info("ALERT CONFIG STRUCTURE:", apiResponseJSON)
	log.Info(summary)
}

func main() {
	alerts := getAlerts()
	alertConfigs := getConfig()
	alertsToBeCreated, alertsToBeUpdated, alertsToBeDeleted := compareAlertsWithConfig(alerts, alertConfigs)
	setAlerts(alertsToBeCreated, alertsToBeUpdated, alertsToBeDeleted, alertConfigs)
}

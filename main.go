package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"time"

	"crypto/sha1"

	"encoding/base64"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Notification struct {
	NotificationType     string `json:"notificationType" yaml:"notificationType"`
	NotificationChannel  string `json:"notificationChannel" yaml:"notificationChannel"`
	NotificationSchedule string `json:"notificationSchedule" yaml:"notificationSchedule"`

	DelayMin    int `json:"delayMin" yaml:"delayMin"`
	IntervalMin int `json:"intervalMin" yaml:"intervalMin"`
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
	Created time.Time `json:"created" yaml:"created"`
	Updated time.Time `json:"updated" yaml:"updated"`
	AlertConfig
}

type AlertConfigs struct {
	Alerts []AlertConfig `yaml:"alerts"`
}

type Alerts struct {
	Alerts []Alert `json:"results"`
}

type AlertAction struct {
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

func updateAlert(alertID string, config AlertConfig) AlertAction {
	result := AlertAction{AlertID: alertID, Action: "update", Body: config}
	return result
}

func createAlert(alertToCreate AlertConfig) AlertAction {
	h := sha1.New()
	h.Write([]byte(alertToCreate.AlertName))
	idFull := base64.URLEncoding.EncodeToString(h.Sum(nil))
	id := string(idFull[:24])
	result := AlertAction{AlertID: id, Action: "create", Body: alertToCreate}
	return result
}
func deleteAlert(alertToDelete Alert) AlertAction {
	result := AlertAction{AlertID: alertToDelete.Id, Action: "delete", Body: alertToDelete.AlertConfig}
	return result
}

func CreateApiStruct(alertsToCreate []AlertConfig, alertsToUpdate []Alert, alertsToDelete []Alert, config map[string]AlertConfig) []AlertAction {
	alertActions := make([]AlertAction, 0)
	for _, alert := range alertsToCreate {
		alertActions = append(alertActions, createAlert(alert))
	}

	for _, alert := range alertsToUpdate {
		configAlert := config[alert.AlertName]
		alertActions = append(alertActions, updateAlert(alert.Id, configAlert))
	}

	for _, alert := range alertsToDelete {
		alertActions = append(alertActions, deleteAlert(alert))
	}
	return alertActions
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
	alertActions := CreateApiStruct(alertsToCreate, alertsToUpdate, alertsToDelete, alertConfigs)
	summary := fmt.Sprintf("Summary:\n Created: %d \n Updated: %d \n Deleted: %d \n", len(alertsToCreate), len(alertsToUpdate), len(alertsToDelete))
	alertActionsJSON := PrettyStructJSON(alertActions)
	log.Info("API Alert Config Structure:", alertActionsJSON)
	log.Info(summary)
}

func main() {
	alerts := getAlerts()
	alertConfigs := getConfig()
	alertsToBeCreated, alertsToBeUpdated, alertsToBeDeleted := compareAlertsWithConfig(alerts, alertConfigs)
	setAlerts(alertsToBeCreated, alertsToBeUpdated, alertsToBeDeleted, alertConfigs)
}

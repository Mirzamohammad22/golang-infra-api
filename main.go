package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/google/uuid"
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
	Id      string    `json:"id" yaml:"id"`
	Created time.Time `json:"Created" yaml:"Created"`
	Updated time.Time `json:"updated" yaml:"updated"`
	Alert
}

type Alerts struct {
	Alerts []Alert `yaml:"alerts"`
}

type CurrentAlerts struct {
	Alerts []CurrentAlert `json:"results"`
}

type ApiResponse struct {
	AlertID string `json:"alertID"`
	Action  string `json:"action"`
	Body    Alert  `json:"body"`
}

func getCurrentAlerts() map[string]CurrentAlert {
	var alertsArray CurrentAlerts
	jsonFile, err := ioutil.ReadFile("./input.json")
	if err != nil {
		// fmt.Printf("\n %v", err)
		// panic(err)
		log.Fatalf("error: %v", err)
	}
	err = json.Unmarshal(jsonFile, &alertsArray)
	if err != nil {
		// fmt.Printf("\n %v", err)
		log.Fatalf("error: %v", err)
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
		fmt.Printf("\n %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &alertsArray)
	if err != nil {
		fmt.Printf("\n %v", err)
	}
	alerts := make(map[string]Alert)

	for _, alert := range alertsArray.Alerts {
		alerts[alert.AlertName] = alert
	}
	return alerts
}

func PrettyStructJSON(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func updateAlert(alertToUpdateId string, config Alert) ApiResponse {
	response := ApiResponse{AlertID: alertToUpdateId, Action: "Update", Body: config}
	return response
}

func createAlert(alertToCreate Alert) ApiResponse {
	id := strings.Replace(uuid.New().String(), "-", "", -1)
	response := ApiResponse{AlertID: id, Action: "Create", Body: alertToCreate}
	return response
}
func deleteAlert(alertToDelete CurrentAlert) ApiResponse {
	response := ApiResponse{AlertID: alertToDelete.Id, Action: "Delete", Body: alertToDelete.Alert}
	return response
}

func CreateApiResponse(alertsToCreate []Alert, alertsToUpdate []CurrentAlert, alertsToDelete []CurrentAlert, config map[string]Alert) (string, string) {
	var apiResponse []ApiResponse

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
	summary := fmt.Sprintf("SUMMARY:\n CREATED:%d \n UPDATED:%d \n DELETED: %d \n", len(alertsToCreate), len(alertsToUpdate), len(alertsToDelete))
	apiResponseJSON, err := PrettyStructJSON(apiResponse)
	if err != nil {
		fmt.Println(err)
	}

	return apiResponseJSON, summary
}

func compareCurrentAndConfig(current map[string]CurrentAlert, config map[string]Alert) (string, string) {
	var alertsToBeCreated []Alert
	var alertsToBeUpdated []CurrentAlert
	var alertsToBeDeleted []CurrentAlert
	for key, configAlert := range config {
		if currentAlert, found := current[key]; found {
			if !reflect.DeepEqual(currentAlert.Alert, configAlert) {
				alertsToBeUpdated = append(alertsToBeUpdated, currentAlert)
			}
			delete(current, key)
		} else {
			alertsToBeCreated = append(alertsToBeCreated, configAlert)
		}
	}
	for _, currentAlert := range current {
		alertsToBeDeleted = append(alertsToBeDeleted, currentAlert)
	}
	result, summary := CreateApiResponse(alertsToBeCreated, alertsToBeUpdated, alertsToBeDeleted, config)
	return result, summary
}

func main() {
	currentAlerts := getCurrentAlerts()
	configuredAlerts := getAlertConfig()
	result, summary := compareCurrentAndConfig(currentAlerts, configuredAlerts)
	log.Info(result, summary)

}

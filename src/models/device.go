package models

import (
	"gopkg.in/mgo.v2/bson"
)

type (
	// User represents the structure of our resource
	Device struct {
		Id         bson.ObjectId `json:"id" bson:"_id"`
		AccountId  bson.ObjectId `json:"account_id" bson:"account_id"`
		DeviceId   string        `json:"device_id" bson:"device_id"`
		Serial     string        `json:"serial_number" bson:"serial"`
		Model      string      	 `json:"device_model" bson:"model"`
		Battery    int	         `json:"battery" bson:"battery"`
		Mode 	   int           `json:"mode" bson:"mode"`
		ActivationDate string 	 `json:"activation_date" bson:"activation_date"`
		LastCommDate   string 	 `json:"last_comm_date" bson:"last_comm_date"`
		LastConfigDate string 	 `json:"last_config_date" bson:"last_config_date"`
	}


	SensorData struct {
		Id             string `json:"id"`
		SensorType    string `json:"sensor_type"`
	}

	DeviceSensorsData struct {
		Id             string `json:"id"`
		DeviceModel    string `json:"device_model"`
		SerialNumber   string `json:"serial_number"`
		Battery        int    `json:"battery"`
		LastConfigDate string `json:"last_config_date"`
		Mode           int    `json:"mode"`
		ActivationDate string `json:"activation_date"`
		LastCommDate   string `json:"last_comm_date"`
		//Sensors    Sensors    `json:"sensors"`
		SensorData    []SensorData `json:"sensors"`

	}

	DevicesSensors struct {
		Results int          `json:"results"`
		Data    DeviceSensorsData `json:"data"`
	}


	DeviceData struct {
		Id             string `json:"id"`
		DeviceModel    string `json:"device_model"`
		SerialNumber   string `json:"serial_number"`
		Battery        int    `json:"battery"`
		LastConfigDate string `json:"last_config_date"`
		Mode           int    `json:"mode"`
		ActivationDate string `json:"activation_date"`
		LastCommDate   string `json:"last_comm_date"`
	}

	Devices struct {
		Results int          `json:"results"`
		Data    []DeviceData `json:"data"`
	}
)

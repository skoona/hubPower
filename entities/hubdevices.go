package entities

import "fyne.io/fyne/v2/data/binding"

type DeviceList struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Label string `json:"label"`
	Type  string `json:"type"`
	Room  string `json:"room"`
}

type Device struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Label      string `json:"label"`
	Type       string `json:"type"`
	Room       string `json:"room"`
	Attributes []struct {
		Name         string      `json:"name"`
		CurrentValue interface{} `json:"currentValue"`
		DataType     string      `json:"dataType"`
		Values       []string    `json:"values,omitempty"`
	} `json:"attributes"`
	Capabilities []interface{} `json:"capabilities"`
	Commands     []string      `json:"commands"`
}

type DeviceDetails struct {
	Name         string      `json:"name"`
	Label        string      `json:"label"`
	Type         string      `json:"type"`
	Id           string      `json:"id"`
	Date         *string     `json:"date"`
	Model        interface{} `json:"model"`
	Manufacturer interface{} `json:"manufacturer"`
	Room         string      `json:"room"`
	Capabilities []string    `json:"capabilities"`
	Attributes   struct {
		Voltage   string      `json:"voltage"`
		DataType  string      `json:"dataType"`
		Values    interface{} `json:"values"`
		Energy    string      `json:"energy"`
		Amperage  string      `json:"amperage"`
		Frequency interface{} `json:"frequency"`
		Switch    string      `json:"switch"`
		Power     string      `json:"power"`
	} `json:"attributes"`
	Commands []struct {
		Command string `json:"command"`
	} `json:"commands"`
	BWattValue binding.Float `json:"-"`
}

type DeviceEvent struct {
	DeviceId      string `json:"device_id"`
	Label         string `json:"label"`
	Room          string `json:"room"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	Date          string `json:"date"`
	Unit          string `json:"unit"`
	IsStateChange bool   `json:"isStateChange"`
	Source        string `json:"source"`
}

type DeviceEventStream struct {
	Content struct {
		Name            string      `json:"name"`
		Value           string      `json:"value"`
		DisplayName     string      `json:"displayName"`
		DeviceId        string      `json:"deviceId"`
		DescriptionText string      `json:"descriptionText"`
		Unit            string      `json:"unit"`
		Type            interface{} `json:"type"`
		Data            interface{} `json:"data"`
	} `json:"content"`
}

type DeviceCapabilities struct {
	Capabilities []interface{} `json:"capabilities"`
}

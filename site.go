package solaredge

import (
	"fmt"
	"net/url"
	"time"
)

type Site struct {
	Id        int     `json:"id"`
	Name      string  `json:"name,omitempty"`
	AccountId int     `json:"accountId,omitempty"`
	Status    string  `json:"status,omitempty"`
	PeakPower float64 `json:"peakPower,omitempty"`
	Location  struct {
		Country     string `json:"country,omitempty"`
		City        string `json:"city,omitempty"`
		Address     string `json:"address,omitempty"`
		Zip         string `json:"zip,omitempty"`
		TimeZone    string `json:"timeZone,omitempty"`
		CountryCode string `json:"countryCode,omitempty"`
	} `json:"location,omitempty"`
}

type Inventory struct {
	Meters    []Meter    `json:"meters,omitempty"`
	Sensors   []Sensor   `json:"sensors,omitempty"`
	Gateways  []Gateway  `json:"gateways,omitempty"`
	Batteries []Battery  `json:"batteries,omitempty"`
	Inverters []Inverter `json:"inverters,omitempty"`
}

type storageData struct {
	Batteries []StorageBattery `json:"batteries,omitempty"`
}

func (sc *SiteClient) Details() (*Site, error) {
	var res Site
	details := struct {
		Details *Site `json:"details"`
	}{
		Details: &res,
	}
	return &res, sc.Get(fmt.Sprintf("/site/%s/details.json", sc.siteid), nil, &details)
}

func (sc *SiteClient) Inventory() (*Inventory, error) {
	var res Inventory
	details := struct {
		Inventory *Inventory `json:"Inventory"`
	}{
		Inventory: &res,
	}
	return &res, sc.Get(fmt.Sprintf("/site/%s/inventory.json", sc.siteid), nil, &details)
}

func (sc *SiteClient) StorageData(start, end time.Time) ([]StorageBattery, error) {
	var res storageData
	details := struct {
		Data *storageData `json:"storageData"`
	}{
		Data: &res,
	}
	parms := url.Values{
		"startTime": []string{start.Format(datetimePattern)},
		"endTime":   []string{end.Format(datetimePattern)},
	}
	return res.Batteries, sc.Get(fmt.Sprintf("/site/%s/storageData.json", sc.siteid), parms, &details)
}

func (sc *SiteClient) PowerDetails(start, end time.Time) (*PowerDetails, error) {
	var res PowerDetails
	details := struct {
		Data *PowerDetails `json:"powerDetails"`
	}{
		Data: &res,
	}
	parms := url.Values{
		"startTime": []string{start.Format(datetimePattern)},
		"endTime":   []string{end.Format(datetimePattern)},
	}
	return &res, sc.Get(fmt.Sprintf("/site/%s/powerDetails.json", sc.siteid), parms, &details)
}

func (sc *SiteClient) EnergyDetails(tu TimeUnit, start, end time.Time) (*EngergyDetails, error) {
	var res EngergyDetails
	details := struct {
		Data *EngergyDetails `json:"energyDetails"`
	}{
		Data: &res,
	}
	parms := url.Values{
		"startTime": []string{start.Format(datetimePattern)},
		"endTime":   []string{end.Format(datetimePattern)},
		"timeUnit":  []string{string(tu)},
	}
	return &res, sc.Get(fmt.Sprintf("/site/%s/energyDetails.json", sc.siteid), parms, &details)
}

func (sc *SiteClient) PowerFlow() (*PowerFlow, error) {
	var res PowerFlow
	details := struct {
		Flow *PowerFlow `json:"siteCurrentPowerFlow"`
	}{
		Flow: &res,
	}
	return &res, sc.Get(fmt.Sprintf("/site/%s/currentPowerFlow.json", sc.siteid), nil, &details)
}

func (sc *SiteClient) Overview() (*OverviewData, error) {
	var res OverviewData
	details := struct {
		Data *OverviewData `json:"overview"`
	}{
		Data: &res,
	}
	return &res, sc.Get(fmt.Sprintf("/site/%s/overview.json", sc.siteid), nil, &details)
}

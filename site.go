package solaredge

import (
	"fmt"
	"net/url"
	"time"
)

type storageData struct {
	Batteries []StorageBattery `json:"batteries,omitempty"`
}

// Details returns site information.
func (sc *SiteClient) Details() (*Site, error) {
	var res Site
	details := struct {
		Details *Site `json:"details"`
	}{
		Details: &res,
	}
	return &res, sc.get(fmt.Sprintf("/site/%s/details.json", sc.siteid), nil, &details)
}

// Inventory returns the inventory of a site.
func (sc *SiteClient) Inventory() (*Inventory, error) {
	var res Inventory
	details := struct {
		Inventory *Inventory `json:"Inventory"`
	}{
		Inventory: &res,
	}
	return &res, sc.get(fmt.Sprintf("/site/%s/inventory.json", sc.siteid), nil, &details)
}

// StorageData returns a list of battery elements.
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
	return res.Batteries, sc.get(fmt.Sprintf("/site/%s/storageData.json", sc.siteid), parms, &details)
}

// PowerDetails returns the power details.
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
	return &res, sc.get(fmt.Sprintf("/site/%s/powerDetails.json", sc.siteid), parms, &details)
}

// EnergyDetails returns the energy details.
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
	return &res, sc.get(fmt.Sprintf("/site/%s/energyDetails.json", sc.siteid), parms, &details)
}

// PowerFlow returns the current powerflow.
func (sc *SiteClient) PowerFlow() (*PowerFlow, error) {
	var res PowerFlow
	details := struct {
		Flow *PowerFlow `json:"siteCurrentPowerFlow"`
	}{
		Flow: &res,
	}
	return &res, sc.get(fmt.Sprintf("/site/%s/currentPowerFlow.json", sc.siteid), nil, &details)
}

// Overview returns the current overview of the site.
func (sc *SiteClient) Overview() (*OverviewData, error) {
	var res OverviewData
	details := struct {
		Data *OverviewData `json:"overview"`
	}{
		Data: &res,
	}
	return &res, sc.get(fmt.Sprintf("/site/%s/overview.json", sc.siteid), nil, &details)
}

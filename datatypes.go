package solaredge

import (
	"fmt"
	"strings"
	"time"
)

const (
	datetimePattern = "2006-01-02 15:04:05"
)

var (
	SiteZone string
)

type TimeUnit string

var (
	Quarter_Of_An_Hour TimeUnit = "QUARTER_OF_AN_HOUR"
	Hour               TimeUnit = "HOUR"
	Day                TimeUnit = "DAY"
	Week               TimeUnit = "WEEK"
	Month              TimeUnit = "MONTH"
	Year               TimeUnit = "YEAR"
)

func init() {
	t := time.Now()
	SiteZone, _ = t.Zone()
}

type SETime time.Time

func (f *SETime) MarshalJSON() ([]byte, error) {
	t := time.Time(*f)
	return []byte(fmt.Sprintf("%q", t.Format(datetimePattern))), nil
}

func (f *SETime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	loc, _ := time.LoadLocation(SiteZone)
	t, err := time.ParseInLocation(datetimePattern, s, loc)
	if err != nil {
		return err
	}
	*f = SETime(t)
	return nil
}

type Meter struct {
	Name                       string `json:"name,omitempty"`
	Manufacturer               string `json:"manufacturer,omitempty"`
	Model                      string `json:"model,omitempty"`
	FirmwareVersion            string `json:"firmwareVersion,omitempty"`
	ConnectedTo                string `json:"connectedTo,omitempty"`
	ConnectedSolaredgeDeviceSN string `json:"connectedSolaredgeDeviceSN,omitempty"`
	Type                       string `json:"type,omitempty"`
	Form                       string `json:"form,omitempty"`
	SN                         string `json:"SN,omitempty"`
}

type Sensor struct {
	Category                   string `json:"category,omitempty"`
	Type                       string `json:"type,omitempty"`
	ConnectedTo                string `json:"connectedTo,omitempty"`
	ConnectedSolaredgeDeviceSN string `json:"connectedSolaredgeDeviceSN,omitempty"`
}

type Gateway struct {
	Name            string `json:"name,omitempty"`
	SerialNumber    string `json:"serialNumber,omitempty"`
	FirmwareVersion string `json:"firmwareVersion,omitempty"`
}

type Battery struct {
	Name                       string  `json:"name,omitempty"`
	Manufacturer               string  `json:"manufacturer,omitempty"`
	Model                      string  `json:"model,omitempty"`
	FirmwareVersion            string  `json:"firmwareVersion,omitempty"`
	ConnectedTo                string  `json:"connectedTo,omitempty"`
	ConnectedSolaredgeDeviceSN string  `json:"connectedSolaredgeDeviceSN,omitempty"`
	ConnectedInverterSn        string  `json:"connectedInverterSn,omitempty"`
	NameplateCapacity          float64 `json:"nameplateCapacity,omitempty"`
	SN                         string  `json:"SN,omitempty"`
}

type Inverter struct {
	Name                string `json:"name,omitempty"`
	Manufacturer        string `json:"manufacturer,omitempty"`
	Model               string `json:"model,omitempty"`
	CommunicationMethod string `json:"communicationMethod,omitempty"`
	CPUVersion          string `json:"cpuVersion,omitempty"`
	SN                  string `json:"SN,omitempty"`
	ConnectedOptimizers int    `json:"connectedOptimizers,omitempty"`
}

type StorageBatteryTelemetry struct {
	Timestamp                SETime  `json:"timeStamp,omitempty"`
	Power                    float64 `json:"power,omitempty"`
	State                    int     `json:"batteryState,omitempty"`
	LifetimeEnergyDischarged int64   `json:"lifeTimeEnergyDischarged,omitempty"`
	LifetimeEnergyCharged    int64   `json:"lifeTimeEnergyCharged,omitempty"`
	PercentageState          float64 `json:"batteryPercentageState,omitempty"`
	FullPackEngergyAvailable float64 `json:"fullPackEnergyAvailable,omitempty"`
	InternalTemp             float64 `json:"internalTemp,omitempty"`
	ACGridCharging           float64 `json:"ACGridCharging,omitempty"`
}
type StorageBattery struct {
	Nameplate   float64                   `json:"nameplate,omitempty"`
	SN          string                    `json:"serialNumber,omitempty"`
	ModelNumber string                    `json:"modelNumber,omitempty"`
	Telemetries []StorageBatteryTelemetry `json:"telemetries,omitempty"`
}

type MeterValue struct {
	Date  SETime  `json:"date"`
	Value float64 `json:"value"`
}

type MeteredValue struct {
	Type   string       `json:"type"`
	Values []MeterValue `json:"values"`
}

type PowerDetails struct {
	TimeUnit TimeUnit       `json:"timeUnit"`
	Unit     string         `json:"unit"`
	Meters   []MeteredValue `json:"meters"`
}

type EngergyDetails struct {
	TimeUnit TimeUnit       `json:"timeUnit"`
	Unit     string         `json:"unit"`
	Meters   []MeteredValue `json:"meters"`
}

type PowerFlowConnection struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type PowerFlowStatus struct {
	Status       string  `json:"status"`
	CurrentPower float64 `json:"currentPower"`
}

type StoragePowerFlowStatus struct {
	PowerFlowStatus
	ChargeLevel int  `json:"chargeLevel"`
	Critical    bool `json:"critical"`
}
type PowerFlow struct {
	Unit        string                  `json:"unit"`
	Connections []PowerFlowConnection   `json:"connections,omitempty"`
	Grid        PowerFlowStatus         `json:"GRID,omitempty"`
	Load        PowerFlowStatus         `json:"LOAD,omitempty"`
	PV          *PowerFlowStatus        `json:"PV,omitempty"`
	Storage     *StoragePowerFlowStatus `json:"STORAGE,omitempty"`
}

type OverviewEnergy struct {
	Energy float64 `json:"energy"`
}

type OverviewPower struct {
	Power float64 `json:"power"`
}

type OverviewData struct {
	LastUpdateTime SETime         `json:"lastUpdateTime"`
	LifetimeData   OverviewEnergy `json:"lifeTimeData"`
	LastYearData   OverviewEnergy `json:"lastYearData"`
	LastMonthData  OverviewEnergy `json:"lastMonthData"`
	LastDayData    OverviewEnergy `json:"lastDayData"`
	CurrentPower   OverviewPower  `json:"currentPower"`
	MeasuredBy     string         `json:"measuredBy"`
}

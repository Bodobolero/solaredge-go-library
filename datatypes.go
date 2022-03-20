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

// SETime supports the datetime format from solaredge to be interpreted as a normal go
// time. You should set the `SiteZone` variable otherwise the zone of the current system
// will be used. SolarEdge sends datetimes in the zone of the site.
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

//  A Site contains the stored site information, like location, etcpp.
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

// Inventory lists the different systems available at a site.
type Inventory struct {
	Meters    []Meter    `json:"meters,omitempty"`
	Sensors   []Sensor   `json:"sensors,omitempty"`
	Gateways  []Gateway  `json:"gateways,omitempty"`
	Batteries []Battery  `json:"batteries,omitempty"`
	Inverters []Inverter `json:"inverters,omitempty"`
}

// Meter data
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

// Sensor data
type Sensor struct {
	Category                   string `json:"category,omitempty"`
	Type                       string `json:"type,omitempty"`
	ConnectedTo                string `json:"connectedTo,omitempty"`
	ConnectedSolaredgeDeviceSN string `json:"connectedSolaredgeDeviceSN,omitempty"`
}

// Gateway data
type Gateway struct {
	Name            string `json:"name,omitempty"`
	SerialNumber    string `json:"serialNumber,omitempty"`
	FirmwareVersion string `json:"firmwareVersion,omitempty"`
}

// Battery data.
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

// Inverter data.
type Inverter struct {
	Name                string `json:"name,omitempty"`
	Manufacturer        string `json:"manufacturer,omitempty"`
	Model               string `json:"model,omitempty"`
	CommunicationMethod string `json:"communicationMethod,omitempty"`
	CPUVersion          string `json:"cpuVersion,omitempty"`
	SN                  string `json:"SN,omitempty"`
	ConnectedOptimizers int    `json:"connectedOptimizers,omitempty"`
}

// StorageBatteryTelemetry contains telemetry data of the battery.
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

// StorageBattery data.
type StorageBattery struct {
	Nameplate   float64                   `json:"nameplate,omitempty"`
	SN          string                    `json:"serialNumber,omitempty"`
	ModelNumber string                    `json:"modelNumber,omitempty"`
	Telemetries []StorageBatteryTelemetry `json:"telemetries,omitempty"`
}

// MeterValue.
type MeterValue struct {
	Date  SETime  `json:"date"`
	Value float64 `json:"value"`
}

// MeteredValue is a collection of MeterValue's
type MeteredValue struct {
	Type   string       `json:"type"`
	Values []MeterValue `json:"values"`
}

// PowerDetails.
type PowerDetails struct {
	TimeUnit TimeUnit       `json:"timeUnit"`
	Unit     string         `json:"unit"`
	Meters   []MeteredValue `json:"meters"`
}

// EnergyDetails
type EngergyDetails struct {
	TimeUnit TimeUnit       `json:"timeUnit"`
	Unit     string         `json:"unit"`
	Meters   []MeteredValue `json:"meters"`
}

// PowerFlowConnection shows the direction of the power flow.
type PowerFlowConnection struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// PowerFlowStatus gives a status and a current power vvalue.
type PowerFlowStatus struct {
	Status       string  `json:"status"`
	CurrentPower float64 `json:"currentPower"`
}

// StoragePowerFlowStatus gives a powerflowstatus as well as a ChargeLevel of the storage
// system (battery).
type StoragePowerFlowStatus struct {
	PowerFlowStatus
	ChargeLevel int  `json:"chargeLevel"`
	Critical    bool `json:"critical"`
}

// PowerFlow describes the current flow of power in the whole system
type PowerFlow struct {
	Unit        string                  `json:"unit"`
	Connections []PowerFlowConnection   `json:"connections,omitempty"`
	Grid        PowerFlowStatus         `json:"GRID,omitempty"`
	Load        PowerFlowStatus         `json:"LOAD,omitempty"`
	PV          *PowerFlowStatus        `json:"PV,omitempty"`
	Storage     *StoragePowerFlowStatus `json:"STORAGE,omitempty"`
}

// OverviewEnergy wraps a energy value
type OverviewEnergy struct {
	Energy float64 `json:"energy"`
}

// OverviewPower wraps a power value
type OverviewPower struct {
	Power float64 `json:"power"`
}

// OverviewData returns many overview values of the whole site.
type OverviewData struct {
	LastUpdateTime SETime         `json:"lastUpdateTime"`
	LifetimeData   OverviewEnergy `json:"lifeTimeData"`
	LastYearData   OverviewEnergy `json:"lastYearData"`
	LastMonthData  OverviewEnergy `json:"lastMonthData"`
	LastDayData    OverviewEnergy `json:"lastDayData"`
	CurrentPower   OverviewPower  `json:"currentPower"`
	MeasuredBy     string         `json:"measuredBy"`
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/ulrichSchreiner/solaredge"
)

var (
	listen   string
	flow     time.Duration
	poll     time.Duration
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "starts a http service for a site",
		Run: func(cmd *cobra.Command, args []string) {
			siteid := viper.GetString("siteid")
			if len(args) > 0 {
				siteid = args[0]
			}
			serveService(siteid)
		},
	}

	pvGauge      prometheus.Gauge
	gridGauge    prometheus.Gauge
	batteryGauge prometheus.Gauge
	socGauge     prometheus.Gauge
)

func init() {
	serveCmd.PersistentFlags().StringVar(&listen, "listen", "localhost:7777", "the listen address for the service")
	serveCmd.PersistentFlags().DurationVar(&flow, "flow", 60*time.Second, "the poll duration for the powerflow call")
	serveCmd.PersistentFlags().DurationVar(&poll, "poll", 15*time.Minute, "the poll duration for standard API calls")
}

type solaredgeService struct {
	lock             sync.RWMutex
	site             *solaredge.SiteClient
	flowTimer        time.Duration
	pollTimer        time.Duration
	currentPowerFlow solaredge.PowerFlow
	currentOverview  solaredge.OverviewData
	staticDetails    solaredge.Site
}

func newSolaredgeService(sc *solaredge.SiteClient) (*solaredgeService, error) {
	res := &solaredgeService{
		site:      sc,
		flowTimer: flow,
		pollTimer: poll,
	}

	if err := res.fetchSiteDetails(); err != nil {
		return nil, err
	}

	pvGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: strings.ToLower(fmt.Sprintf("site_%d", res.staticDetails.Id)),
		Subsystem: "pv",
		Name:      "current_power",
		Help:      "the current power of the pv",
	})
	prometheus.MustRegister(pvGauge)
	gridGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: strings.ToLower(fmt.Sprintf("site_%d", res.staticDetails.Id)),
		Subsystem: "grid",
		Name:      "current_power",
		Help:      "the current power of the grid",
	})
	prometheus.MustRegister(gridGauge)
	batteryGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: strings.ToLower(fmt.Sprintf("site_%d", res.staticDetails.Id)),
		Subsystem: "battery",
		Name:      "current_power",
		Help:      "the current power of the battery",
	})
	prometheus.MustRegister(batteryGauge)
	socGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: strings.ToLower(fmt.Sprintf("site_%d", res.staticDetails.Id)),
		Subsystem: "soc",
		Name:      "current_value",
		Help:      "the current state of charge of the battery",
	})
	prometheus.MustRegister(socGauge)

	http.HandleFunc("/powerflow", res.sitePowerFlow)
	http.HandleFunc("/flow", res.siteFlow)
	http.HandleFunc("/overview", res.siteOverview)
	http.HandleFunc("/details", res.siteDetails)

	http.Handle("/metrics", promhttp.Handler())

	go res.start()
	return res, nil
}

func (ses *solaredgeService) listen(l string) {
	_ = http.ListenAndServe(listen, nil)
}

func (ses *solaredgeService) fetchPowerFlow() {
	ses.lock.Lock()
	defer ses.lock.Unlock()
	det, err := ses.site.PowerFlow()
	if err != nil {
		log.Error().Err(err).Msg("cannot query powerflow")
	} else {
		log.Info().
			Interface("powerflow", *det).
			Msg("fetched new powerflow")
		ses.currentPowerFlow = *det
	}

	fd := genFlowData(ses.currentPowerFlow)
	pvGauge.Set(fd.PV)
	gridGauge.Set(fd.Grid)
	batteryGauge.Set(fd.Battery)
	socGauge.Set(fd.SoC)
}

func (ses *solaredgeService) fetchOverview() {
	ses.lock.Lock()
	defer ses.lock.Unlock()
	det, err := ses.site.Overview()
	if err != nil {
		log.Error().Err(err).Msg("cannot query overview")
	} else {
		log.Info().
			Interface("overview", *det).
			Msg("fetched new overview")
		ses.currentOverview = *det
	}
}

func (ses *solaredgeService) fetchSiteDetails() error {
	ses.lock.Lock()
	defer ses.lock.Unlock()
	det, err := ses.site.Details()
	if err != nil {
		return fmt.Errorf("cannot query site details: %w", err)
	} else {
		log.Info().
			Interface("details", *det).
			Msg("fetched new site details")
		ses.staticDetails = *det
	}
	return nil
}

func (ses *solaredgeService) start() {
	flowtick := time.Tick(ses.flowTimer)
	polltick := time.Tick(ses.pollTimer)

	// first initialize our state
	ses.fetchPowerFlow()
	ses.fetchOverview()

	for {
		select {
		case <-flowtick:
			ses.fetchPowerFlow()
		case <-polltick:
			ses.fetchOverview()
		}
	}
}

func (ses *solaredgeService) sitePowerFlow(rw http.ResponseWriter, rq *http.Request) {
	ses.lock.RLock()
	defer ses.lock.RUnlock()

	rw.Header().Add("content-type", "application/json")
	_ = json.NewEncoder(rw).Encode(ses.currentPowerFlow)
}

func (ses *solaredgeService) siteFlow(rw http.ResponseWriter, rq *http.Request) {
	ses.lock.RLock()
	defer ses.lock.RUnlock()

	rw.Header().Add("content-type", "application/json")
	_ = json.NewEncoder(rw).Encode(genFlowData(ses.currentPowerFlow))
}

func (ses *solaredgeService) siteOverview(rw http.ResponseWriter, rq *http.Request) {
	ses.lock.RLock()
	defer ses.lock.RUnlock()

	rw.Header().Add("content-type", "application/json")
	_ = json.NewEncoder(rw).Encode(ses.currentOverview)
}

func (ses *solaredgeService) siteDetails(rw http.ResponseWriter, rq *http.Request) {
	ses.lock.RLock()
	defer ses.lock.RUnlock()

	rw.Header().Add("content-type", "application/json")
	_ = json.NewEncoder(rw).Encode(ses.staticDetails)
}

func serveService(siteid string) {
	sic, err := solaredge.SiteFromIDs(viper.GetString("apikey"), siteid, solaredge.WithBaseURL(viper.GetString("baseurl")))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create client")
	}

	srv, err := newSolaredgeService(sic)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start solaredge service")
	}
	srv.listen(listen)
}

func unitFactor(unit string) float64 {
	switch strings.ToLower(unit) {
	case "w":
		return 1.0
	case "kw":
		return 1000
	case "mw":
		return 1000000
	}
	return 1.0
}

func flowDirection(c solaredge.PowerFlowConnection) (float64, bool) {
	from := strings.ToLower(c.From)
	to := strings.ToLower(c.To)

	if from == "load" && to == "grid" {
		return -1, true
	}
	if from == "grid" && to == "load" {
		return 1, true
	}
	return 0, false
}

type flowdata struct {
	PV      float64 `json:"pv"`
	Grid    float64 `json:"grid"`
	Battery float64 `json:"battery"`
	SoC     float64 `json:"soc"`
}

func genFlowData(pf solaredge.PowerFlow) flowdata {
	battscale := -1.0
	unitscale := unitFactor(pf.Unit)

	var res flowdata
	if pf.PV != nil {
		res.PV = pf.PV.CurrentPower * unitscale
	}
	if pf.Storage != nil {
		if pf.Storage.Status == "Discharging" || pf.Storage.Status == "Idle" {
			battscale = 1.0
		}
		res.Battery = pf.Storage.CurrentPower * battscale * unitscale
		res.SoC = float64(pf.Storage.ChargeLevel)
	}
	for _, c := range pf.Connections {
		if fact, ok := flowDirection(c); ok {
			res.Grid = pf.Grid.CurrentPower * fact * unitscale
			break
		}
	}
	return res
}

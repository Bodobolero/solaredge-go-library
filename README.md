# solaredge

simple library/tool to query solaredge webservice API

forked from Ulrich Schreiner's

https://gitlab.com/ulrichSchreiner/solaredge

## Usage

Export your API-KEY and your SITE-ID:
~~~
export SOLAREDGE_APIKEY=123123123123
export SOLAREDGE_SITEID=987654
~~~

Use the `help` command to see the available subcommands:

~~~
❯ solaredge help
solaredge is a client for the solaredge webservice API

Usage:
  solaredge [flags]
  solaredge [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  site        site related actions

Flags:
      --apikey string     Your API key
      --baseurl string    The base URL for the webservices (default "https://monitoringapi.solaredge.com")
  -h, --help              help for solaredge
      --timezone string   The timezone to use for timestamps (default "CET")

Use "solaredge [command] --help" for more information about a command.
~~~

At the moment there are some site specific commands:

~~~
❯ solaredge help site
site related actions

Usage:
  solaredge site [flags]
  solaredge site [command]

Available Commands:
  details       query site details
  energydetails query energy details
  inventory     query site inventory
  powerdetails  query power details
  storagedata   query battery storage data

Flags:
  -h, --help            help for site
      --siteid string   your site id to query

Global Flags:
      --apikey string     Your API key
      --baseurl string    The base URL for the webservices (default "https://monitoringapi.solaredge.com")
      --timezone string   The timezone to use for timestamps (default "CET")

Use "solaredge site [command] --help" for more information about a command.
~~~

To query energydetails for the last 30 minutes you can do

~~~
❯ solaredge site energydetails --since 30m
{
  "timeUnit": "QUARTER_OF_AN_HOUR",
  "unit": "Wh",
  "meters": [
    {
      "type": "SelfConsumption",
      "values": [
        {
          "date": "2022-03-18 16:45:00",
          "value": 151
        },
        {
          "date": "2022-03-18 17:00:00",
          "value": 89
        },
        {
          "date": "2022-03-18 17:15:00",
          "value": 0
        }
      ]
    },
    {
      "type": "Production",
      "values": [
        {
 ...
~~~

Or some power details:

~~~
❯ solaredge site powerdetails --since 60m
{
  "timeUnit": "QUARTER_OF_AN_HOUR",
  "unit": "W",
  "meters": [
    {
      "type": "FeedIn",
      "values": [
        {
          "date": "2022-03-18 16:15:00",
          "value": 2037.3552
        },
        {
          "date": "2022-03-18 16:30:00",
          "value": 1607.2875
        },
        {
          "date": "2022-03-18 16:45:00",
          "value": 928.3824
        },
        {
          "date": "2022-03-18 17:00:00",
          "value": 516.12494
...
~~~

To query specific values, you can use `jq`:
~~~
❯ solaredge site powerflow | jq
{
  "unit": "kW",
  "connections": [
    {
      "from": "LOAD",
      "to": "Grid"
    },
    {
      "from": "PV",
      "to": "Load"
    }
  ],
  "GRID": {
    "status": "Active",
    "currentPower": 0.23
  },
  "LOAD": {
    "status": "Active",
    "currentPower": 0.6
  },
  "PV": {
    "status": "Active",
    "currentPower": 0.83
  },
  "STORAGE": {
    "status": "Idle",
    "currentPower": 0,
    "chargeLevel": 74,
    "critical": false
  }
}
❯ solaredge site powerflow | jq .STORAGE.chargeLevel
74
~~~

## Data service

You can start the embedded server to publish some data via http as JSON values.
So starting the service with
~~~
❯ export SOLAREDGE_SITEID=xxxxx
❯ solaredge serve
~~~
Runs a daemon which fetches some data from solaredge regularly. The powerflow
data is fetched every 60sec while the overview is fetched only every 15min. You
can change these intervalls as parameters.

To query the data you can do a simple GET request:
~~~
❯ xh localhost:7777/flow
HTTP/1.1 200 OK
Content-Length: 41
Content-Type: application/json
Date: Sat, 26 Mar 2022 18:57:09 GMT

{
    "pv": 0,
    "grid": 0,
    "battery": 900,
    "soc": 50
}
~~~

And you will receive the values

 - `pv`<br>
   the generated power
 - `grid`<br>
   the power which goes into the grid or is consumed
 - `battery`<br>
   the current power of the battery
 - `soc`<br>
   the state of charge of the battery
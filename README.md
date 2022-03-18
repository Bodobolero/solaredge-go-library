# solaredge

simple library/tool to query solaredge webservice API

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

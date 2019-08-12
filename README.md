# brewery3

[![Build Status](https://travis-ci.org/mkuchenbecker/brewery3.svg?branch=master)](https://travis-ci.org/mkuchenbecker/brewery3)

[![Coverage Status](https://coveralls.io/repos/github/mkuchenbecker/brewery3/badge.svg?branch=master)](https://coveralls.io/github/mkuchenbecker/brewery3?branch=master)

brewery3 is a project to build a modular microservice RaspberryPi brewery.

It's the third set of brewery software I've built. Brewery1 was an arduino brewery,
and brewery2 was a RaspberryPi python brewery with a Django frontend.

The brewery software is designed so each component lives on its own service. The `Brewery` server makes GRPC calls to `Thermometer` to get temeratures and `Switch` heating elements.

Information is logged to InfluxDb where it can then be used with Grafana.

## Major TODOs

- Spin up the InfluxDb server to ingest logs and metrics.
- Get Grafana dashboard up.
- Create GUI (currently controlled by CLI).

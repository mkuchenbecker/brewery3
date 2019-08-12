# brewery3

[![Build Status](https://travis-ci.org/mkuchenbecker/brewery3.svg?branch=master)](https://travis-ci.org/mkuchenbecker/brewery3)

[![Coverage Status](https://coveralls.io/repos/github/mkuchenbecker/brewery3/badge.svg?branch=master)](https://coveralls.io/github/mkuchenbecker/brewery3?branch=master)

brewery3 is a project to build a modular microservice RaspberryPi brewery.

It's the third set of brewery software I've built. Brewery1 was an arduino brewery,
and brewery2 was a RaspberryPi python brewery with a Django frontend.

The brewery software is designed so each component lives on its own service, and the services are managed via Kubernetes. The `Brewery` server makes GRPC calls to `Thermometer` to get temeratures and `Switch` heating elements.

The "goal" for using kubernetes is I can flash a small raspberry pi and have it register to the brewpi master node. Once it's registered, it will automatically get the payload to run. In addition, I can then easilly split each component, such as sensors and elements, onto their own RPi's that are controlled by replicated `brewery` services.

The ACTUAL goal for using kubernetes is so I can use kubernetes, and using it on this project is mostly a silly frill.

Information is logged to InfluxDb where it can then be used with Grafana.

## Major TODOs

- Spin up the InfluxDb server to ingest logs and metrics.
- Get Grafana dashboard up.
- Create GUI (currently controlled by CLI).

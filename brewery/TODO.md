Coding:
* get to 90% test coverage
* enforce 90% test coverage on master
* enable linters on pre-merge
* make brewery into own repo?


Persistance layer:

* Persist config to (remote?) database and retrieve.
* Persist all temperature reads to database

Visualization:
* Dashboard to visualize temperatures (grafana)?

Deploys:

* Deploy to RPI via Kubernetes
* Deploy serverless code to Google Cloud

Experiment/ Fun:

* Use redis cluster to cache current brewery setup


Stabilize

* Make all GRPC calls retry.
* Make the number of retries a constant.
* Make an implementation to the element interface that locks, unlocks, and turns itself off.


Hardware:

* Get DS18B20 running.
* Get holes cut in box for ports.

Feature work:

* Pumps
* Pump config
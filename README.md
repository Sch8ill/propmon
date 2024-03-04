# propmon

[![Release](https://img.shields.io/github/release/sch8ill/propmon.svg?style=flat-square)](https://github.com/sch8ill/propmon/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/sch8ill/propmon)](https://goreportcard.com/report/github.com/sch8ill/propmon)
![MIT license](https://img.shields.io/badge/license-MIT-green)

---

`propmon` is a Prometheus metrics exporter designed to monitor service proposals within
the [Mysterium Network](https://mysterium.network). Each [node](https://github.com/mysteriumnetwork/node) in
the [Mysterium Network](https://mysterium.network) advertises its services by transmitting service proposals to a
message broker through [nats](https://github.com/nats-io).
The [discovery](https://github.com/mysteriumnetwork/discovery) service captures these service proposals and compiles
them into a list, which is then made accessible to Mysterium Network clients through
a [REST API](https://discovery.mysterium.network/api/v4/proposals).

Similar to the [discovery](https://github.com/mysteriumnetwork/discovery) service, `propmon` listens for incoming
service proposals. It stores these proposals and generates metrics related to them, which are exposed on
port [9500](http://localhost:9500/metrics).

A service proposal includes:

- Provider ID of the node
- Service type
- Approximate GEO-location of the node
- Access policies for the service
- Compatibility information
- ...

---

## Installation

### Docker

```bash
docker run -p 9500:9500 sch8ill/propmon:latest
```

### Build

Requires:

```
go >= 1.21
make
```

Build command:

```bash
make build
```

---

## Usage

### Prometheus config

Example `prometheus.yml` scrape config:

```yaml
scrape_configs:
  - job_name: propmon
    scrape_interval: 15s
    static_configs:
      - targets: [ "localhost:9500" ]
```

### Metrics

| name                          | description                               | labels             | type    |
|-------------------------------|-------------------------------------------|--------------------|---------|
| propmon_proposal_ping         | Service Proposal ping                     |                    | counter |
| propmon_proposal_registered   | Service Proposal registered               |                    | counter |
| propmon_proposal_unregistered | Service Proposal unregistered             |                    | counter |
| propmon_proposal_expired      | Service Proposal expired                  |                    | counter |
| propmon_proposal_invalid      | Service Proposal invalid                  |                    | counter |
| propmon_proposal_count        | Service Proposal count                    | service_type       | gauge   |
| propmon_provider_count        | Provider count                            | country, node_type | gauge   |
| propmon_nats_bytes_rx         | Number of bytes received by NATS listener | subject            | counter |

### CLI flags

```
   --broker-address value        broker address to listen for proposals (default: "nats://broker.mysterium.network:4222")
   --proposal-lifetime value     lifetime of a proposal until it expires if not renewed (default: 3m10s)
   --expiration-job-delay value  delay between expiration job runs (default: 20s)
   --metrics-address value       address the prometheus metrics exporter listens on (default: ":9500")
   --help, -h                    show help
```

## License

This package is licensed under the [MIT License](LICENSE).

---

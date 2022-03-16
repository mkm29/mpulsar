# Hermes

[![Go 1.19](https://img.shields.io/badge/Go-v1.18-blue)](https://golang.org/doc/go1.18)

```yaml
Author: Mitch Murphy
Date: 2022 March 16
```

## Introduction

This project is a lightning fast, highly scalable API written in Golang (intended to be deployed to a Kubernetes environment) that:

  * Publishes/subscribes to topics on a [Pulsar](https://pulsar.apache.org/) server  
  * Processes message in real time to add geo metadata  
  * Geospatial data is provided by Pelias server  
  * Geohashes are added to metadata  
  * Data is inserted/updated in Cassandra  

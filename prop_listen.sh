#!/bin/bash

if ! command -v nats &> /dev/null
then
    echo "nats cli not found"
    echo "try: go install github.com/nats-io/natscli/nats@latest"
    exit 1
fi

nats sub -s nats://broker.mysterium.network:4222 "*.proposal-ping.v3"
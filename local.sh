#!/bin/bash

# Pull the latest version of the image
docker pull ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib:latest

# Run the container with the provided config file
docker run -v "${PWD}/config.yaml:/etc/otel/config.yaml" -p 4317:4317 -p 12000:12000 -p 12001:12001 --name otel_contrib ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib:latest --config /etc/otel/config.yaml
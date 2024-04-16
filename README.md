# Dropsonde

Dropsonde is a terminal UI designed to complement the [OpenTelemetry
Collector](https://github.com/open-telemetry/opentelemetry-collector-contrib)
. It leverages
the [Remote Tap Processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/remotetapprocessor)
and [OpAMP Extension](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/extension/opampextension).

## Features

- Real-time output of signals from Remote Tap (Metrics, Traces, Logs)
- View and edit Collector configuration via OpAMP

## Requirements

- OpenTelemetry Collector with Remote Tap Processor + OpAMP Extension

## Usage

Configure your Collector pipeline as seen in the `config.yaml` file in this 
repository. 

Launch dropsonde with the `-tap-endpoint` flag set to your Remote Tap endpoint.

An example of the appropriate OpAMP configuration is also in `config.yaml`.
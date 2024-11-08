# Sequence Pipeline

Sequence pipeline process data from the given data sources to a data warehouse table (BigQuery), transform and load the data into a table using Go and GCP. Its extract, normalize, and calculate daily marketplace volume, daily transactions, and aggregated volume data, and send it to an API endpoint for data visualization.

## Table of Contents
- [Table of Contents](#table-of-contents)
- [Overview](#overview)
- [Installation](#installation)
  - [Manual installation](#manual-installation)
- [Usage](#usage)
  - [Examples](#examples)
  - [Running the pipeline cli](#running-the-pipeline-cli)
  - [Running the pipeline cli locally](#running-the-pipeline-cli-locally)
  - [Running the pipeline K8s](#running-the-pipeline-k8s)
- [Enhancement](#enhancement)
- [Contributing](#contributing)


## Overview

The current architecture of the pipeline is described in the [ARCHITECTURE.md](./ARCHITECTURE.md) document.

For more in-depth explanations and considerations on architectural choices for the service, please refer to our [Architecture Decision Records](./resources/adr) folder.

If you want to submit an architectural change to the service, please create a new entry in the ADR folder [using the template provided](./resources/adr/template.md) and open a new Pull Request for review. Each ADR should have a prefix with the consecutive number and a name. For example `002-implement-server-streaming.md`

## Installation

Before to use the pipeline, you must have Go installed and configured properly on the computer. Please see [https://golang.org/doc/install](https://golang.org/doc/install)

```shell
go get github.com/dohernandez/horizon-blockchain-games
```

Or download binary from [releases](https://github.com/dohernandez/horizon-blockchain-games/releases).

[[table of contents]](#table-of-contents)

#### Manual installation

Clone the repository and build the binary by running the command

```shell
make build
```

[[table of contents]](#table-of-contents)

the binary will be located in the folder `bin` at the root of the project.

## Usage

The **file** with the data to process is a required and it is controller with the flag`--file`; the **dir** or **bucket** where the data is located, is controller with the flag`--dir`.

The step `calculator` can be executed in parallel by setting the flag `--workers`. Default value is 1.

```shell
% bin/sequence help
NAME:
   sequence - Run a sequence of steps to process data

USAGE:
   sequence [global options] command [command options]

COMMANDS:
   run      
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

```

```shell
% bin/sequence run help
NAME:
   sequence run

USAGE:
   sequence run [command options]

DESCRIPTION:
   Run pipeline, or a specific step depending on options

OPTIONS:
   --env value                     environment (default: dev) [$ENVIRONMENT]
   --extractor, -e                 run only pipeline step extractor (default: false) [$EXTRACTOR_ENABLED]
   --calculator, -c                run only pipeline step calculator (default: false) [$CALCULATOR_ENABLED]
   --insertion, -i                 run only pipeline step insertion (default: false) [$INSERTION_ENABLED]
   --all, -a                       run all pipeline steps (default: true)
   --workers value, -w value       number of workers to run the pipeline (default: 1) [$CALCULATOR_WORKERS]
   --dir value                     folder or bucket to read/store the intermediate step data when required (default: 2024-11-08) [$DIR, $BUCKET]
   --file value                    file to read the data from (default: transactions.csv) [$FILE, $DATA_FILE]
   --test                          run the pipeline in test mode using local file system as providers (default: false)
   --conversor value               conversor to use to convert the currency [coingecko hardcoded] (default: coingecko)
   --coingecko-api-key-type value  API key type to use with the coingecko conversor (default: x_cg_demo_api_key) [$CG_API_KEY_TYPE]
   --coingecko-api-key value       API key to use with the coingecko conversor [$CG_API_KEY]
   --storage-type value            storage type to use to load/store [file bucket] (default: bucket)
   --verbose, -v                   enable verbose output (default: false) [$VERBOSE]
   --warehouse value               target type to use to load/store [print bigquery] (default: bigquery)
   --bigquery-dataset value        BigQuery dataset in the following format <project_id.dataset.table> [$BIGQUERY_DATASET]
   --help, -h                      show help
```

[[table of contents]](#table-of-contents)

#### Examples

1. To use the pipeline in testing mode.
```shell
bin/sequence run --dir ./resources/sample-bucket --file sample_data.csv --test
```

2. To use the pipeline split across multiple runs.
```shell
bin/sequence run --calculator -w 10 --dir ./resources/sample-bucket --file sample_data.csv --test
```

[[table of contents]](#table-of-contents)

#### Running the pipeline cli

The pipeline can be executed simultaneously in a single run or split across multiple runs, the default to run all the steps in a single run. To run the pipeline split across multiple runs use any of the flag (`--extractor`, `--calculator`, `--insertion`).

To configure the pipeline to use different storage such `bucket` to load the data to process from GCP, it is required to export `GOOGLE_APPLICATION_CREDENTIALS=/path/to/sa-json` and use the flag `--storage-type`. Default value is `bucket`.

To configure the pipeline to use different warehouse such `BigQuery` to save the output of the step, it is required to export `GOOGLE_APPLICATION_CREDENTIALS=/path/to/sa-json` and use the flag `--warehouse`. Default value is `bigquery`.

The pipeline can be also configurable to use different conversor such as `coingecko` to convert the currency. When using  `coingecko` convertor (flag `--conversor` default value is `coingecko`), it is required to export `CG_API_KEY=<ch-api-key>` and set the flag `--coingecko-api-key-type` to `x_cg_demo_api_key` or `x-cg-pro-api-key` depending on the API key type. (default `x_cg_demo_api_key`). 

The pipeline can be also run in test mode using the local file system as a provider and output the result to the os.Stdout.

[[table of contents]](#table-of-contents)


#### Running the pipeline cli locally

To run the pipeline locally, using `fsouza/fake-gcs-server` to emulate GCS storage bucket:

Spin up docker-compose services:

```shell
docker-compose up
```

Export the additional following environment variables:

```shell
export GCP_BUCKET_ENDPOINT=http://storage.gcs.127.0.0.1.nip.io:4443/storage/v1/
export STORAGE_EMULATOR_HOST=http://localhost:4443
```

Run the pipeline:

```shell
bin/sequence run --dir sample-bucket --file sample_data.csv --verbose
```

[[table of contents]](#table-of-contents)

#### Running the pipeline K8s

To run the pipeline in a K8s cluster, you must build the docker image to ship it to a docker registry.

```shell
make docker-build
```

## Enhancement

* Improve test suite.
* Improve error handling.
* Improve obliquity language.
* Add local BigQuery emulator.

[[table of contents]](#table-of-contents)

## Contributing

Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

[[table of contents]](#table-of-contents)
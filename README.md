# gcloud-ssh

`gcloud-ssh` is a Go-based command-line tool that simplifies the process of connecting to Google Cloud Platform (GCP) bastion servers using the `gcloud` command. It supports multiple environments and SOCKS proxy configuration.

## Features

- Easy connection to GCP bastion servers across different environments
- YAML-based configuration for managing multiple environments
- Support for default values to minimize repetitive configuration
- Optional SOCKS proxy support
- Customizable configuration file path

## Prerequisites

- `gcloud` command-line tool installed and configured
- Access to GCP projects and compute instances

## Installation

Build for tool:

```
$ go build
```

## Configuration

Create a YAML configuration file at `~/.config/gcloud-ssh.yaml` with the following structure:

```yaml
defaults:
  host: bastion
  zone: us-central1-a
  user: your-username
  socks_port: 8005

environments:
  dev:
    project: your-dev-project-id
  staging:
    project: your-staging-project-id
    zone: us-west1-b
  prod:
    project: your-prod-project-id
    host: prod-bastion
    user: prod-user
    socks_port: 8010
```

Adjust the values according to your GCP setup.

## Usage

Basic usage:

```
$ gcloud-ssh -env <environment_name>
```

With SOCKS proxy:

```
$ gcloud-ssh -env <environment_name> -socks
```

Using a custom configuration file:

```
$ gcloud-ssh -env <environment_name> -config /path/to/custom/config.yaml
```

## Examples

Connect to the development environment:

```
$ gcloud-ssh -env dev
```

Connect to the production environment with SOCKS proxy:

```
$ gcloud-ssh -env prod -socks
```

## License

This project is licensed under the [MIT License](./LICENSE).

# Ketoz 

- [Introduction](#introduction)
- [Features](#features)
- [Getting Started](#getting-started)
    - [Installation](#installation)
    - [Configuration](#configuration)
- [Usage](#usage)
    - [API Reference](#api-reference)
- [License](#license)

## Introduction

Ketoz is a microservice that extends Ory Keto to provide fine-grained, hierarchical role-based access control (HRBAC) for modern applications. It enables flexible permission management and scalable authorization policies.

## Features

- Hierarchical roles and permissions
- Fine-grained access control
- RESTful APIs

## Getting Started

### Installation
1. Pull the Docker image:
    ```sh
    docker pull kimnguyenlong/ketoz:latest
    ```
2. Run the service:
    ```sh
    docker run -d \
        --name ketoz \
        --env-file /path/to/.env \
        -p 8000:8000 \
        kimnguyenlong/ketoz:latest
    ```

### Configuration

#### Keto OPL

To enable Ketoz to function correctly, you must apply the [OPL](https://www.ory.sh/docs/keto/reference/ory-permission-language) from the [`keto/namespaces.ts`](keto/namespaces.ts) file to your Keto instance. This ensures that the required namespaces and permission structures are available for the service.

#### Environment variables

Ketoz loads its configuration from environment variables at startup. 

```env
# Service
SERVICE_HOST=0.0.0.0
SERVICE_PORT=8000
SERVICE_LOG_LEVEL=DEBUG

# Keto
KETO_HOST=keto
KETO_READ_PORT=4466
KETO_WRITE_PORT=4467
```


## Usage

### API Reference

Ketoz exposes a RESTful API for managing roles, permissions, and access policies. See the [API documentation](docs/api.md) for detailed endpoints and request/response formats.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
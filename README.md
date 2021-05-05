# Minter Hub Sentinel

Essential application for running validator in [Minter Hub](https://minter.org) network.

## Installation

### From source

```bash
git clone https://github.com/friendstrust/minter-hub-sentinel && \
cd minter-hub-sentinel && \
make build
```
### Docker

Docker image `friendstrust/minter-hub-sentinel` is available on [Docker Hub](https://hub.docker.com/r/friendstrust/minter-hub-sentinel).

`docker-compose.yml` example:

```yaml

version: '3'

services:
  app:
    image: friendstrust/minter-hub-sentinel:latest
    volumes:
      - ./config.yaml:/config.yaml
    restart: unless-stopped
    logging:
      options:
        max-size: "10m"
        max-file: "3"
    deploy:
      resources:
        limits:
          cpus: '0.50'
          memory: 50M
        reservations:
          cpus: '0.25'
          memory: 20M
```

### Binary

Download the latest release version from https://github.com/friendstrust/minter-hub-sentinel/releases

## Configuration

Create config.yaml based on the [config.example.yaml](https://github.com/FriendsTRUST/minter-hub-sentinel/blob/master/config.example.yaml) file.

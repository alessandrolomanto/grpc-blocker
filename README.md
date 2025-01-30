# Traefik Plugin - gRPC Blocker

<p align="center">
  <img src="assets/img/grpc-blocker-banner.svg" alt="gRPC Blocker Banner" width="800" />
</p>

A Traefik plugin that allows blocking specific gRPC services based on their fully qualified service names.

## Features

- Block specific gRPC services by their fully qualified names
- Optional logging of blocked requests
- Compatible with Traefik v2.x

## Configuration

### Static Configuration

To enable the plugin in your Traefik static configuration:

```yaml
# Static configuration
pilot:
  token: "xxxxx"
  experimental:
    plugins:
      grpc-blocker:
        moduleName: "github.com/alessandrolomanto/traefik-grpc-blocker"
        version: "v0.0.1-beta"
```

### Dynamic Configuration

```yaml
# Dynamic configuration
http:
  middlewares:
    my-blocker:
      plugin:
        grpc-blocker:
          enableLogging: false
          blockedServices:
            - "grpc.reflection.v1"
            - "grpc.health.v1.Health"
```

# Saber

A security management platform that integrates data collection and host control.

## Features

- **Agent**: Maintains gRPC long connections with Controller and Transfer, periodically reports host metrics (CPU, memory, etc.)
- **Transfer**: Receives streaming data and writes to Kafka
- **Controller**: Central gRPC service for agent connections

## Architecture

```
┌─────────┐     Connect (stream)      ┌────────────┐
│  Agent  │◄─────────────────────────►│ Controller │
└────┬────┘  AgentRequest/Response    └────────────┘
     │
     │ PushData (stream)
     ▼
┌─────────┐     Kafka                 ┌───────┐
│ Transfer│──────────────────────────►│ Kafka │
└─────────┘  TransferRequest payload  └───────┘
```

- **Agent**: Host agent that maintains gRPC long connections with Controller (keepalive/heartbeat) and Transfer, periodically reporting host metrics (CPU, memory, etc.)
- **Controller**: gRPC service that accepts Agent connections via `Connect(stream AgentRequest) returns (stream AgentResponse)`
- **Transfer**: gRPC service that receives streaming data from Agent via `PushData(stream TransferRequest)` and writes the payload to Kafka
- **Kafka**: Transfer writes each metric as a message to the configured topic; ClientID can optionally be used as the message key

## Prerequisites

- Go 1.23+
- Docker, Helm, kubectl (for Kubernetes deployment)
- Kafka (optional, for Transfer sink)

## Build

```bash
make              # Generate proto and build agent, controller, transfer
make agent        # Build agent only
make controller   # Build controller only
make transfer     # Build transfer only
```

## Configuration

| Component | Config File | Key Settings |
|-----------|-------------|--------------|
| Agent | `etc/agent.yaml` | controller/transfer endpoints, collector.interval |
| Transfer | `etc/transfer.yaml` | service.listenAddress, kafka.brokers, kafka.topic |
| Controller | `etc/controller.yaml` | service.listenAddress |

## Deployment (Helm + kubectl)

Deploy controller, transfer, and agent to an existing Kubernetes cluster.

### Build image (optional, for local use)

```bash
make docker-build
# Image defaults to saber/saber:$(VERSION), override with REGISTRY
```

### Install

```bash
./deploy/install.sh
```

With a specific namespace:

```bash
./deploy/install.sh -n my-ns
```

With custom values (e.g., Kafka brokers, image):

```bash
./deploy/install.sh -n my-ns -f deploy/helm/saber/values-prod.yaml
```

The install script waits for controller, transfer, and agent Deployments to be ready, then prints `kubectl get svc,pods` examples.

### Uninstall

```bash
helm uninstall saber -n <namespace>
```

Chart and configuration details: `deploy/helm/saber/`. Controller and Transfer are exposed only via cluster-internal Services; use Ingress or LoadBalancer for external access.

## License

See [LICENSE](LICENSE) for details.

# saber

这是一个可以让自己念头通达的项目！

## 架构

- **Probe**：主机探针，与 Controller 保持 gRPC 长连接（保活/心跳），与 Transfer 保持 gRPC 长连接并定时上报主机性能数据（CPU、内存等）。
- **Controller**：gRPC 服务，接收 Probe 的 `Connect(stream ProbeRequest) returns (stream ProbeResponse)` 长连接。
- **Transfer**：gRPC 服务，接收 Probe 的 `PushData(stream TransferRequest)` 流式数据，将收到的 Payload 写入 Kafka。
- **Kafka**：Transfer 将每条性能数据作为一条消息写入配置的 Topic，可选以 ClientID 为 key。

## 构建

```bash
make          # 生成 proto 并构建 probe、controller、transfer
make probe    # 仅构建 probe
make controller
make transfer
```

## 配置

- Probe：`etc/probe.yaml`（controller/transfer 的 endpoints、collector.interval 等）。
- Transfer：`etc/transfer.yaml`（service.listenAddress、kafka.brokers、kafka.topic 等）。
- Controller：`etc/controller.yaml`（service.listenAddress 等）。

## 部署（Helm + kubectl）

在已有 Kubernetes 集群中一键部署 controller、transfer、probe。

**依赖**：Docker、Helm、kubectl、可用的 Kubernetes 集群。

**构建镜像**（可选，若使用本地镜像）：

```bash
make docker-build
# 镜像默认为 saber/saber:$(VERSION)，可通过 REGISTRY 覆盖
```

**一键安装**：

```bash
./deploy/install.sh
```

指定命名空间：

```bash
./deploy/install.sh -n my-ns
```

使用自定义 values 覆盖（如 Kafka brokers、镜像等）：

```bash
./deploy/install.sh -n my-ns -f deploy/helm/saber/values-prod.yaml
```

安装后脚本会等待 controller、transfer（及 probe）Deployment 就绪，并输出 `kubectl get svc,pods` 示例。

**卸载**：

```bash
helm uninstall saber -n <namespace>
```

Chart 与配置说明见 `deploy/helm/saber/`；Controller/Transfer 仅集群内 Service 暴露，需对外访问时可自行配置 Ingress 或 LoadBalancer。

# 功能列表
- Probe：双 gRPC 长连接 + 主机性能采集上报
- Transfer：接收数据并写入 Kafka

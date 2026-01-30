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

# 功能列表
- Probe：双 gRPC 长连接 + 主机性能采集上报
- Transfer：接收数据并写入 Kafka

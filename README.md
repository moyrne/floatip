# floatip

[**中文**](README.md) | [**English**](README.en.md)

自动同步当前主机的公网 IP 到云防火墙模板，并应用到指定实例。

## 工作原理

1. **检测 IP** — 定时请求 IP 回显服务获取当前公网 IP
2. **同步规则** — 在指定的防火墙模板中查找或更新一条规则，来源 IP 设为当前公网 IP
3. **应用模板** — 将模板应用到目标实例，使实例防火墙与模板保持一致

## 快速开始

```bash
make build    # 编译
make run      # 运行（需 config.yaml）
make test     # 测试
```

## Docker

```bash
docker build -t floatip .
docker run -v $(pwd)/config.yaml:/etc/floatip/config.yaml floatip
```

## 部署

Kubernetes 部署清单位于 `deploy/` 目录。

## 当前适配

| 云厂商 | 状态 |
|--------|------|
| 腾讯云轻量云服务器 | ✅ 已适配 |

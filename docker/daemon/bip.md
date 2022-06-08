# 修改docker默认网卡网段

## 查看当前配置

```go
ifconfig // 查看以下字段
	docker0
```

## 修改配置

```go
sudo vim /etc/docker/daemon.json

{
  "bip": "10.101.0.1/16"
}

```

## 重启验证

```go
systemctl restart docker
ifconfig // 检查
```
# 修改docker默认存储路径

## 查看当前配置

```go
docker info // 查看以下字段
	Docker Root Dir
```

## 修改配置

```go
vi /etc/docker/daemon.json

{
  "graph": "/data/docker"
}

// 版本比较旧的字段可能是 data-root
```

## 重启验证

```go
systemctl restart docker
docker info // 检查
```


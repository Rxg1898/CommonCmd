# 修改docker日志记录驱动

## 修改配置

```go
sudo vim /etc/docker/daemon.json

{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3" 
  }
}

```

## 重启生效

```go
systemctl reload docker
systemctl restart docker
```
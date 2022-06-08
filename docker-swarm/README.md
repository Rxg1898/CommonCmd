Docker swarm常用命令	

## 1、强制重启服务

- `--force`：强制更新
- `--update-parallelism`：同时更新最大副本，默认0表示全部一起更新
- `--update-delay`：更新之间的延迟时间(ns|us|ms|s|m|h)

```
docker service update --force --update-parallelism 1 --update-delay 30s service名称
```

## 2、查看日志

- `-f`：跟踪日志输出
- `--tail`：从日志末尾倒数第几行开始打印(默认全部)

```
docker service logs -f --tail 100 service名称
```


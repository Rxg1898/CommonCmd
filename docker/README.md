Docker常用命令	

## 1、强制删除容器服务

有时候我们会遇到docker rm删除不掉的容器，这时又不想重启整个docker进程(systemctl restart docker)影响到其他正在运行的容器，这时我们可以这样操作：

### 1.1 docker维护模式(在守护进程停机期间保持容器正常运行)

官网链接：https://docs.docker.com/config/containers/live-restore/

编辑：`/etc/docker/daemon.json`添加以下内容

```
{
	"live-restore": true
}
```

然后重启docker

`systemctl reload docker`

`systemctl restart docker`

### 1.2 以上方法不行，就先删除容器containers下对应容器ID的目录

Docker Root Dir没有修改过的情况下，containers默认路径`/var/lib/docker/containers`

删除对应的容器ID目录：rm -rf 容器ID

再执行以上1.1的方法 

## 2、查看日志

- `-f`：跟踪日志输出
- `--tail`：从日志末尾倒数第几行开始打印(默认全部)

```
docker  logs -f --tail 100 容器名称｜容器ID
```

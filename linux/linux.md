## 其他

- [git命令](./git/git.md)
- [curl](./curl/README.md)

## 修改服务器时区为UTC-0时区

```
ln -sf /usr/share/zoneinfo/UTC /etc/localtime
```



## 快速删除文件

```
// 删除日志文件，只保留近7天的
find ./ -iname "*.log" -type f -mtime +7 -delete
```



## 当前服务器公网IP地址

```
curl -L ip.tool.lu

curl ip.sb
```


## 其他

- [git命令](./git/git.md)

## 修改服务器时区为UTC-0时区

```
ln -sf /usr/share/zoneinfo/UTC /etc/localtime
```



## 快速删除文件

```
// 删除日志文件，只保留近7天的
find ./ -iname "*.log" -type f -mtime +7 -delete
```


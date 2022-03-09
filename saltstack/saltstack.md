# saltstack常用命令

saltstack批量主机管理工具，安装官网地址：https://repo.saltproject.io/index.html#rhel

## 查看当前所能管理主机的key

```
salt-key
```

## 节点加入

```
// 允许node1加入
salt-key -a node1
// 允许所有加入
salt-key -A
```

## 删除节点

```
// 删除单个节点
salt-key -d node1
// 批量删除
salt-key -d "*node*"
// 删除所有节点
salt-key -D
```

## 测试当前节点连接性

```
salt "*" test.ping
```

## 远程执行命令

```
salt  '*'  cmd.run "date"
// 使用root用户，加载环境变量执行
salt  '*'  cmd.run 'date' shell='/bin/bash' runas='root'
```

## 定时任务

```
// 查看主机root用户的定时任务有哪些
salt "*"  cron.list_tab root
// 添加定时任务
salt "*" cron.set_job nginx '*/2' '*' '*' '*' '*' '/bin/bash  /opt/nginx-monitor.sh'
```

## 执行编排sls脚本

脚本路径以/etc/salt/master以下配置作为相对路径

```
file_roots:
  base:
    - /srv/salt
```

### 对pre环境的主机执行http.sls

```
salt "*pre*" state.sls http
salt "*pre*" state.sls pre.http
salt "*pre*" state.sls pre/http
```

### 对demo应用的主机执行http.sls

主机组在/etc/salt/master配置

```
nodegroups:
  demo: '*pre*'
  pre-demo: 'L@ali-hk-pre-node-00,ali-sz-pre-node-00'
  prod-demo: 'L@ali-bj-prod-node-00'
```

-N 指定主机组

-b 并发多少个主机执行

```
salt -N demo -b 5  state.sls http
salt -N pre-demo state.sls http
salt -N proe-demo -b 1  state.sls http
```

传参执行

```
salt -N pre-demo state.sls deploy  pillar='{"service": "demo", "version": 1}'
```



## 示例文件

[salt示例sls文件](./salt示例sls文件)


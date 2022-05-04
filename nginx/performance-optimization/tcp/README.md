## 性能优化方法论

- 软件层面提升硬件使用率
  - 增大CPU的利用率
  - 增大内存的利用率
  - 增大硬盘IO的利用率
  - 增大网络带宽的利用率
- 提升硬件
  - 网卡：万兆网卡
  - 硬盘：固体硬盘，关注IOPS和BPS指标
  - CPU：更快主频，更多核心，更大缓存，更优架构
  - 内存：更快访问速度
- 超出硬件性能上限后使用DNS

## TCP基本知识

### 三次握手和四次挥手

![tcp](https://blog-1301758797.cos.accelerate.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/tcp/20161127122713750))

### Tcp Fast Open

起到tcp加快通信的作用

- net.ipv4.tcp_fastopen 系统开启TFO功能
  - 0：关闭
  - 1：作为客户端时
  - 2：作为服务端时
  - 3：无论客户端或服务端，都开启TFO

- Tcp Fast Open nginx配置,指定TFO连接队列长度

```
listen addres[:port][fastopen=number];
```

### 滑动窗口

- 用于限制连接的网速，解决报文乱序和可靠传输问题
- Nginx中limit_rate等限速指令皆依赖它实现
- 由操作系统内核实现
- 连接两端各有发送窗口(用于发送内容)和接收窗口(用于接收内容)

nginx超时指令和滑动窗口

- 两次读操作间的超时

```
client_body_timeout time; // 默认60s
```

- 两次写操作间的超时

```
send_timeout time; // 默认60s
```

- 兼具两次读写操作间的超时

```
proxy_timeout timeout; // 默认10m
```

### 缓冲区

- 读缓存最小值、默认值、最大值、单位字节，覆盖net.core.rmem_max

```
net.ipv4.tcp_rmem=4096 87380 6291456
```

- 写缓存最小值、默认值、最大值、单位字节，覆盖net.core.wmem_max

```
net.ipv4.tcp_wmem=4096 16384 4194304
```

- 系统无内存压力，启动压力模式阈值、最大值，单位为页的数量

```
net.ipv4.tcp_mem=1541646 2055528 3083292
```

- 开启自动调整缓存模式

```
net.ipv4.tcp_moderate_rcvbuf=1
```

- 调整接收窗口与应用缓存

```
net.ipv4.tcp_adv_win_scale=1

应用缓存=buffer/(2^tcp_adv_win_scale)
```

### 吞吐量=窗口/时延

### BDP=带宽*时延

### Nagle算法

作用：

1、避免一个连接上同时存在大量小报文：最多只存在一个小报文，合并多个小报文一起发送

2、提高带宽利用率

- Nagle算法在Nginx配置

```
tcp_nodelay off; // 吞吐量优先，开启Nagle算法。默认on
tcp_nodelay on; // 低时延优先，禁用Nagle算法。默认on
```

- Nginx避免发送小报文配置

```
postpone_output size; // 默认1460
```

### CORK算法

仅针对sendfile on开启时有效，完全禁止小报文的发送，提升网络效率

```
tcp_nopush on|off; // 默认off开启
```

### 拥塞窗口

什么是拥塞窗口？发送方主动限制流量

什么是通告窗口？接收方限制流量

实际流量：拥塞窗口与通告窗口的最小值

拥塞处理：

- 慢启动
  - 指数扩展拥塞窗口（cwnd为拥塞窗口大小）
    - 每收到1个ACK，cwnd=cwnd+1
    - 每过一个RTT，cwnd=cwnd*2

- 拥塞避免
  - 线性扩展拥塞窗口
    - 每收到1个ACK，cwnd=cwnd+1/cwnd
    - 每过一个RTT，窗口加1
- 拥塞发送
  - 急速降低拥塞窗口
    - RTO超时，threshold=cwnd/2，cwnd=1
    - Fast Retransmit ，收到3个duplicate ACK，cwnd=cwnd/2，threshold=cwnd
- 快速恢复
  - 当Fast Retransmit出现时，cwnd调整为threshold+3*MSS

RTT(Round Trip Time)由物理链路传输时间+末端处理时间+路由器排队处理时间组成。

RTO(Retransmission Time Out)正确的重传之前的丢包

### KeepAlive功能

tcp的keepalive和http的keepalive作用是不同的，http是为了复用，tcp则为了尽快释放。

tcp keepalive应用场景：检测实际断掉的连接和用于维持与客户端间防火墙有活跃网络包

linux系统的tcp keepalive：

- 发送心跳周期

```
net.ipv4.tcp_keepalive_time=7200 // 秒
```

- 探测包发送间隔

```
net.ipv4.tcp_keepalive_intvl=75
```

- 探测包重试次数

```
net.ipv4.tcp_keepalive_probes=9
```

Nginx的tcp keepalive：

- so_keepalive=30m::10
- Keepidlea,keepintvl,keepcnt

### time_wait状态

#### 被动关闭连接

查看连接状态，State字段

```
netstat -ant
或者ss -ant
```

- 存在大量CLOSE_WAIT状态：应用进程没有及时响应对端关闭连接
- 存在大量LAST_ACK状态：等待接收主动关闭端操作系统发来的针对FIN的ACK报文



#### 主动关闭连接

fin_wait1状态和fin_wait2状态

- fin_wait1状态，发送FIN报文的重试次数，0相当于8

```
net.ipv4.tcp_orphan_retries=0
```

- fin_wait2状态，保持在FIN_WAIT_2状态的时间

```
net.ipv4.tcp_fin_timeout=60
```

#### time_wait状态有什么用？

time_wait存在的2个原因：

1. 可靠地终止TCP连接
2. 保证让迟来的TCP报文段有足够的时间被识别并丢弃

MSL(Maximum Segment Lifetime)报文最大生存时间,tcp报文生存周期1MSL，2MSL是为了确保网络上的tcp报文已经无了。所以维持TIME_WAIT状态的时长是2MSL，保证至少一次报文的往返时间内端口是不可复用。



#### time_wait优化

- 开启后，作为客户端时新连接可以使用仍然处于time_wait状态端口

```
net.ipv4.tcp_tw_reuse=1
```

- 由于timestamp的存在，操作系统可以拒绝迟到的报文

```
net,ipv4.tcp_timestamps=1
```

- 开启后，同时作为客户端和服务器都可以使用time_wait状态的端口，不安全，无法避免报文延迟、重复等给新连接造成混乱

```
net.ipv4.tcp_tw_recycle=0
```

- Time_wait状态连接的最大数量，超出后直接关闭连接

```
net.ipv4.tcp_max_tw_buckets=262144
```



### lingering_close延迟关闭

意义：当Nginx处理完成调用close关闭连接后，若接收缓冲区仍然收到客户端发来的内容，则服务器会向客户端发送RST包关闭连接，导致客户端由于收到RST而忽略了http response

Nginx配置：

- off 关闭功能
- on由Nginx判断，当用户请求未接收完(根据chunk或者Content-Length头部等)时启用功能，否则及时关闭连接
- always无条件启用功能

```
lingering_close off|on|always; // 默认on
```

- 当功能启用时，最长的读取用户请求内容的时长，达到后立刻关闭连接

```
lingering_time time; //默认30s
```

- 当功能启用时，检测客户端是否仍然请求内容到达，若超时后仍没有数据达到，则立刻关闭连接

```
lingering_timeout time; //默认5s
```

### RST代替正常四次挥手关闭连接

- 当读、写超时指令生效引发连接关闭时产生的异常，通过发送RST立刻释放端口、内存等资源来关闭连接

```
reset_timedout_connection on|off; // 默认off
```



## SYN攻击

什么是SYN攻击：攻击者短时间伪造不同IP地址的SYN报文，快速占满backlog队列，使服务器不能为正常用户服务

### 对SYN攻击防御的内核参数配置

- 接收自网卡，但未被内核协议栈处理的报文队列长度

```
net.core.netdev_max_backlog
```

- backlog队列长度，nginx配置

```
listen addres[:port][backlog=number]
```

- SYN_RCVD状态连接的最大个数(半连接)

```
net.ipv4.tcp_max_syn_backlog
```

- 超出处理能力时，对新来的SYN直接回包RST，丢弃连接

```
net.ipv4.tcp_abort_on_overflow
```

- 当SYN队列满后，新的SYN不进入队列，计算出cookie再以SYN+ACK中的序列号返回客户端，正常客户端发报文时，服务器根据报文中携带的cookie重新恢复连接
  - 由于cookie占用序列号空间，导致此时所有TCP可选功能失效，例如扩充窗口、时间戳等

```
net.ipv4.tcp_syncookies = 1
```

## 文件句柄数

linux系统一切皆文件，并发的时候文件句柄数有上限

### 系统全局

- 操作系统可使用的最大句柄数

```
fs.file-max
```

查看命令：sysctl -a|grep file-max

- 使用fs.file-nr可以查看当前已分配、正在使用、上限

```
fs.file-nr=21632 0 40000500
```

查看命令：sysctl -a|grep file-nr

### 用户级别

```
# /etc/security/limits.conf
root soft nofile 65535
root hard nofile 65535
```

### 进程级别

- nginx配置文件

```
worker_rlimit_nofile number; // 默认无
```

## 相关配置

- 客户端 主动建立连接时，发SYN的重试次数

```
net.ipv4.tcp_syn_retries=6
```

- 客户端 建立连接时的本地端口可用范围

```
net.ipv4.ip_local_port_range=32768 60999
```

- 主动建立连接时应用层超时时间

```
proxy_connect_timeout time; // 默认是60s
```

- 服务端 SYN_RCVD状态连接的最大个数(半连接)

```
net.ipv4.tcp_max_syn_backlog
```

- 服务端  被动建立连接时，发SYN/ACK的重试次数

```
net.ipv4.tcp_synack_retries
```

- 接收自网卡，但未被内核协议栈处理的报文队列长度

```
net.core.netdev_max_backlog = 262144
```

- backlog队列长度，nginx配置

```
listen addres[:port][backlog=number]
```

- ACCEPT队列已完成握手

```
net.core.somaxconn
```

- 超出处理能力时，对新来的SYN直接回包RST，丢弃连接

```
net.ipv4.tcp_abort_on_overflow
```

- 当SYN队列满后，新的SYN不进入队列，计算出cookie再以SYN+ACK中的序列号返回客户端，正常客户端发报文时，服务器根据报文中携带的cookie重新恢复连接
  - 由于cookie占用序列号空间，导致此时所有TCP可选功能失效，例如扩充窗口、时间戳等

```
net.ipv4.tcp_syncookies = 1
```

- 设置worker进程最大连接数量

```
worker_connections number; // 默认512
```

- net.ipv4.tcp_fastopen 系统开启TFO功能
  - 0：关闭
  - 1：作为客户端时
  - 2：作为服务端时
  - 3：无论客户端或服务端，都开启TFO

- Tcp Fast Open nginx配置,指定TFO连接队列长度

```
listen addres[:port][fastopen=number];
```

- 两次读操作间的超时

```
client_body_timeout time; // 默认60s
```

- 两次写操作间的超时

```
send_timeout time; // 默认60s
```

- 兼具两次读写操作间的超时

```
proxy_timeout timeout; // 默认10m
```

- 限制重传次数-丢包重传达到上限后，更新路由缓存

```
net.ipv4.tcp_retries1=3
```

- 限制重传次数-丢包重传达到上限后，关闭TCP连接

```
net.ipv4.tcp_retries2=15
```

- Nagle算法在Nginx配置

```
tcp_nodelay off; // 吞吐量优先，开启Nagle算法。默认on
tcp_nodelay on; // 低时延优先，禁用Nagle算法。默认on
```

- Nginx避免发送小报文配置

```
postpone_output size; // 默认1460
```

- 完全禁止小报文的发送

```
tcp_nopush on|off; // 默认off开启
```

- 发送心跳周期

```
net.ipv4.tcp_keepalive_time=7200 // 秒
```

- 探测包发送间隔

```
net.ipv4.tcp_keepalive_intvl=75
```

- 探测包重试次数

```
net.ipv4.tcp_keepalive_probes=9
```

- fin_wait1状态，发送FIN报文的重试次数，0相当于8

```
net.ipv4.tcp_orphan_retries=0
```

- fin_wait2状态，保持在FIN_WAIT_2状态的时间

```
net.ipv4.tcp_fin_timeout=60
```

- 开启后，作为客户端时新连接可以使用仍然处于time_wait状态端口

```
net.ipv4.tcp_tw_reuse=1
```

- 由于timestamp的存在，操作系统可以拒绝迟到的报文

```
net,ipv4.tcp_timestamps=1
```

- 开启后，同时作为客户端和服务器都可以使用time_wait状态的端口，不安全，无法避免报文延迟、重复等给新连接造成混乱

```
net.ipv4.tcp_tw_recycle=0
```

- Time_wait状态连接的最大数量，超出后直接关闭连接

```
net.ipv4.tcp_max_tw_buckets=262144
```

- off 关闭功能
- on由Nginx判断，当用户请求未接收完(根据chunk或者Content-Length头部等)时启用功能，否则及时关闭连接
- always无条件启用功能

```
lingering_close off|on|always; // 默认on
```

- 当功能启用时，最长的读取用户请求内容的时长，达到后立刻关闭连接

```
lingering_time time; //默认30s
```

- 当功能启用时，检测客户端是否仍然请求内容到达，若超时后仍没有数据达到，则立刻关闭连接

```
lingering_timeout time; //默认5s
```

- 当读写超时指令生效引发连接关闭时，通过发送RST立刻释放端口、内存等资源来关闭连接

```
reset_timedout_connection on|off; // 默认off
```


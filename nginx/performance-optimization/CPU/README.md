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

## CPU基本知识

### 为何一个CPU就可以同时运行多个进程？

1. 宏观上并行，微观上串行
   - 把进程的运行时间分为一段段的时间片
   - OS调度系统依次选择每个进程，最多执行时间片指定的时长
2. 阻塞API引发的时间片内主动让出CPU
   - 速度不一致引发的阻塞API：例如CPU和硬盘
   - 业务场景产生的阻塞API：例如同步读网络报文

### 何为进程间切换

是指CPU从一个进程或线程切换到另外一个进程或者线程，类别分为主动切换和被动切换：

- 主动切换：进程或线程IO阻塞
- 被动切换：时间片耗尽

#### 如何查看上下文切换次数

```
// 每隔2秒打印一次
vmstat 2 

cs字段：所有进程上下文切换次数
in字段：所有进程中断次数
```

```
dstat

system.csw 进程上下文切换次数
```

```
// 每隔2秒打印一次
pidstat -w 2
// 指定2655371 pid进程
pidstat -w 1 -p 2655371

cswch/s字段：主动切换
nvcswch/s字段：被动切换
```

### 什么决定CPU时间片的大小？

```
top命令

NI字段：Nice静态优先级(-20 -- 19),number值越小,CPU时间片占用越长
PR字段：Priority动态优先级(0-139),number值越小,优先级高于其他正常进程
```

### 查看进程状态

```
ps -aux
```

STAT字段：

- R运行：正在运行或在运行队列中等待
- S中断：休眠中, 受阻, 在等待某个条件的形成或接受到信号
- D不可中断：收到信号不唤醒和不可运行, 进程必须等待直到有中断发生
- Z僵死：进程已终止, 但进程描述符存在, 直到父进程调用wait4()系统调用后释放
- T停止：进程收到SIGSTOP, SIGSTP, SIGTIN, SIGTOU信号后停止运行运行

### 多核CPU间的负载均衡

1. worker进程间负载均衡

惊群问题：一个网络请求报文进来，激活了所有work进程，然后只有一个work进程处理了这个请求报文，其他进程则进入sleep状态，从而造成了浪费。

由于当时的linux内核处理很有问题，nginx实现了accept_mutex处理惊群问题。早期默认开启accept_mutex on，目前默认关闭 accept_mutex off。linux内核3.9以上采用reuseport处理大大提升了效率。

经测试得出以下数据：

- 吞吐量由小到大：accept_mutex on<accept_mutex off<reuseport

- 时延由长到短：accept_mutex off<accept_mutex on<reuseport
- 处理的波动情况由大到小，越大有些用户体验就很差：accept_mutex on<accept_mutex off<reuseport

可以看出reuseport全面的优于前面的方案，linux3.9以上、centos7以上可以打开reuseport大大提升性能。

reuseport是怎么做到的呢?它是在kernel层面上实际所有的worker都处于listener

			2. 多队列网卡对多核CPU的优化

处理网络报文分为硬中断和软中断

- RSS(Receive Side Scaling)硬中断负载均衡：需要网卡支持，现在大部分网卡都支持。网卡在各个CPU上创建队列可以并发处理硬中断
- RPS(Receive Packet Steering)软中断负载均衡
- RFS(Receive Flow Steering)

不是说把RSS、RPS、RFS都打开，性能就能得到提升。而是在做性能优化的时候，可以考虑这一点，根据场景进行配置优化性能。

			3. 提升CPU缓存命中率：worker_cpu_affinity

每个CPU都有L1、L2、L3的多级缓存,一个程序的运行都是从内存->L3>L2>L1>cpu0。其中L3、L2、L1上有缓存信息，如果一个处理过程因为CPU的切换，导致程序在不同的CPU上执行往往浪费了之前在多级缓存中的信息(速度比没有缓存要慢点)。

查看L1、L2、L3缓存大小

```
~# cat /sys/devices/system/cpu/cpu0/cache/index1/size
32K
~# cat /sys/devices/system/cpu/cpu0/cache/index2/size
4096K
~# cat /sys/devices/system/cpu/cpu0/cache/index3/size
16384K
```

查看多级缓存与哪个CPU核心共享的空间,奇数之间共享，偶数之间共享。下面因为云主机只有2核，所以只有一个

```
~# cat /sys/devices/system/cpu/cpu0/cache/index3/shared_cpu_list
0
~# cat /sys/devices/system/cpu/cpu1/cache/index3/shared_cpu_list
1
```

			4. NUMA架构

随着CPU核数增多16核、32核、内存总线就这么一个，这么多核心并发访问的时候就受限了。NUMA怎把内存和CPU核心分成两部分，比如8核CPU64G的分成2个32G总线分别给CPU0-3、CPU4-7使用。

NUMA节点概念：

- 本地节点:对于某个节点中的所有CPU，此节点称为本地节点。
- 邻居节点:与本地节点相邻的节点称为邻居节点。
- 远端节点:非本地节点或邻居节点的节点，称为远端节点。
- 邻居节点和远端节点,都称作非本地节点(Off Node)。

列举系统上的NUMA节点

```
~# numactl --hardware
available: 1 nodes (0)
node 0 cpus: 0 1
node 0 size: 1987 MB
node 0 free: 208 MB
node distances:
node   0
  0:  10
```

查看NUMA状态

```
~# numastat
                           node0
numa_hit              1097599967
numa_miss                      0
numa_foreign                   0
interleave_hit             21312
local_node            1097599967
other_node                     0
```

numa_hit—命中的，也就是为这个节点成功分配本地内存访问的内存大小
numa_miss—把内存访问分配到另一个node节点的内存大小，这个值和另一个node的numa_foreign相对应。
numa_foreign–另一个Node访问我的内存大小，与对方node的numa_miss相对应
local_node----这个节点的进程成功在这个节点上分配内存访问的大小
other_node----这个节点的进程 在其它节点上分配的内存访问大小
很明显，miss值和foreign值越高，就要考虑绑定的问题。比如配置：禁止访问远程节点

## 如何增大Nginx使用CPU的有效时长？

1. 能够使用全部CPU资源
   - master-worker多进程架构
   - worker进程数量应当大于等于CPU核数
2. Nginx进程间不做无用功浪费CPU资源
   - worker进程不应在繁忙时，主动让出CPU
   - worker进程间不应由于争抢造成资源耗散：worker进程数量应当等于CPU核数
   - worker进程不应调用一些API导致主动让出CPU：拒绝类似的第三方模块
3. 不被其他进程争抢资源
   - 提升进程优先级占用CPU更长时间
   - 减少操作系统上耗资源的非Nginx进程

## 相关配置

- 设置worker进程的数量

```
worker_processes number|auto; // 默认1
```

- 设置worker进程的静态优先级

```
worker_priority number; // 默认0
```

- 延迟处理新连接

```
TCP_DEFER_ACCEPT 延迟处理新连接-TODO
```

- 绑定worker 到指定CPU

```
worker_cpu_affinity auto/cpumask; //默认无绑定
```




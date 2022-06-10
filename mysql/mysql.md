## 1、导出数据

```
mysqldump -h 127.0.0.1 -u 账号 -p密码 --set-gtid-purged=OFF --no-create-db -R demo > demo.sql
```

- 在原数据库导出时加了--set-gtid-purged=OFF，导入新数据库时，才会触发记录到新数据库的binlog日志。如果不加，则新数据库不记录binlog日志。
- --no-create-db 不创建db数据库
- -R 导出存储过程以及自定义函数



### 对正在运行的业务影响最小，不锁表

```
mysqldump -h 127.0.0.1 -u 账号 -p密码 --set-gtid-purged=OFF --master-data --single-transaction -R demo > demo.sql
```

- --master-data 排它锁，需要先开binlog
  - 如果有增删改或事物语句进行到一半还没提交，加了`--master-data`则可以顺利备份。
- --single-transaction

说明start transaction和start transaction with sonsistent snapshot的区别

1. `start transaction`：不是跟执行start transaction命令的时候保存的数据一致，而是跟执行start transaction命令后第一次查询的时候数据保持一致，如果执行start transaction后增加了数据再读，也可以读到增加的数据。但是读一次后再增加数据就看不到了
2. `start transaction with sonsistent snapshot`：是跟执行完这条命令的状态保持一致，后续的增删改都看不到。

加了`--single-transaction`的时候，只会阻塞毫秒级，对对应的表照了个快照从快照里面备份数据，不影响表的增删改。采用的是start transaction with sonsistent snapshot方式。



### 显示耗时

```
time mysqldump -h 127.0.0.1 -u 账号 -p密码 --set-gtid-purged=OFF --no-create-db -R demo > demo.sql
```



## 2、导入数据

### demo 指定导入demo库

```
mysql -h 127.0.0.1 -u 账号 -p密码  demo < demo.sql
```

### 显示耗时

```
time mysql -h 127.0.0.1 -u 账号 -p密码  demo < demo.sql
```



## 3、查看数据库存储使用情况

### 查看指定库所有表情况-demo库

```sql
SELECT table_name AS 'Table Name', CONCAT(ROUND(table_rows/1000000,4),'M') AS 'Number of Rows', CONCAT(ROUND(data_length/(1024*1024*1024),4),'G') AS 'Data Size', CONCAT(ROUND(index_length/(1024*1024*1024),4),'G') AS 'Index Size', CONCAT(ROUND((data_length+index_length)/(1024*1024*1024),4),'G') AS'Total'FROM information_schema.TABLES WHERE table_schema LIKE 'demo';
```



## 4、建表

### 完整复制表

```
// 复制表结构
CREATE TABLE `库`.目标表 LIKE `库`.源表;
// 查入数据
INSERT INTO 目标表 SELECT * FROM 源表;
```



## 5、查看当前session会话

```
show processlist;

// 结束指定会话ID
kill 会话ID
```



## 6、建库

```
CREATE DATABASE `demo` CHARACTER SET utf8;
```



## 7、删除数据

### 清空表

```
TRUNCATE TABLE `demo`;
```



## 8、索引

### 千万级大表在线添加索引

切记业务低谷操作

```
ALTER TABLE tbl_name ADD PRIMARY KEY (column), ALGORITHM=INPLACE, LOCK=NONE;
ALTER TABLE tbl_name ADD INDEX `idx_userid_type_status` (`created_by`,`Type`,`Status`), ALGORITHM=INPLACE, LOCK=NONE;

// 参数说明
ALGORITHM=INPLACE
更优秀的解决方案，在当前表加索引，步骤：
1.创建索引(二级索引)数据字典
2.加共享表锁，禁止DML，允许查询
3.读取聚簇索引，构造新的索引项，排序并插
入新索引
4.等待打开当前表的所有只读事务提交
5.创建索引结束

ALGORITHM=COPY
通过临时表创建索引，需要多一倍存储，还有更多的IO，步骤：
1.新建带索引（主键索引）的临时表
2.锁原表，禁止DML，允许查询
3.将原表数据拷贝到临时表
4.禁止读写,进行rename，升级字典锁
5.完成创建索引操作

LOCK=DEFAULT：默认方式，MySQL自行判断使用哪种LOCK模式，尽量不锁表
LOCK=NONE：无锁：允许Online DDL期间进行并发读写操作。如果Online DDL操
作不支持对表的继续写入，则DDL操作失败，对表修改无效
LOCK=SHARED：共享锁：Online DDL操作期间堵塞写入，不影响读取
LOCK=EXCLUSIVE：排它锁：Online DDL操作期间不允许对锁表进行任何操作
```

## 9、查看当前数据库设置的最大连接数

`show variables like '%max_connection%';`

```sql
mysql> show variables like '%max_connection%';

+-----------------+-------+
| Variable_name   | Value |
+-----------------+-------+
| max_connections | 6000  |
+-----------------+-------+
1 row in set (0.00 sec)

```

## 10、查看当前连接数

`show status like 'Threads%';`

```sql
mysql> show status like 'Threads%';
+-------------------+--------+
| Variable_name     | Value  |
+-------------------+--------+
| Threads_cached    | 33     |
| Threads_connected | 4185   |
| Threads_created   | 795017 |
| Threads_running   | 3      |
+-------------------+--------+
4 rows in set (0.02 sec)

```

## 11、concat拼接kill 用户的会话进程，释放连接数

```sql
mysql>select concat('KILL ',id,';') from information_schema.processlist where user='root' into outfile '/tmp/2022.sql';


mysql>source /tmp/2022.sql;
```

## 12、临时修改变量值-示例wait_timeout

```sql
mysql>show global variables like ‘wait_timeout’;

mysql>set global wait_timeout=120;
```

### 12.1 修改wait_timeout不生效

- 如果使用`show variables like 'wait_timeout';`，查询的是会话变量。然后使用了`set global wait_timeout=120;`修改全局变量，应该使用`show global variables like 'wait_timeout`。这样就可以看到变化了
- interactive_timeout和wait_timeout的值都是86400(24小时)，当两个参数同时出现时值以interactive_timeout为准，所以使用以下操作可以达到修改wait_timeout的效果：

```sql
set global interactive_timeout=120;
exit
再登录查show variables like 'wait_timeout';
```


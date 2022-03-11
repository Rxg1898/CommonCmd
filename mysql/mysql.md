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


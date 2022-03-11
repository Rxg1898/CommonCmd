## 导出数据

```
mysqldump -h 127.0.0.1 -u 账号 -p密码 --set-gtid-purged=OFF --no-create-db -R demo > demo.sql
```

- 在原数据库导出时加了--set-gtid-purged=OFF，导入新数据库时，才会触发记录到新数据库的binlog日志。如果不加，则新数据库不记录binlog日志。
- --no-create-db 不创建db数据库
- -R 导出存储过程以及自定义函数

### 显示耗时

```
time mysqldump -h 127.0.0.1 -u账号 -p密码 --set-gtid-purged=OFF --no-create-db -R demo > demo.sql
```

## 导入数据

### demo 指定导入demo库

```
mysql -h 127.0.0.1 -u 账号 -p密码  demo < demo.sql
```

### 显示耗时

```
time mysql -h 127.0.0.1 -u 账号 -p密码  demo < demo.sql
```

## 查看数据库存储使用情况

### 查看指定库所有表情况-demo库

```sql
SELECT table_name AS 'Table Name', CONCAT(ROUND(table_rows/1000000,4),'M') AS 'Number of Rows', CONCAT(ROUND(data_length/(1024*1024*1024),4),'G') AS 'Data Size', CONCAT(ROUND(index_length/(1024*1024*1024),4),'G') AS 'Index Size', CONCAT(ROUND((data_length+index_length)/(1024*1024*1024),4),'G') AS'Total'FROM information_schema.TABLES WHERE table_schema LIKE 'demo';
```

## 建表

### 完整复制表

```
// 复制表结构
CREATE TABLE `库`.目标表 LIKE `库`.源表;
// 查入数据
INSERT INTO 目标表 SELECT * FROM 源表;
```


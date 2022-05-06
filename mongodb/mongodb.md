**命令行**中密码中有些特殊字符的需要 **\\** 转义，比如：

```
mongosh -u root -p X9\$kI^qsvf --authenticationDatabase admin 
```

## 1、登录

- 登录时,必须明确指定验证库才能登录

- 通常,管理员用的验证库是admin,普通用户的验证库一般是所管理的库设置为验证库
- 如果直接登录到数据库,不进行use,默认的验证库是test,生产不建议使用test库

```
mongosh --host 127.0.0.1 --port 27017 -u root -p mongo123 --authenticationDatabase admin
```

```
mongosh --host 127.0.0.1 --port 27017 -u app -p mongo123 --authenticationDatabase app
```



## 2、创建用户



- 建用户时,use到的库,就是此用户的验证库



role用户角色权限：root(超级管理员)、dbAdmin(库管理员)、readWrite(读写)、read(只读)

db：作用对象(库)

```
use demo
db.createUser(
	{
        user: "app_rw",
        pwd: "123456",
        roles: [ { role: "readWrite", db: "demo" }]
	}
)
```



## 3、查看用户

```
use admin
db.system.users.find().pretty()
```



## 4、删除用户

```
use app
db.dropUser("app02")
```



## 5、导出数据

默认所有数据(json、bson格式文件)，创建一个新目录

```
mkdir /mongodb
mongodump  --port 27017  -o /mongodb
```

-d 指定导出哪个库

-c 指明collection的名字

```
mongodump  --host dds-wz9cc8c0e4afffb41.mongodb.rds.aliyuncs.com --port 27017  -u root -p 123456 --authenticationDatabase admin -d demo -o /mongodb
```



## 6、导入数据

### mongorestore

--drop 表示删除原有数据后导入

```
mongorestore -h 127.0.0.1 -u root -p 123456 --port 27017 --authenticationDatabase admin --drop  /mongodb
```

-d 指定导入哪个库

-c 指明collection的名字

```
mongorestore -u root -p 123456 --host 127.0.0.1 --port 27017 --authenticationDatabase admin -d demo2 /mongodb
```


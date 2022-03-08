## 创建secret

```bash
kubectl create secret docker-registry -n harbor docker-secret --docker-server=registry.cn-shenzhen.aliyuncs.com --docker-username=admin --docker-password=12345
```

- -n 指定harbor命名空间，**该密钥只能在对应namespace使用**
- dockrer-secret: 指定密钥的键名称，**自行定义**
- --docker-server 指定镜像仓库地址
- --docker-username 指定镜像仓库账号
- --docker-password 指定镜像仓库密码

文件如下：

[docker-secret.yaml](./docker-secret.yaml)

## 拉取镜像yaml文件配置此项

```
spec.template.spec.imagePullSecrets

imagePullSecrets:
- name: docker-secret
```




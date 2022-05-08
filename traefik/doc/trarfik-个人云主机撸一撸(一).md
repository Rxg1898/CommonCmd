## tracfik是什么？

[Traefik](https://doc.traefik.io/traefik/) 是一个开源的Edge Router，它使发布你的服务成为一种有趣和简单的体验。它代表你的系统接收请求，并找出负责处理这些请求的组件。
除了它的许多功能外，Traefik的与众不同之处在于它能自动为你的服务发现正确的配置。当Traefik检查你的基础设施时，神奇的事情发生了，它发现了相关信息，并发现哪个服务为哪个请求服务。
Traefik原生兼容每一种主要的集群技术，如Kubernetes、Docker、Docker Swarm、AWS、Mesos、Marathon等；并且可以同时处理许多集群。(它甚至适用于在裸机上运行的传统软件）。
有了Traefik，就不需要维护和同步一个单独的配置文件：一切都会自动、实时地发生（没有重新启动，没有连接中断）。有了Traefik，你可以把时间花在开发和部署新功能上，而不是配置和维护其工作状态上。

## 云主机安装

### 前置准备

- 准备好k8s环境
- 开放相对安全组端口

这里使用的是单机版k8s，安全组策略已经开放了30000-40000端口

### 安装

一共4个[yaml](https://github.com/rxg456/skills-base/tree/main/traefik/k8s-yaml)文件：

- crd.yaml
- rbac.yaml
- deployment.yaml
- dashboard.yaml

#### CRD对象

crd.yaml用于创建CRD对象

```yaml
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: ingressroutes.traefik.containo.us
spec:
  group: traefik.containo.us
  version: v1alpha1
  names:
    kind: IngressRoute
    plural: ingressroutes
    singular: ingressroute
  scope: Namespaced

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: middlewares.traefik.containo.us
spec:
  group: traefik.containo.us
  version: v1alpha1
  names:
    kind: Middleware
    plural: middlewares
    singular: middleware
  scope: Namespaced

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: ingressroutetcps.traefik.containo.us
spec:
  group: traefik.containo.us
  version: v1alpha1
  names:
    kind: IngressRouteTCP
    plural: ingressroutetcps
    singular: ingressroutetcp
  scope: Namespaced

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: ingressrouteudps.traefik.containo.us
spec:
  group: traefik.containo.us
  version: v1alpha1
  names:
    kind: IngressRouteUDP
    plural: ingressrouteudps
    singular: ingressrouteudp
  scope: Namespaced

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: tlsoptions.traefik.containo.us
spec:
  group: traefik.containo.us
  version: v1alpha1
  names:
    kind: TLSOption
    plural: tlsoptions
    singular: tlsoption
  scope: Namespaced

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: tlsstores.traefik.containo.us
spec:
  group: traefik.containo.us
  version: v1alpha1
  names:
    kind: TLSStore
    plural: tlsstores
    singular: tlsstore
  scope: Namespaced

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: traefikservices.traefik.containo.us
spec:
  group: traefik.containo.us
  version: v1alpha1
  names:
    kind: TraefikService
    plural: traefikservices
    singular: traefikservice
  scope: Namespaced
```

#### RBAC权限

rbac.yaml用于给traefik授权k8s集群权限，这里ServiceAccount用户位于kube-system命名空间下

```yaml
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: traefik-ingress-controller
rules:
  - apiGroups:
      - ""
    resources:
      - services
      - endpoints
      - secrets
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - extensions
    resources:
      - ingresses
      - ingressclasses
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - extensions
    resources:
      - ingresses/status
    verbs:
      - update
  - apiGroups:
      - traefik.containo.us
    resources:
      - middlewares
      - ingressroutes
      - traefikservices
      - ingressroutetcps
      - ingressrouteudps
      - tlsoptions
      - tlsstores
    verbs:
      - get
      - list
      - watch

---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: traefik-ingress-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: traefik-ingress-controller
subjects:
  - kind: ServiceAccount
    name: traefik-ingress-controller
    namespace: kube-system

---
kind: ServiceAccount
apiVersion: v1
metadata:
  name: traefik-ingress-controller
  namespace: kube-system
```

#### Deployment控制器

deployment.yaml中的args是traefik的启动参数可按需修改，其中前两项配置是来定义**web**和**websecure**这两个入口点的，**--api=true**开启,就会创建一个名为**api@internal**的特殊 service，在 dashboard 中可以直接使用这个 service 来访问，然后其他比较重要的就是开启 **kubernetesingress** 和 这两个 **kubernetescrd** provider。
这里因为单机云主机所以直接使用了主机网络，且端口也从80、443修改成了30080、30443

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: traefik
  namespace: kube-system
  labels:
    app: traefik
spec:
  selector:
    matchLabels:
      app: traefik
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: traefik
    spec:
      serviceAccountName: traefik-ingress-controller
      tolerations:
      - operator: "Exists"
      containers:
      - image: traefik:2.3
        name: traefik
        ports:
        - name: web
          containerPort: 30080
          hostPort: 30080
        - name: websecure
          containerPort: 30443
          hostPort: 30443
        - name: mysql
          containerPort: 33306
          hostPort: 33306
        args:
        - --entryPoints.web.address=:30080
        - --entryPoints.websecure.address=:30443
        - --entryPoints.mysql.address=:33306        
        - --log.level=INFO
        - --accesslog
        - --api=true
        - --api.dashboard=true
        - --ping=true
        - --providers.kubernetesingress
        - --providers.kubernetescrd
        resources:
          requests:
            cpu: "50m"
            memory: "50Mi"
          limits:
            cpu: "200m"
            memory: "100Mi"
        securityContext:
          allowPrivilegeEscalation: true
          capabilities:
            drop:
            - ALL
            add:
            - NET_BIND_SERVICE
        readinessProbe:
          httpGet:
            path: /ping
            port: 8080
          failureThreshold: 1
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 2
        livenessProbe:
          httpGet:
            path: /ping
            port: 8080
          failureThreshold: 3
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 2
      # 使用主机网络
      hostNetwork: true
```

#### Dashboard

dashboard.yaml为管理页面，其中的dashboard.gitee.com修改对应自己域名，记得在**/etc/hosts**加上对应解析

```yaml
---
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: traefik-dashboard
  namespace: kube-system
spec:
  entryPoints:
  - web
  routes:
  - match: Host(`dashboard.gitee.com`)
    kind: Rule
    services:
    - name: api@internal
      kind: TraefikService
```

启动顺序

```yaml
kubectl apply -f crd.yaml
kubectl apply -f rbac.yaml 
kubectl apply -f deployment.yaml
kubectl apply -f dashboard.yaml

// 显示
root@i-t66xixhz:/opt/traefik# kubectl apply -f crd.yaml 
customresourcedefinition.apiextensions.k8s.io/ingressroutes.traefik.containo.us created
customresourcedefinition.apiextensions.k8s.io/middlewares.traefik.containo.us created
customresourcedefinition.apiextensions.k8s.io/ingressroutetcps.traefik.containo.us created
customresourcedefinition.apiextensions.k8s.io/ingressrouteudps.traefik.containo.us created
customresourcedefinition.apiextensions.k8s.io/tlsoptions.traefik.containo.us created
customresourcedefinition.apiextensions.k8s.io/tlsstores.traefik.containo.us created
customresourcedefinition.apiextensions.k8s.io/traefikservices.traefik.containo.us created
root@i-t66xixhz:/opt/traefik# kubectl apply -f rbac.yaml 
clusterrole.rbac.authorization.k8s.io/traefik-ingress-controller created
clusterrolebinding.rbac.authorization.k8s.io/traefik-ingress-controller created
serviceaccount/traefik-ingress-controller created
root@i-t66xixhz:/opt/traefik# kubectl apply -f deployment.yaml 
deployment.apps/traefik created
root@i-t66xixhz:/opt/traefik# kubectl apply -f dashboard.yaml 
ingressroute.traefik.containo.us/traefik-dashboard created
```

### 访问

[http://dashboard.gitee.com:30080/dashboard/](http://dashboard.gitee.com:30080/dashboard/) 
![2022-05-05 19-54-59屏幕截图.png](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/2022-05-05%2019-54-59%E5%B1%8F%E5%B9%95%E6%88%AA%E5%9B%BE.png)
![2022-05-05 19-57-47屏幕截图.png](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/2022-05-05%2019-57-47%E5%B1%8F%E5%B9%95%E6%88%AA%E5%9B%BE.png)

## 两个版本web

准备了两个Nginx web，通过configmap挂载了不同index内容

```yaml
---
apiVersion: v1
kind: Namespace
metadata:
  name: web

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: web-config
  namespace: web
data:
  index-v1.html: |
    <h1>web v1</h1>
  index-v2.html: |
    <h1>web v2</h1>

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-v1
  namespace: web
spec:
  selector:
    matchLabels:
      app: nginx
      version: v1
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx
        version: v1
    spec:
      containers:
      - name: nginx
        image: nginx:alpine
        ports:
        - containerPort: 80
          name: port-v1
        volumeMounts:
        - name: config
          mountPath: "/usr/share/nginx/html/index.html"
          subPath: index-v1.html
          readOnly: true
      volumes:
        - name: config
          configMap:
            name: web-config

---
apiVersion: v1
kind: Service
metadata:
  name: app-v1
  namespace: web
spec:
  selector:
    version: v1
  ports:
  - name: http
    port: 80
    targetPort: port-v1

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-v2
  namespace: web
spec:
  selector:
    matchLabels:
      app: nginx
      version: v2
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx
        version: v2
    spec:
      containers:
      - name: nginx
        image: nginx:alpine
        ports:
        - containerPort: 80
          name: port-v2
        volumeMounts:
        - name: config
          mountPath: "/usr/share/nginx/html/index.html"
          subPath: index-v2.html
          readOnly: true
      volumes:
        - name: config
          configMap:
            name: web-config

---
apiVersion: v1
kind: Service
metadata:
  name: app-v2
  namespace: web
spec:
  selector:
    version: v2
  ports:
  - name: http
    port: 80
    targetPort: port-v2
```

## 通过traefik暴露服务

创建一个 IngressRoute 资源对象：(web-ingressroute.yaml),配置的 **Service** 是**Kubernetes  Service**对象。

```yaml
---
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: app-ingressroute
  namespace: web
spec:
  entryPoints:
    - web
  routes:
  - match: Host(`web.gitee.com`)
    kind: Rule
    services:
    - name: app-v1
      namespace: web # 命名空间
      port: 80
```

这个时候我们对域名**web.gitee.com**做**/etc/hosts**解析
访问[http://web.gitee.com:30080/](http://web.gitee.com:30080/)

## HTTPS的暴露

基于上面的基础创建证书对应的secret，已申请到的证书文件：tls.crt和tls.key

```yaml
kubectl create secret tls web-tls --cert=tls.crt --key=tls.key
```

IngressRoute配置tls，使用websecure

```yaml
---
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: app-ingressroute
  namespace: web
spec:
  entryPoints:
    - websecure
  routes:
  - match: Host(`web.gitee.com`)
    kind: Rule
    services:
    - name: app-v1
      namespace: web # 命名空间
      port: 80
  tls:
    secretName: web-tls
```

## 中间件

中间件是 Traefik2.0 中一个非常有特色的功能，我们可以根据自己的各种需求去选择不同的中间件来满足服务，Traefik 官方已经内置了许多不同功能的中间件，其中一些可以修改请求，头信息，一些负责重定向，一些添加身份验证等等，而且中间件还可以通过链式组合(洋葱模式)的方式来适用各种情况。
![overview.png](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/overview.png)

### 示例：http跳转https(redirectScheme)

web-middleware.yaml

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: redirect-https
  namespace: web
spec:
  redirectScheme:
    scheme: https
```

然后将这个中间件附加到 http 服务的**IngressRoute**上面去，因为 https 的不需要跳转

```yaml
---
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: app-ingressroute
  namespace: web
spec:
  entryPoints:
    - web
  routes:
  - match: Host(`web.gitee.com`)
    kind: Rule
    services:
    - name: app-v1
      namespace: web # 命名空间
      port: 80
    middlewares: 
    - name: redirect-https
```

这个时候我们再去访问 http 服务可以发现就会自动跳转到 https 去了。更多中间件的用法可以看：[http中间件](https://doc.traefik.io/traefik/middlewares/http/overview/)和[tcp中间件](https://doc.traefik.io/traefik/middlewares/tcp/overview/)

## 灰度发布

Traefik2.0 的一个更强大的功能就是灰度发布，灰度发布我们有时候也会称为金丝雀发布（Canary），主要就是让一部分测试的服务也参与到线上去，经过测试观察看是否符号上线要求。
2.3版本新增了一个 TraefikService 的 CRD 资源，我们可以直接利用这个对象来配置 web app，之前的版本需要通过 File Provider，比较麻烦。
下面利用 Traefik2.0 中提供的**带权重的轮询（WRR）**功能来控制我们的流量，将3/4的流量路由到 app-v1，1/4 的流量路由到 app-v2 。新建一个描述app的资源清单：(web-traefikservice.yaml)

```yaml
---
apiVersion: traefik.containo.us/v1alpha1
kind: TraefikService
metadata:
  name: app
  namespace: web
spec:
  weighted:
    services:
      - name: app-v1
        weight: 3  # 定义权重
        port: 80
        kind: Service  # 可选，默认就是 Service
      - name: app-v2
        weight: 1
        port: 80
```

接着创建一个 IngressRoute 资源对象：(web-ingressroute.yaml),不过需要注意的是现在我们配置的 Service 不再是直接的 Kubernetes 对象了，而是上面我们定义的 TraefikService 对象，直接创建上面的资源对象

```yaml
---
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: app-ingressroute
  namespace: web
spec:
  entryPoints:
    - web
  routes:
  - match: Host(`web.gitee.com`)
    kind: Rule
    services:
    - name: app
      namespace: web # 命名空间
      kind: TraefikService
```

访问[http://web.gitee.com:30080/](http://web.gitee.com:30080/)，去浏览器中连续访问 4 次，我们可以观察到 app-v1 这应用会收到 3 次请求，而 app-v2 这个应用只收到 1 次请求，符合上面我们的权重配置。

## 流量复制

Traefik 2.0 还引入了流量镜像服务，是一种可以将流入流量复制并同时将其发送给其他服务的方法，镜像服务可以获得给定百分比的请求同时也会忽略这部分请求的响应。
在 2.3 版本中我们已经可以通过**TraefikService**资源对象中的**mirroring**来进行配置，下面将服务 v1 的流量复制 50% 到服务 v2

```yaml
---
apiVersion: traefik.containo.us/v1alpha1
kind: TraefikService
metadata:
  name: app
  namespace: web
spec:
  mirroring:
    name: app-v1 # 发送 100% 的请求到 K8S 的 Service "app-v1"
    port: 80
    mirrors:
    - name: app-v2 # 然后复制 50% 的请求到 app-v2 
      percent: 50
      port: 80


---
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: app-ingressroute
  namespace: web
spec:
  entryPoints:
    - web
  routes:
  - match: Host(`web.gitee.com`)
    kind: Rule
    services:
    - name: app
      namespace: web # 命名空间
      kind: TraefikService
```

访问**2**遍浏览器，通过查看日志如下：

```
v1: kubectl logs -f nginx-v1-64c95b8dd6-5sltw -n web 
```


![2022-05-06 14-00-13屏幕截图.png](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/2022-05-06%2014-00-13%E5%B1%8F%E5%B9%95%E6%88%AA%E5%9B%BE.png)

```
v2: kubectl logs -f nginx-v2-57588c6859-qwg5s -n web
```

![2022-05-06 14-00-50屏幕截图.png](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/2022-05-06%2014-00-50%E5%B1%8F%E5%B9%95%E6%88%AA%E5%9B%BE.png)

## TCP服务

首先创建一个普通的mysql服务(mysql-tcp.yaml)

```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: mysql-config
data:
  my.cnf: |
    [mysqld]
    pid-file        = /var/run/mysqld/mysqld.pid
    socket          = /var/run/mysqld/mysqld.sock
    datadir         = /var/lib/mysql
    symbolic-links=0


---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
spec:
  selector:
    matchLabels:
      app: mysql
  strategy:
    type: Recreate # 在创建新 Pods 之前，所有现有的 Pods 会被杀死
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - image: mysql:5.7
        name: mysql
        env:
        - name: MYSQL_ROOT_PASSWORD
          value: a123456
        ports:
        - containerPort: 3306
          name: mysql
        volumeMounts:
        - name: mysqlcnf
          mountPath: /etc/mysql/my.cnf
          subPath: my.cnf
      volumes:
      - name: mysqlcnf
        configMap:
          name: mysql-config

---
apiVersion: v1
kind: Service
metadata:
  name: mysql-traefik
spec:
  selector:
    app: mysql
  ports:
  - port: 3306
```

创建成功后就可以来为 mysql 服务配置一个路由了。我们这里创建一个**IngressRouteTCP** 类型的 CRD 对象（前面我们就已经安装了对应的 CRD 资源），因为没有配置证书，所以HostSNI使用通配符 ***** 进行配置(mysql-ingressroute-tcp.yaml)

```yaml
---
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRouteTCP
metadata:
  name: mysql-traefik-tcp
spec:
  entryPoints:
    - mysql
  routes:
  - match: HostSNI(`*`)
    services:
    - name: mysql-traefik
      port: 3306
```

要注意上面的** entryPoints **部分，是根据我们启动的 Traefik 的静态配置中的 entryPoints 来决定的，我们当然可以使用前面我们定义得 30080(web) 和 30443(websecure) 这两个入口点，但是也可以可以自己添加一个用于 mysql 服务的专门入口点。[更多EntryPoints知识](https://www.qikqiak.com/traefik-book/routing/entrypoints/)

```yaml
- --entryPoints.web.address=:30080
- --entryPoints.websecure.address=:30443
- --entryPoints.mysql.address=:33306  
```

创建成功

```yaml
# kubectl get IngressRouteTCP 
NAME                AGE
mysql-traefik-tcp   9m22s
```

尝试连接

```yaml
# mysql -h 1xx.1xx.1xx.172 -P 33306 -u root -pa123456
mysql: [Warning] Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 2
Server version: 5.7.36 MySQL Community Server (GPL)

Copyright (c) 2000, 2017, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| sys                |
+--------------------+
4 rows in set (0.02 sec)
```

### 安全的TCP(TLS)

创建TLS的TCP需要根据证书文件cert.pem和key.key，创建secret

```yaml
kubectl create secret tls traefik-mysql-certs --cert=cert.pem --key=key.key
```

然后在 **IngressRouteTCP** 对象，增加 TLS 配置：

```yaml
---
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRouteTCP
metadata:
  name: mysql-traefik-tcp
spec:
  entryPoints:
    - mysql
  routes:
  - match: HostSNI(`mysql.traefik.com`)
    services:
    - name: mysql-traefik
      port: 3306
  tls: 
    secretName: traefik-mysql-certs
```

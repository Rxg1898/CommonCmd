## 1.概述
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/middleware.png#crop=0&crop=0&crop=1&crop=1&id=UfRgi&originHeight=1730&originWidth=3165&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)
附加到路由器的中间件是一种在请求发送到您的[服务](https://doc.traefik.io/traefik/routing/services/)之前（或在服务的答案发送到客户端之前）调整请求的方法。
Traefik 中有几个可用的中间件，有的可以修改请求、headers，有的负责重定向，有的添加认证等等。			
使用相同协议的中间件可以组合成链以适应各种场景。

### 1.1可用的中间件
| **中间件** | **目的** | **区域** |
| --- | --- | --- |
| AddPrefix | 添加一个路径前缀 | 路径修改器 |
| BasicAuth | 添加基本认证 | 安全、认证 |
| Buffering | 缓冲请求/应答 | 请求生命周期 |
| Chain | 结合多个中间件 | 其他的 |
| CircuitBreaker | 防止调用不健康的服务 | 请求生命周期 |
| Compress | 压缩响应 | 内容修饰符 |
| ContentType | 处理 Content-Type 内容类型的自动检测 | 其他的 |
| DigestAuth | 增加摘要认证 | 安全、认证 |
| Errors | 自定义错误页面 | 请求生命周期 |
| ForwardAuth | 代表身份验证 | 安全、认证 |
| Headers | 添加/更新Headers头 | 安全 |
| IPWhiteList | 限制允许的IP | 安全、请求生命周期 |
| InFlightReq | 限制同时连接的数量 | 安全、请求生命周期 |
| PassTLSClientCert | 在Headers头中添加Client证书 | 安全 |
| RateLimit | 限制调用频率 | 安全、请求生命周期 |
| RedirectScheme | 基于方案的重定向 | 请求生命周期 |
| RedirectRegex | 基于正则表达式的重定向 | 请求生命周期 |
| ReplacePath | 改变请求路径 | 路径修改器 |
| ReplacePathRegex | 正则改变请求路径 | 路径修改器 |
| Retry | 发生错误时自动重试 | 请求生命周期 |
| StripPrefix | 改变请求路径 | 路径修改器 |
| StripPrefixRegex | 正则改变请求路径 | 路径修改器 |

## 2.Add Prefix
路径前缀
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/AddPrefix.png#crop=0&crop=0&crop=1&crop=1&id=gaDTb&originHeight=239&originWidth=2466&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)
AddPrefix中间件在转发请求之前更新请求的路径
### 2.1配置示例
```yaml
# Prefixing with /foo
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: add-foo
spec:
  addPrefix:
    prefix: /foo
```
### 2.2配置选项
#### 2.2.1prefix
prefix是要在请求的 URL 中的当前路径之前添加的字符串。它应该包括一个前导斜杠 ( /)

## 3.BasicAuth
添加基本身份验证
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/BasicAuth.png#crop=0&crop=0&crop=1&crop=1&id=HeCrK&originHeight=554&originWidth=2445&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)
BasicAuth 中间件将您的服务的访问权限限制为已知用户
### 3.1配置示例
```yaml
# Declaring the user list
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  basicAuth:
    secret: secretName
```
### 3.2配置选项
#### 3.2.1General
密码必须使用 MD5、SHA1 或 BCrypt 进行哈希处理(用于htpasswd生成密码)
#### 3.2.2users
该**users**选项是一组授权用户。必须使用该**name:hashed-password**格式声明每个用户。
```yaml
如果同时提供users和usersFile，则合并这两个文件。usersFile的内容优先于users中的值。

出于安全原因，Kubernetes IngressRoute中用户字段不存在，应该使用secret字段。
```
```yaml
kubernetes支持特殊的kubernetes.io/basic-auth秘密类型。此秘密必须包含两个键：用户名和密码。请注意，这些键不会以任何方式进行哈希或加密，因此比其他方法更安全。您可以找到有关Kubernetes基本身份验证秘密文档的更多信息
```
```yaml
# 声明用户列表
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  basicAuth:
    secret: authsecret

---
# 注意:在kubernetes secret中，字符串(例如由htpasswd生成的)必须首先是base64编码的.
# 要创建一个编码的用户密码对, 可以使用以下命令:
# htpasswd -nb user password | openssl base64

apiVersion: v1
kind: Secret
metadata:
  name: authsecret
  namespace: default
data:
  users: |2
    dGVzdDokYXByMSRINnVza2trVyRJZ1hMUDZld1RyU3VCa1RycUU4d2ovCnRlc3QyOiRhcHIxJGQ5
    aHI5SEJCJDRIeHdnVWlyM0hQNEVzZ2dQL1FObzAK

---
# 这是另一种身份验证秘密，它演示了基本身份验证秘密类型.
# 注意:密码不是哈希的，只是base64编码.

apiVersion: v1
kind: Secret
metadata:
  name: authsecret2
  namespace: default
type: kubernetes.io/basic-auth
data:
  username: dXNlcg== # username: user
  password: cGFzc3dvcmQ= # password: password
```
#### 3.2.3usersFile
usersFile选项是指向外部文件的路径，该文件包含中间件的授权用户
文件内容是一个名称列表**name:hashed-password,**例如以下文件内容：
```yaml
test:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/
test2:$apr1$d9hr9HBB$4HxwgUir3HP4EsggP/QNo0
```

```yaml
因为引用Kubernetes上的文件路径没有多大意义，所以Kubernetes IngressRoute的usersFile字段不存在，应该使用secret字段。
```

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  basicAuth:
    secret: authsecret

---
apiVersion: v1
kind: Secret
metadata:
  name: authsecret
  namespace: default

data:
  users: |2
    dGVzdDokYXByMSRINnVza2trVyRJZ1hMUDZld1RyU3VCa1RycUU4d2ovCnRlc3QyOiRhcHIxJGQ5
    aHI5SEJCJDRIeHdnVWlyM0hQNEVzZ2dQL1FObzAK
```

#### 3.2.4realm
你可以使用领域选项自定义用于身份验证的领域.缺省值为**traefik**.
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  basicAuth:
    realm: MyRealm
```
#### 
#### 3.2.5headerField
你可以使用**headerField**定义一个报头字段来存储经过身份验证的用户.
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: my-auth
spec:
  basicAuth:
    # ...
    headerField: X-WebAuth-User
```
#### 3.2.6removeHeader
设置 **removeHeader** 选项为 **true**，可以将请求转发到你的服务之前删除授权得 Header.（默认值为 false）
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  basicAuth:
    removeHeader: true
```
## 4.Buffering
如何在转发前阅读请求
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/image.png)
缓冲中间件限制了可以转发给服务的请求的大小.
通过缓冲，Traefik将整个请求读入内存（可能将大型请求缓冲到磁盘），并拒绝超过指定大小限制的请求.
这可以帮助服务避免大量的数据（例如multipart/form-data)，并且可以最大程度地减少将数据发送到服务所花费的时间。

### 4.1配置示例
```yaml
# 将最大请求体设置为2Mb
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
name: limit
spec:
buffering:
    maxRequestBodyBytes: 2000000
```
### 4.2配置选项
#### 4.2.1maxRequestBodyBytes
使用 maxRequestBodyBytes 选项， 可以配置允许的请求的最大请求体（以字节为单位）.
如果请求超过允许的大小，则不会将其转发到服务，并且客户端会收到 413 (Request Entity Too Large) 响应.
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: limit
spec:
  buffering:
    maxRequestBodyBytes: 2000000
```
#### 4.2.2memRequestBodyBytes
可以使用 memRequestBodyBytes 选项配置一个阈值（以字节为单位），低于该阈值的请求会缓存到内存，而超过该阈值的请求将缓存到磁盘上.
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: limit
spec:
  buffering:
    memRequestBodyBytes: 2000000
```
#### 4.2.3maxResponseBodyBytes
使用 maxReesponseBodyBytes 选项，可以配置服务允许的最大响应大小（以字节为单位）。
如果响应超出允许的大小，则不会转发给客户端。 客户端改为收到 413 (Request Entity Too Large) 响应.
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: limit
spec:
  buffering:
    maxResponseBodyBytes: 2000000
```
#### 4.2.4memResponseBodyBytes
可以使用 memResponseBodyBytes 选项配置一个阈值（以字节为单位），低于该阈值的响应会缓存到内存，而超过该阈值的响应将缓存到磁盘上.
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: limit
spec:
  buffering:
    memResponseBodyBytes: 2000000
```
#### 4.2.5retryExpression
可以在 retryExpression 选项的帮助下使缓存中间件重试该请求.
```yaml
# 发生网络错误时重试一次
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: limit
spec:
  buffering:
    retryExpression: "IsNetworkError() && Attempts() < 2"
```
重试表达式定义为以下函数与运算符 AND（&&）和 OR（||）的逻辑组合。 至少需要一个函数：

- Attempts() 尝试次数（第一个计数）
- ResponseCode() 服务的响应码
- IsNetworkError() - 如果响应码与网络错误有关

## 5.Chain
当一个中间件还不够时
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/image%20%281%29.png)
链式中间件使您能够定义其它中间件的可重用组合, 这使得重用相同的组更加容易.

### 5.1配置示例
例如，由 WhiteList、BasicAuth 和 HTTPS 组成链式中间件
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: test
  namespace: default

spec:
  entryPoints:
    - web

  routes:
    - match: Host(`mydomain`)
      kind: Rule
      services:
        - name: whoami
          port: 80
      middlewares:
        - name: secured
---
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: secured
spec:
  chain:
    middlewares:
    - name: https-only
    - name: known-ips
    - name: auth-users
---
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: auth-users
spec:
  basicAuth:
    users:
    - test:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/
---
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: https-only
spec:
  redirectScheme:
    scheme: https
---
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: known-ips
spec:
  ipWhiteList:
    sourceRange:
    - 192.168.1.7
    - 127.0.0.1/32
```

## 6.CircuitBreaker
别浪费时间调用不健康得服务
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/image%20%282%29.png)
断路器可以保护你的系统免于将请求堆积到不正常的服务（导致级联故障）上.
当系统运行状况良好时，电路处于关闭状态（正常运行）.当系统运行不正常的时候，电路将断开，并且不再转发请求（而是由后备机制进行处理)
为了评估系统的健康状态，断路器会不断监测你的服务

```yaml
断路器只分析它在中间件链中的位置之后发生的事情。之前发生的事情对其状态没有影响
```
```yaml
每个路由器都有自己的断路器实例。一个断路器实例可以打开，而另一个保持关闭：它们的状态是不共享的.

这是预期的行为，我们希望你能够定义什么是健康的服务，而不需要为每个路由声明一个断路器.
```
### 6.1配置示例
```yaml
# 延迟检查
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: latency-check
spec:
  circuitBreaker:
    expression: LatencyAtQuantileMS(50.0) > 100
```
#### 6.1.1断路器有三种可能的状态

- Closed：服务正常运行
- Open：后备机制接管你的服务
- Recovering：断路器试图通过逐步向你的服务发送请求来恢复正常运行
#### 6.1.2Closed
当电路关闭时，断路器只收集指标来分析请求的行为.
在指定的时间间隔（checkPeriod），断路器评估表达式以决定其状态是否必须改变.
#### 6.1.3Open
当开放时，回退机制在FallbackDuration的持续时间内接管正常的服务调用。在这个持续时间之后，它进入恢复状态
#### 6.1.4Recovering
在恢复期间，断路器向你的服务发送线性增加的请求量（为RecoveryDuration）。如果你的服务在恢复期间出现故障，断路器会再次打开。如果服务在整个恢复期间正常运行，那么断路器就会关闭.
### 6.2配置选项
#### 6.2.1Configuring the Trigger 配置触发器
你可以指定一个表达式，一旦匹配，就会打开断路器并应用回退机制而不是调用你的服务.
表达式选项可以检查三个不同的指标:

- 网络错误率（NetworkErrorRatio）
- 状态代码比率（ResponseCodeRatio）
- 以毫秒为单位的四分位数的延迟（LatencyAtQuantileMS）
##### NetworkErrorRatio
如果你想让断路器在30%的网络错误率下触发，表达式将是
```yaml
NetworkErrorRatio() > 0.30
```
##### ResponseCodeRatio
你可以根据给定的状态代码范围的比率来触发断路器,ResponseCodeRatio接受四个参数:

- from
- to
- dividedByFrom
- dividedByTo

将被计算的操作是sum(to -> from) / sum (dividedByFrom -> dividedByTo)
```yaml
如果总和（dividedByFrom -> dividedByTo）等于0，那么ResponseCodeRatio返回0

from是包容的，to是排斥的
```
例如，表达式**ResponseCodeRatio(500, 600, 0, 600)>0.25**，如果25%的请求返回5XX状态（在返回0到5XX状态代码的请求中），将触发断路器.
##### LatencyAtQuantileMS
你可以在给定比例的请求变得太慢时触发断路器。例如，表达式**LatencyAtQuantileMS(50.0)>100**将在中位延迟（50分位数）达到100MS时触发断路器。
你必须为四分位值提供一个浮点数（尾数为.0）

##### Using multiple metrics
你可以在表达式中使用运算符组合多个指标，支持的运算符是：

- AND (&&)
- OR (||)

例如：**ResponseCodeRatio(500, 600, 0, 600) > 0.30 || NetworkErrorRatio() > 0.10** 当30%的请求返回5XX状态代码，或者当网络错误比例达到10%时，就会触发断路器

##### Operators
下面是支持的运算符的列表：

- 大于 (**>**)
- 大于或等于 (**>=**)
- 小于 (**<**)
- 小于或等于 (**<=**)
- 等于 (**==**)
- 不等于 (**!=**)
#### 6.2.2Fallback mechanism
回退机制向客户端返回 HTTP 503 Service Unavailable（而不是调用目标服务）。这种行为不能被配置

#### 6.2.3CheckPeriod
用于评估表达式并决定断路器的状态是否必须改变的时间间隔。默认情况下，CheckPeriod是100ms。这个值不能被配置

#### 6.2.4FallbackDuration
**回撤时间**默认情况下，**FallbackDuration**是10秒。这个值不能被配置.

#### 6.2.5RecoveringDuration
恢复模式（恢复状态）的持续时间默认情况下，RecoveringDuration是10秒。这个值不能被配置.

## 7.Compress
将响应压缩之后在发送给客户端
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/compress.png)
Compress 中间件启用 gzip 压缩

### 7.1配置示例
```yaml
# 开启 gzip 压缩
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-compress
spec:
  compress: {}
```
响应在以下情况下会被压缩：

- 响应体大于 1400 字节
- Accept-Encoding 请求头包含 gzip
- 响应尚未压缩，即尚未设置 Content-Encoding 响应头
### 7.2配置选项
#### 7.2.1excludedContentTypes
**excludedContentTypes** 指定一系列内容类型，以便在压缩之前将传入请求的 Content-Type 请求头与之对比.
请求中 **excludedContentTypes** 定义的内容类型不会被压缩.
内容类型压缩时忽略大小写和空格.
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-compress
spec:
  compress:
    excludedContentTypes:
      - text/event-stream
```
## 8.DigestAuth
增加 Digest 身份认证
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/digestauth.png)
DigestAuth 中间件是一种将访问权限限制到已知用户的快速方法

### 8.1配置示例
```yaml
# 声明用户列表
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  digestAuth:
    secret: userssecret

# Kubernetes 中使用 secret 来进行认证，可以用下面得命令来生成：
# 用户 admin123 在 qq.com 的认证文件
# htdigets -c auth qq.com admin123 
# kubectl create secret generic userssecret --from-file=auth
```
### 8.2配置选项
**使用 htdigest 生成密码**
#### 8.2.1users
The **users** 是授权用户的数组。每个用户需要使用这种 **name:realm:encoded-password **格式声明。

- 如果同时提供了 users 和 usersFile，则两者将合并。usersFile 的内容优先于 users 中的值。
- 出于安全原因，Kubernetes IngressRoute 的用户字段 users 不存在，而应该使用 secret 字段。
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  digestAuth:
    secret: authsecret

---
apiVersion: v1
kind: Secret
metadata:
  name: authsecret
  namespace: default

data:
  users: |2
    dGVzdDp0cmFlZmlrOmEyNjg4ZTAzMWVkYjRiZTZhMzc5N2YzODgyNjU1YzA1CnRlc3QyOnRyYWVmaWs6NTE4ODQ1ODAwZjllMmJmYjFmMWY3NDBlYzI0ZjA3NGUKCg==
```
#### 8.2.2usersFile
**usersFile** 选项是指向外部文件的路径，该文件包含中间件的授权用户。
文件内容是 **name:realm:encoded-password **格式的列表。
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  digestAuth:
    secret: authsecret

---
apiVersion: v1
kind: Secret
metadata:
  name: authsecret
  namespace: default

data:
  users: |2
    dGVzdDokYXByMSRINnVza2trVyRJZ1hMUDZld1RyU3VCa1RycUU4d2ovCnRlc3QyOiRhcHIxJGQ5
    aHI5SEJCJDRIeHdnVWlyM0hQNEVzZ2dQL1FObzAK
```
一个包含test/test和test2/test2的文件:
```yaml
test:traefik:a2688e031edb4be6a3797f3882655c05
test2:traefik:518845800f9e2bfb1f1f740ec24f074e
```
#### 8.2.3realm
可以使用 **realm** 选项自定义身份验证领域。 默认值为 **traefik**。
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  digestAuth:
    realm: MyRealm
```
#### 8.2.4headerField
可以使用 **headerField** 选项为经过身份验证的用户自定义标题字段。
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: my-auth
spec:
  digestAuth:
    # ...
    headerField: X-WebAuth-User
```
#### 8.2.5removeHeader
将 **removeHeader** 选项设置为 **true** 以在将请求转发到您的服务之前删除授权标头。 （默认值为 **false** ）
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  digestAuth:
    removeHeader: true
```
## 9.ErrorPage
说出错误从来没有这么容易过！
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/errorpages.png)
**ErrorPage**中间件根据HTTP状态码的配置范围返回一个自定义页面来代替默认页面。
**注意：**错误页面本身不是由Traefik托管的

### 9.1配置示例
在例子中，错误页面URL基于状态代码 (query=/{status}.html).
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-errorpage
spec:
  errors:
    status:
      - "500-599"
    query: /{status}.html
    service:
      name: whoami
      port: 80
```
### 9.2配置选项
#### 9.2.1status
status选项定义了哪些状态或状态范围会导致错误页面。状态代码范围是包含的(**500-599**将触发500到599之间的每个代码，**包括**500到599)。
**注意：**可以将状态码定义为一个数字(500)、多个用逗号分隔的数字(**500,502**)、两个代码之间用破折号(**500-599)**分隔的范围或两者的组合(**404,418,500-599**)。
#### 9.2.2service
将提供新请求的错误页面的服务，在Kubernetes中，需要引用Kubernetes服务而不是Traefik服务。
**注意：Host Header** 缺省情况下，客户端**Host**头值被转发到配置的错误服务。要转发与配置的错误服务URL对应的**Host**值，**passHostHeader**选项必须设置为**false**
#### 9.2.3query
错误页面的URL（由**service**托管）。你可以在**query**选项中使用{status}变量，以便在URL中插入状态代码.

## 10.ForwardAuth
使用外部服务来转发认证
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/authforward.png)
ForwardAuth中间件将认证委托给一个外部服务。如果该服务的回答是2XX代码，则允许访问，并执行原始请求。否则，将返回认证服务器的响应。

### 10.1配置示例
```yaml
# 转发认证到example.com
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  forwardAuth:
    address: https://example.com/auth
```
### 10.2Forward-Request Headers
以下请求属性将作为 **X-Forwarded- **头信息提供给forward-auth目标端点。

| 属性 | 转发请求标头 |
| --- | --- |
| HTTP Method | X-Forwarded-Method |
| Protocol | X-Forwarded-Proto |
| Host | X-Forwarded-Host |
| Request URI | X-Forwarded-Uri |
| Source IP-Address | X-Forwarded-For |

### 10.3配置选项
#### 10.3.1address
地址选项定义了认证服务器地址
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  forwardAuth:
    address: https://example.com/auth
```
#### 10.3.2trustForwardHeader
将**trustForwardHeader**选项设置为**true**，以信任所有**X-Forwarded-***头
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  forwardAuth:
    address: https://example.com/auth
    trustForwardHeader: true
```
#### 10.3.3authResponseHeaders
authResponseHeaders选项是要从认证服务器响应中复制并在转发请求中设置的头信息列表，替换任何现有的冲突头信息。
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  forwardAuth:
    address: https://example.com/auth
    authResponseHeaders:
      - X-Auth-User
      - X-Secret
```
#### 10.3.4authResponseHeadersRegex
authResponseHeadersRegex选项是匹配头信息的正则表达式，用于从认证服务器的响应中复制并在转发的请求中设置，在剥离所有匹配正则表达式的头信息后。它允许部分匹配正则表达式与头文件的密钥。应该使用字符串的开头（^）和字符串的结尾（$）锚，以确保与头文件键完全匹配。
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  forwardAuth:
    address: https://example.com/auth
    authResponseHeadersRegex: ^X-
```
**提示：**

- 正则表达式和替换可以使用在线工具，如[Go Playground](https://go.dev/play/)或[Regex101](https://regex101.com/)进行测试。
- 在YAML中定义正则表达式时，任何转义字符都需要转义两次： example\.com需要写成 example\\.com。
#### 10.3.5authRequestHeaders
**authRequestHeaders**选项是要从请求中复制到认证服务器的头文件列表。它允许过滤那些不应该被传递给认证服务器的头信息。如果没有设置或为空，那么所有的请求头信息都会被传递。
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  forwardAuth:
    address: https://example.com/auth
    authRequestHeaders:
      - "Accept"
      - "X-CustomHeader"
```
#### 10.3.6tls(可选)
定义了用于与认证服务器安全连接的TLS配置
##### ca(可选)
ca是用于与认证服务器安全连接的证书颁发机构的路径，它默认为系统捆绑。
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  forwardAuth:
    address: https://example.com/auth
    tls:
      caSecret: mycasercret

---
apiVersion: v1
kind: Secret
metadata:
  name: mycasercret
  namespace: default

data:
  # 必须包含 "tls.ca "或 "ca.crt "密钥下的证书. 
  tls.ca: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0=
```
##### caOptional(可选)
**caOptional**的值定义了TLS客户端认证与认证服务器的安全连接应该使用哪个策略。

**注意：**如果**ca**是未定义的，这个选项将被忽略，在握手过程中不会要求客户证书。因此，任何提供的证书将永远不会被验证。

当该选项被设置为 **true** 时，在握手过程中会要求提供客户证书，但不一定需要。如果发送了证书，它必须是有效的。
当该选项设置为 **false **时，在握手过程中会要求提供客户证书，并且客户应该至少发送一份有效的证书。
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  forwardAuth:
    address: https://example.com/auth
    tls:
      caOptional: true
```
##### cert(可选)
**cert**是用于与认证服务器安全连接的公共证书的路径。使用该选项时，需要设置**key**选项。

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  forwardAuth:
    address: https://example.com/auth
    tls:
      certSecret: mytlscert

---
apiVersion: v1
kind: Secret
metadata:
  name: mytlscert
  namespace: default

data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0=
  tls.key: LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0=
```
**注意：**出于安全原因，Kubernetes IngressRoute不存在该字段，应该使用秘密字段来代替。
##### key(可选)
**key**是用于与认证服务器安全连接的私人密钥的路径。使用该选项时，需要设置**cert**选项。
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  forwardAuth:
    address: https://example.com/auth
    tls:
      certSecret: mytlscert

---
apiVersion: v1
kind: Secret
metadata:
  name: mytlscert
  namespace: default

data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0=
  tls.key: LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0=
```
**注意：**出于安全原因，Kubernetes IngressRoute不存在该字段，应该使用秘密字段来代替。
##### insecureSkipVerify(可选，默认false)
如果 **insecureSkipVerify **为 **true **，则与认证服务器的TLS连接接受服务器提交的任何证书，而不考虑其涵盖的主机名。
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-auth
spec:
  forwardAuth:
    address: https://example.com/auth
    tls:
      insecureSkipVerify: true
```


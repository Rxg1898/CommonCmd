## 1.概述

![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/middleware.png#crop=0&crop=0&crop=1&crop=1&id=UfRgi&originHeight=1730&originWidth=3165&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=#crop=0&crop=0&crop=1&crop=1&id=jYsPu&originHeight=1730&originWidth=3165&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)
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


## 2.Headers

管理请求/响应报头

![](https://doc.traefik.io/traefik/assets/img/middleware/headers.png#crop=0&crop=0&crop=1&crop=1&id=OdVm7&originHeight=410&originWidth=2420&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)

报头中间件管理请求和响应的报头。默认情况下，会自动添加一组转发的标题。更多请看 [FAQ](https://doc.traefik.io/traefik/getting-started/faq/#what-are-the-forwarded-headers-when-proxying-http-requests) .

### 2.1配置示例

#### 2.1.1向请求和响应添加header

下面的示例将X-Script-Name头添加到代理请求中，将X-Custom-Response-Header头添加到响应中

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-header
spec:
  headers:
    customRequestHeaders:
      X-Script-Name: "test"
    customResponseHeaders:
      X-Custom-Response-Header: "value"
```

#### 2.1.2添加和删除头信息

在下面的示例中，请求使用一个额外的X-Script-Name头进行代理，同时去掉它们的X-Custom-Request-Header头，而响应则去掉它们的X-Custom-Response-Header头.

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-header
spec:
  headers:
    customRequestHeaders:
      X-Script-Name: "test" # 添加
      X-Custom-Request-Header: "" # 移除
    customResponseHeaders:
      X-Custom-Response-Header: "" # 移除
```

#### 2.1.3使用安全Headers

安全相关的报头(HSTS headers，浏览器XSS过滤器等)可以像上面所示的那样管理自定义报头。通过添加标题，这个功能可以方便地使用安全特性.

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-header
spec:
  headers:
    frameDeny: true
    browserxssfilter: true
```

#### 2.1.4跨域 CORS Headers

CORS (Cross-Origin Resource Sharing) （跨源资源共享）headers，可以以类似于上述自定义标头的方式进行添加和配置。这一功能允许快速设置更高级的安全功能。如果设置了CORS头，那么中间件不会将预检请求传递给任何服务，相反，响应将被生成并直接发回给客户端。

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-header
spec:
  headers:
    accessControlAllowMethods:
      - "GET"
      - "OPTIONS"
      - "PUT"
    accessControlAllowOriginList:
      - "https://foo.bar.org"
      - "https://example.org"
    accessControlMaxAge: 100
    addVaryHeader: true
```

### 2.2配置选项

#### 2.2.1常规

**注意：**

- 如果现有头文件具有相同的名称，则自定义头文件将覆盖它们。
- 关于安全头的详细文档可以在[unrolled/secure](https://github.com/unrolled/secure#available-options).

#### 2.2.2customRequestHeaders

customRequestHeaders选项列出了要应用于请求的报头名称和值

#### 2.2.3customResponseHeaders

customResponseHeaders选项列出要应用于响应的报头名称和值

#### 2.2.4accessControlAllowCredentials

accessControlAllowCredentials指示请求是否可以包含用户凭据

#### 2.2.5accessControlAllowHeaders

accessControlAllowHeaders指出哪些报头字段名可以用作请求的一部分

#### 2.2.6accessControlAllowMethods

accessControlAllowMethods指示在请求期间可以使用哪些方法

#### 2.2.7accessControlAllowOriginList

- accessControlAllowOriginList指示是否可以通过返回不同的值来共享资源
- 还可以配置通配符 *，并匹配所有请求。如果该值是后端服务设置的，则会被Traefik覆盖
- 该值可以包含允许的源的列表
- 更多信息，包括如何使用设置可以找到 
   - [Mozilla.org](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin)
   - [w3](https://fetch.spec.whatwg.org/#http-access-control-allow-origin)
   - [IETF](https://tools.ietf.org/html/rfc6454#section-7.1)
- Traefik不再支持空值，[因为不再建议将其作为返回值](https://w3c.github.io/webappsec-cors-for-developers/#avoid-returning-access-control-allow-origin-null)

#### 2.2.8accessControlAllowOriginListRegex

**accessControlAllowOriginListRegex**选项是与**accessControlAllowOriginList**选项对应的正则表达式选项，而不是原始值。它允许在**accessControlAllowOriginList**中包含正则表达式的任何匹配的所有源.

**提示：**

- 正则表达式和替换可以使用在线工具，如[Go Playground](https://go.dev/play/)或[Regex101](https://regex101.com/)进行测试。
- 在YAML中定义正则表达式时，任何转义字符都需要转义两次： example.com需要写成 example\.com。

#### 2.2.9accessControlExposeHeaders

**accessControlExposeHeaders**指示哪些头是可以安全公开给CORS api规范的api的

#### 2.2.10addVaryHeader

**addVaryHeader**与**accessControlAllowOriginList**一起使用，以确定是否应该添加或修改Vary报头，以证明服务器响应可以根据原始报头的值而不同.

#### 2.2.11allowedHosts

**allowedHosts**选项列出了允许的完全限定的域名

#### 2.2.12hostsProxyHeaders

**hostsProxyHeaders**选项是一组头键，其中可能包含请求的代理主机名值

#### 2.2.13sslRedirect

**提示：**已弃用，转而使用[ntryPoint redirection](https://doc.traefik.io/traefik/routing/entrypoints/#redirection) or the [RedirectScheme middleware](https://doc.traefik.io/traefik/middlewares/http/redirectscheme/)

当设置为true时，sslRedirect只允许HTTPS请求

#### 2.2.14sslTemporaryRedirect

**提示：**已弃用，转而使用[ntryPoint redirection](https://doc.traefik.io/traefik/routing/entrypoints/#redirection) or the [RedirectScheme middleware](https://doc.traefik.io/traefik/middlewares/http/redirectscheme/)

将sslTemporaryRedirect设置为true以强制使用302(而不是301)进行SSL重定向

#### 2.2.15sslHost

**提示：**已弃用，转而使用[RedirectRegex middleware](https://doc.traefik.io/traefik/middlewares/http/redirectregex/).

sslHost选项是用于将HTTP请求重定向到HTTPS的主机名

#### 2.2.16sslProxyHeaders

**sslProxyHeaders**选项是一组标头键，它们带有指示有效HTTPS请求的相关值。当使用其他代理时它可能很有用(例如:**"X-Forwarded-Proto":"https"**)

#### 2.2.17sslForceHost

**提示：**已弃用，转而使用[RedirectRegex middleware](https://doc.traefik.io/traefik/middlewares/http/redirectregex/).

将**sslForceHost**设置为**true**，并设置**sslHost**强制请求使用**sslHost**，而不管它们是否已经使用SSL.

#### 2.2.18stsSeconds

**stsSeconds**是**Strict-Transport-Security**报头的最大age。如果设置为0，则不设置报头.

#### 2.2.19stsIncludeSubdomains

如果**stsincluubdomains**设置为**true**，则**incluubdomains**指令会被追加到**Strict-Transport-Security**报头.

#### 2.2.20stsPreload

设置**stsPreload**为**true**，将**preload**标志附加到**Strict-Transport-Security**报头

#### 2.2.21forceSTSHeader

设置**forceSTSHeader**为**true**，即使连接是HTTP也可以添加STS头

#### 2.2.22frameDeny

设置**frameDeny**为**true**，添加值为**DENY**的**X-Frame-Options**报头

#### 2.2.23customFrameOptionsValue

**customFrameOptionsValue**允许使用自定义值设置**X-Frame-Options**报头值。这将覆盖**FrameDeny**选项

#### 2.2.24contentTypeNosniff

设置**contentTypeNosniff**为true，以添加值为**nosniff**的**X-Content-Type-Options**报头

#### 2.2.25browserXssFilter

将**browserXssFilter**设为**true**，以增加**X-XSS-Protection**头的值1；mode=block

#### 2.2.26customBrowserXSSValue

**customBrowserXssValue**选项允许使用自定义值设置**X-XSS-Protection**头值。这将覆盖**BrowserXssFilter**选项.

#### 2.2.27contentSecurityPolicy

**contentSecurityPolicy**选项允许**Content-Security-Policy**报头值设置为自定义值

#### 2.2.28publicKey

**publicKey**通过HPKP协议防止伪造证书的MITM攻击

#### 2.2.29referrerPolicy

**referrerPolicy**允许站点控制浏览器是否将**Referer**头转发给其他站点

#### 2.2.30featurePolicy

**提示：**已弃用，转而使用**permissionpolicy**

featurePolicy允许站点控制浏览器特性

#### 2.2.31permissionsPolicy

permissionpolicy允许站点控制浏览器特性

#### 2.2.32isDevelopment

在开发时将**isDevelopment**设置为**true**，以减轻**AllowedHosts**、SSL和STS选项带来的不必要的影响。通常，测试是使用HTTP，而不是HTTPS，并且是在**localhost**，而不是你的生产域。如果你想让你的开发环境模仿生产环境，有完整的主机封锁、SSL重定向和STS头文件，请将此设置为**false**.

## 3.IPWhiteList
限制客户端使用特定ip
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/IPWhiteList.png)
**IPwhitelist**接受 /拒绝基于客户端IP的请求

### 3.1配置示例
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-ipwhitelist
spec:
  ipWhiteList:
    sourceRange:
      - 127.0.0.1/32
      - 192.168.1.7
```
### 3.2配置选项
#### 3.2.1 sourceRange
**sourceRange**选项设置允许的ip地址(或使用CIDR表示法允许的ip地址范围).
#### 3.2.2 ipStrategy
**ipStrategy**选项定义了两个参数，用于设置Traefik如何确定客户端IP: **depth**和**excludedip**.
##### ipStrategy.depth
depth选项告诉Traefik使用X-Forwarded-For报头并获取位于depth位置的IP(从右侧开始).

- 如果深度大于X-Forwarded-For中的IP总数，则客户端IP为空
- 如果深度的值小于或等于0，则忽略深度

depth和X-Forwarded-For的例子：如果深度设置为2，请求X-Forwarded-For的头部是“10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1”，那么真正的客户端IP是“10.0.0.1”(深度4)，但是用于白名单的IP是“12.0.0.1”(深度2)

| X-Forwarded-For | depth | `clientIP` |
| --- | --- | --- |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | 1 | "13.0.0.1" |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | 3 | "11.0.0.1" |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | 5 | "" |

```yaml
# 基于depth=2的X-Forwarded-For的白名单
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-ipwhitelist
spec:
  ipWhiteList:
    sourceRange:
      - 127.0.0.1/32
      - 192.168.1.7
    ipStrategy:
      depth: 2
```
##### ipStrategy.excludedIPs
**excludedIPs**配置Traefik扫描**X-Forwarded-For**报头，并选择第一个不在列表中的IP
**提示：**如果指定了depth，则忽略excludedip

**excludedIPs**和**X-Forwarded-For**示例：

| X-Forwarded-For | excludedIPs | clientIP |
| --- | --- | --- |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | "12.0.0.1,13.0.0.1" | "11.0.0.1" |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | "15.0.0.1,13.0.0.1" | "12.0.0.1" |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | "10.0.0.1,13.0.0.1" | "12.0.0.1" |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | "15.0.0.1,16.0.0.1" | “13.0.0.1" |
| "10.0.0.1,11.0.0.1" | "10.0.0.1,11.0.0.1" | "" |

```yaml
# Exclude from `X-Forwarded-For`
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-ipwhitelist
spec:
  ipWhiteList:
    ipStrategy:
      excludedIPs:
        - 127.0.0.1/32
        - 192.168.1.7
```
## 4.InFlightReq
限制同时请求的数量
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/InFlightReq.png)
为了主动防止服务因高负载而不堪重负，可以限制允许同时进行的在线请求的数量。

### 4.1配置示例
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-inflightreq
spec:
  inFlightReq:
    amount: 10
```
### 4.2配置选项
#### 4.2.1 amount
**amount**选项定义允许同时进行的请求的最大数量。如果已经有大量的请求在进行中(基于相同的sourceCriterion策略)，中间件将以**HTTP 429 Too Many Requests**响应.
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-inflightreq
spec:
  inFlightReq:
    amount: 10
```
#### 4.2.2 sourceCriterion
**sourceCriterion**选项定义了用什么标准将请求分组为来自一个共同的来源。如果同时定义了几个策略，将产生一个错误。如果没有设置，默认是使用 **requestHost**
##### sourceCriterion.ipStrategy
**ipStrategy**选项定义了两个参数，用于配置Traefik如何确定客户的IP：**depth**和**excludedIPs**排除的IP
###### IPSTRATEGY.DEPTH
depth选项告诉Traefik使用X-Forwarded-For报头并获取位于depth位置的IP(从右侧开始).

- 如果depth大于X-Forwarded-For中的IP总数，则客户端IP为空
- 如果深度的值小于或等于0，则忽略深度

depth和X-Forwarded-For的例子：如果深度设置为2，请求X-Forwarded-For的头部是“10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1”，那么真正的客户端IP是“10.0.0.1”(深度4)，但是用于白名单的IP是“12.0.0.1”(深度2)

| X-Forwarded-For | depth | `clientIP` |
| --- | --- | --- |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | 1 | "13.0.0.1" |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | 3 | "11.0.0.1" |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | 5 | "" |

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-inflightreq
spec:
  inFlightReq:
    sourceCriterion:
      ipStrategy:
        depth: 2
```
###### IPSTRATEGY.EXCLUDEDIPS
**excludedIPs**配置Traefik扫描**X-Forwarded-For**报头，并选择第一个不在列表中的IP
**提示：**如果指定了depth，则忽略excludedip

**excludedIPs**和**X-Forwarded-For**示例：

| X-Forwarded-For | excludedIPs | clientIP |
| --- | --- | --- |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | "12.0.0.1,13.0.0.1" | "11.0.0.1" |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | "15.0.0.1,13.0.0.1" | "12.0.0.1" |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | "10.0.0.1,13.0.0.1" | "12.0.0.1" |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | "15.0.0.1,16.0.0.1" | “13.0.0.1" |
| "10.0.0.1,11.0.0.1" | "10.0.0.1,11.0.0.1" | "" |

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-inflightreq
spec:
  inFlightReq:
    sourceCriterion:
      ipStrategy:
        excludedIPs:
        - 127.0.0.1/32
        - 192.168.1.7
```
##### sourceCriterion.requestHeaderName
用于对传入请求进行分组的报头的名称
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-inflightreq
spec:
  inFlightReq:
    sourceCriterion:
      requestHeaderName: username
```
##### sourceCriterion.requestHost
是否考虑将请求主机作为源
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-inflightreq
spec:
  inFlightReq:
    sourceCriterion:
      requestHost: true
```
## 5.PassTLSClientCert
在Header中添加客户端证书
PassTLSClientCert将从传递的客户端TLS证书中选择的数据添加到Header
### 5.1配置示例
在X-Forwarded-Tls-Client-Cert头中传递转义的pem
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-passtlsclientcert
spec:
  passTLSClientCert:
    pem: true
```
传递X-Forwarded-Tls-Client-Cert标题中的所有可用信息
```yaml
# 传递X-Forwarded-Tls-Client-Cert标题中的所有可用信息
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-passtlsclientcert
spec:
  passTLSClientCert:
    info:
      notAfter: true
      notBefore: true
      sans: true
      subject:
        country: true
        province: true
        locality: true
        organization: true
        organizationalUnit: true
        commonName: true
        serialNumber: true
        domainComponent: true
      issuer:
        country: true
        province: true
        locality: true
        organization: true
        commonName: true
        serialNumber: true
        domainComponent: true
```
### 5.2配置选项
#### 5.2.1 常规配置
**PassTLSClientCert**可以向请求添加两个头：

- X-Forwarded-Tls-Client-Cert包含转义的pem
- X-Forwarded-Tls-Client-Cert-Info在一个转义字符串中包含所有选定的证书信息

**提示：**

- 每个头的值都是一个字符串，为了成为一个有效的URL查询，已经被转义了.
- 这些选项只对 [MutualTLS configuration](https://doc.traefik.io/traefik/https/tls/#client-authentication-mtls)配置起作用。也就是说，只有符合clientAuth.clientAuthType策略的证书才能被传递.

下面的例子显示了一个完整的证书，并解释了每个中间件的选项:
```yaml
# 完整的客户端TLS证书
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number: 1 (0x1)
        Signature Algorithm: sha1WithRSAEncryption
        Issuer: DC=org, DC=cheese, O=Cheese, O=Cheese 2, OU=Simple Signing Section, OU=Simple Signing Section 2, CN=Simple Signing CA, CN=Simple Signing CA 2, C=FR, C=US, L=TOULOUSE, L=LYON, ST=Signing State, ST=Signing State 2/emailAddress=simple@signing.com/emailAddress=simple2@signing.com
        Validity
            Not Before: Dec  6 11:10:16 2018 GMT
            Not After : Dec  5 11:10:16 2020 GMT
        Subject: DC=org, DC=cheese, O=Cheese, O=Cheese 2, OU=Simple Signing Section, OU=Simple Signing Section 2, CN=*.example.org, CN=*.example.com, C=FR, C=US, L=TOULOUSE, L=LYON, ST=Cheese org state, ST=Cheese com state/emailAddress=cert@example.org/emailAddress=cert@sexample.com
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                RSA Public-Key: (2048 bit)
                Modulus:
                    00:de:77:fa:8d:03:70:30:39:dd:51:1b:cc:60:db:
                    a9:5a:13:b1:af:fe:2c:c6:38:9b:88:0a:0f:8e:d9:
                    1b:a1:1d:af:0d:66:e4:13:5b:bc:5d:36:92:d7:5e:
                    d0:fa:88:29:d3:78:e1:81:de:98:b2:a9:22:3f:bf:
                    8a:af:12:92:63:d4:a9:c3:f2:e4:7e:d2:dc:a2:c5:
                    39:1c:7a:eb:d7:12:70:63:2e:41:47:e0:f0:08:e8:
                    dc:be:09:01:ec:28:09:af:35:d7:79:9c:50:35:d1:
                    6b:e5:87:7b:34:f6:d2:31:65:1d:18:42:69:6c:04:
                    11:83:fe:44:ae:90:92:2d:0b:75:39:57:62:e6:17:
                    2f:47:2b:c7:53:dd:10:2d:c9:e3:06:13:d2:b9:ba:
                    63:2e:3c:7d:83:6b:d6:89:c9:cc:9d:4d:bf:9f:e8:
                    a3:7b:da:c8:99:2b:ba:66:d6:8e:f8:41:41:a0:c9:
                    d0:5e:c8:11:a4:55:4a:93:83:87:63:04:63:41:9c:
                    fb:68:04:67:c2:71:2f:f2:65:1d:02:5d:15:db:2c:
                    d9:04:69:85:c2:7d:0d:ea:3b:ac:85:f8:d4:8f:0f:
                    c5:70:b2:45:e1:ec:b2:54:0b:e9:f7:82:b4:9b:1b:
                    2d:b9:25:d4:ab:ca:8f:5b:44:3e:15:dd:b8:7f:b7:
                    ee:f9
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Key Usage: critical
                Digital Signature, Key Encipherment
            X509v3 Basic Constraints:
                CA:FALSE
            X509v3 Extended Key Usage:
                TLS Web Server Authentication, TLS Web Client Authentication
            X509v3 Subject Key Identifier:
                94:BA:73:78:A2:87:FB:58:28:28:CF:98:3B:C2:45:70:16:6E:29:2F
            X509v3 Authority Key Identifier:
                keyid:1E:52:A2:E8:54:D5:37:EB:D5:A8:1D:E4:C2:04:1D:37:E2:F7:70:03

            X509v3 Subject Alternative Name:
                DNS:*.example.org, DNS:*.example.net, DNS:*.example.com, IP Address:10.0.1.0, IP Address:10.0.1.2, email:test@example.org, email:test@example.net
    Signature Algorithm: sha1WithRSAEncryption
         76:6b:05:b0:0e:34:11:b1:83:99:91:dc:ae:1b:e2:08:15:8b:
         16:b2:9b:27:1c:02:ac:b5:df:1b:d0:d0:75:a4:2b:2c:5c:65:
         ed:99:ab:f7:cd:fe:38:3f:c3:9a:22:31:1b:ac:8c:1c:c2:f9:
         5d:d4:75:7a:2e:72:c7:85:a9:04:af:9f:2a:cc:d3:96:75:f0:
         8e:c7:c6:76:48:ac:45:a4:b9:02:1e:2f:c0:15:c4:07:08:92:
         cb:27:50:67:a1:c8:05:c5:3a:b3:a6:48:be:eb:d5:59:ab:a2:
         1b:95:30:71:13:5b:0a:9a:73:3b:60:cc:10:d0:6a:c7:e5:d7:
         8b:2f:f9:2e:98:f2:ff:81:14:24:09:e3:4b:55:57:09:1a:22:
         74:f1:f6:40:13:31:43:89:71:0a:96:1a:05:82:1f:83:3a:87:
         9b:17:25:ef:5a:55:f2:2d:cd:0d:4d:e4:81:58:b6:e3:8d:09:
         62:9a:0c:bd:e4:e5:5c:f0:95:da:cb:c7:34:2c:34:5f:6d:fc:
         60:7b:12:5b:86:fd:df:21:89:3b:48:08:30:bf:67:ff:8c:e6:
         9b:53:cc:87:36:47:70:40:3b:d9:90:2a:d2:d2:82:c6:9c:f5:
         d1:d8:e0:e6:fd:aa:2f:95:7e:39:ac:fc:4e:d4:ce:65:b3:ec:
         c6:98:8a:31
-----BEGIN CERTIFICATE-----
MIIGWjCCBUKgAwIBAgIBATANBgkqhkiG9w0BAQUFADCCAYQxEzARBgoJkiaJk/Is
ZAEZFgNvcmcxFjAUBgoJkiaJk/IsZAEZFgZjaGVlc2UxDzANBgNVBAoMBkNoZWVz
ZTERMA8GA1UECgwIQ2hlZXNlIDIxHzAdBgNVBAsMFlNpbXBsZSBTaWduaW5nIFNl
Y3Rpb24xITAfBgNVBAsMGFNpbXBsZSBTaWduaW5nIFNlY3Rpb24gMjEaMBgGA1UE
AwwRU2ltcGxlIFNpZ25pbmcgQ0ExHDAaBgNVBAMME1NpbXBsZSBTaWduaW5nIENB
IDIxCzAJBgNVBAYTAkZSMQswCQYDVQQGEwJVUzERMA8GA1UEBwwIVE9VTE9VU0Ux
DTALBgNVBAcMBExZT04xFjAUBgNVBAgMDVNpZ25pbmcgU3RhdGUxGDAWBgNVBAgM
D1NpZ25pbmcgU3RhdGUgMjEhMB8GCSqGSIb3DQEJARYSc2ltcGxlQHNpZ25pbmcu
Y29tMSIwIAYJKoZIhvcNAQkBFhNzaW1wbGUyQHNpZ25pbmcuY29tMB4XDTE4MTIw
NjExMTAxNloXDTIwMTIwNTExMTAxNlowggF2MRMwEQYKCZImiZPyLGQBGRYDb3Jn
MRYwFAYKCZImiZPyLGQBGRYGY2hlZXNlMQ8wDQYDVQQKDAZDaGVlc2UxETAPBgNV
BAoMCENoZWVzZSAyMR8wHQYDVQQLDBZTaW1wbGUgU2lnbmluZyBTZWN0aW9uMSEw
HwYDVQQLDBhTaW1wbGUgU2lnbmluZyBTZWN0aW9uIDIxFTATBgNVBAMMDCouY2hl
ZXNlLm9yZzEVMBMGA1UEAwwMKi5jaGVlc2UuY29tMQswCQYDVQQGEwJGUjELMAkG
A1UEBhMCVVMxETAPBgNVBAcMCFRPVUxPVVNFMQ0wCwYDVQQHDARMWU9OMRkwFwYD
VQQIDBBDaGVlc2Ugb3JnIHN0YXRlMRkwFwYDVQQIDBBDaGVlc2UgY29tIHN0YXRl
MR4wHAYJKoZIhvcNAQkBFg9jZXJ0QGNoZWVzZS5vcmcxHzAdBgkqhkiG9w0BCQEW
EGNlcnRAc2NoZWVzZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQDed/qNA3AwOd1RG8xg26laE7Gv/izGOJuICg+O2RuhHa8NZuQTW7xdNpLXXtD6
iCnTeOGB3piyqSI/v4qvEpJj1KnD8uR+0tyixTkceuvXEnBjLkFH4PAI6Ny+CQHs
KAmvNdd5nFA10Wvlh3s09tIxZR0YQmlsBBGD/kSukJItC3U5V2LmFy9HK8dT3RAt
yeMGE9K5umMuPH2Da9aJycydTb+f6KN72siZK7pm1o74QUGgydBeyBGkVUqTg4dj
BGNBnPtoBGfCcS/yZR0CXRXbLNkEaYXCfQ3qO6yF+NSPD8VwskXh7LJUC+n3grSb
Gy25JdSryo9bRD4V3bh/t+75AgMBAAGjgeAwgd0wDgYDVR0PAQH/BAQDAgWgMAkG
A1UdEwQCMAAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB0GA1UdDgQW
BBSUunN4oof7WCgoz5g7wkVwFm4pLzAfBgNVHSMEGDAWgBQeUqLoVNU369WoHeTC
BB034vdwAzBhBgNVHREEWjBYggwqLmNoZWVzZS5vcmeCDCouY2hlZXNlLm5ldIIM
Ki5jaGVlc2UuY29thwQKAAEAhwQKAAECgQ90ZXN0QGNoZWVzZS5vcmeBD3Rlc3RA
Y2hlZXNlLm5ldDANBgkqhkiG9w0BAQUFAAOCAQEAdmsFsA40EbGDmZHcrhviCBWL
FrKbJxwCrLXfG9DQdaQrLFxl7Zmr983+OD/DmiIxG6yMHML5XdR1ei5yx4WpBK+f
KszTlnXwjsfGdkisRaS5Ah4vwBXEBwiSyydQZ6HIBcU6s6ZIvuvVWauiG5UwcRNb
CppzO2DMENBqx+XXiy/5Lpjy/4EUJAnjS1VXCRoidPH2QBMxQ4lxCpYaBYIfgzqH
mxcl71pV8i3NDU3kgVi2440JYpoMveTlXPCV2svHNCw0X238YHsSW4b93yGJO0gI
ML9n/4zmm1PMhzZHcEA72ZAq0tKCxpz10djg5v2qL5V+Oaz8TtTOZbPsxpiKMQ==
-----END CERTIFICATE-----
```
#### 5.2.2 pem
pem选项将转义的证书设置为X-Forwarded-Tls-Client-Cert头
在这个例子中，它是-----BEGIN CERTIFICATE----- 和-----END CERTIFICATE----- 分界线之间的部分:
```yaml
-----BEGIN CERTIFICATE-----
MIIGWjCCBUKgAwIBAgIBATANBgkqhkiG9w0BAQUFADCCAYQxEzARBgoJkiaJk/Is
ZAEZFgNvcmcxFjAUBgoJkiaJk/IsZAEZFgZjaGVlc2UxDzANBgNVBAoMBkNoZWVz
ZTERMA8GA1UECgwIQ2hlZXNlIDIxHzAdBgNVBAsMFlNpbXBsZSBTaWduaW5nIFNl
Y3Rpb24xITAfBgNVBAsMGFNpbXBsZSBTaWduaW5nIFNlY3Rpb24gMjEaMBgGA1UE
AwwRU2ltcGxlIFNpZ25pbmcgQ0ExHDAaBgNVBAMME1NpbXBsZSBTaWduaW5nIENB
IDIxCzAJBgNVBAYTAkZSMQswCQYDVQQGEwJVUzERMA8GA1UEBwwIVE9VTE9VU0Ux
DTALBgNVBAcMBExZT04xFjAUBgNVBAgMDVNpZ25pbmcgU3RhdGUxGDAWBgNVBAgM
D1NpZ25pbmcgU3RhdGUgMjEhMB8GCSqGSIb3DQEJARYSc2ltcGxlQHNpZ25pbmcu
Y29tMSIwIAYJKoZIhvcNAQkBFhNzaW1wbGUyQHNpZ25pbmcuY29tMB4XDTE4MTIw
NjExMTAxNloXDTIwMTIwNTExMTAxNlowggF2MRMwEQYKCZImiZPyLGQBGRYDb3Jn
MRYwFAYKCZImiZPyLGQBGRYGY2hlZXNlMQ8wDQYDVQQKDAZDaGVlc2UxETAPBgNV
BAoMCENoZWVzZSAyMR8wHQYDVQQLDBZTaW1wbGUgU2lnbmluZyBTZWN0aW9uMSEw
HwYDVQQLDBhTaW1wbGUgU2lnbmluZyBTZWN0aW9uIDIxFTATBgNVBAMMDCouY2hl
ZXNlLm9yZzEVMBMGA1UEAwwMKi5jaGVlc2UuY29tMQswCQYDVQQGEwJGUjELMAkG
A1UEBhMCVVMxETAPBgNVBAcMCFRPVUxPVVNFMQ0wCwYDVQQHDARMWU9OMRkwFwYD
VQQIDBBDaGVlc2Ugb3JnIHN0YXRlMRkwFwYDVQQIDBBDaGVlc2UgY29tIHN0YXRl
MR4wHAYJKoZIhvcNAQkBFg9jZXJ0QGNoZWVzZS5vcmcxHzAdBgkqhkiG9w0BCQEW
EGNlcnRAc2NoZWVzZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQDed/qNA3AwOd1RG8xg26laE7Gv/izGOJuICg+O2RuhHa8NZuQTW7xdNpLXXtD6
iCnTeOGB3piyqSI/v4qvEpJj1KnD8uR+0tyixTkceuvXEnBjLkFH4PAI6Ny+CQHs
KAmvNdd5nFA10Wvlh3s09tIxZR0YQmlsBBGD/kSukJItC3U5V2LmFy9HK8dT3RAt
yeMGE9K5umMuPH2Da9aJycydTb+f6KN72siZK7pm1o74QUGgydBeyBGkVUqTg4dj
BGNBnPtoBGfCcS/yZR0CXRXbLNkEaYXCfQ3qO6yF+NSPD8VwskXh7LJUC+n3grSb
Gy25JdSryo9bRD4V3bh/t+75AgMBAAGjgeAwgd0wDgYDVR0PAQH/BAQDAgWgMAkG
A1UdEwQCMAAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB0GA1UdDgQW
BBSUunN4oof7WCgoz5g7wkVwFm4pLzAfBgNVHSMEGDAWgBQeUqLoVNU369WoHeTC
BB034vdwAzBhBgNVHREEWjBYggwqLmNoZWVzZS5vcmeCDCouY2hlZXNlLm5ldIIM
Ki5jaGVlc2UuY29thwQKAAEAhwQKAAECgQ90ZXN0QGNoZWVzZS5vcmeBD3Rlc3RA
Y2hlZXNlLm5ldDANBgkqhkiG9w0BAQUFAAOCAQEAdmsFsA40EbGDmZHcrhviCBWL
FrKbJxwCrLXfG9DQdaQrLFxl7Zmr983+OD/DmiIxG6yMHML5XdR1ei5yx4WpBK+f
KszTlnXwjsfGdkisRaS5Ah4vwBXEBwiSyydQZ6HIBcU6s6ZIvuvVWauiG5UwcRNb
CppzO2DMENBqx+XXiy/5Lpjy/4EUJAnjS1VXCRoidPH2QBMxQ4lxCpYaBYIfgzqH
mxcl71pV8i3NDU3kgVi2440JYpoMveTlXPCV2svHNCw0X238YHsSW4b93yGJO0gI
ML9n/4zmm1PMhzZHcEA72ZAq0tKCxpz10djg5v2qL5V+Oaz8TtTOZbPsxpiKMQ==
-----END CERTIFICATE-----
```

**提取数据：**分隔符和\n将被删除。如果有多个证书，它们之间用一个","分隔
**X-Forwarded-Tls-Client-Cert 的值可能超过web服务器头部大小限制：**网络服务器的头文件大小限制通常在4kb和8kb之间。如果这变成了一个问题，并且如果重新配置服务器以允许更大的头文件是不可行的，我们可以通过使用下面描述的信息选项，只选择证书中有趣的部分来缓解这个问题。(并将 pem 设置为 false）

#### 5.2.3 info
`info`选项选择了你想添加到`X-Forwarded-Tls-Client-Cert-Info`头的特定客户证书细节.
头部的值是所有选定的证书细节的转义连接。但在下文中，除非另有规定，为了便于阅读，所有标头值的例子都是未转义的.
下面的例子显示了这样一个串联，当所有可用的字段被选中时:
```yaml
Subject="DC=org,DC=cheese,C=FR,C=US,ST=Cheese org state,ST=Cheese com state,L=TOULOUSE,L=LYON,O=Cheese,O=Cheese 2,CN=*.example.com";Issuer="DC=org,DC=cheese,C=FR,C=US,ST=Signing State,ST=Signing State 2,L=TOULOUSE,L=LYON,O=Cheese,O=Cheese 2,CN=Simple Signing CA 2";NB="1544094616";NA="1607166616";SAN="*.example.org,*.example.net,*.example.com,test@example.org,test@example.net,10.0.1.0,10.0.1.2"

```
如果有一个以上的证书，它们之间用","分隔.

##### info.notAfter
将info.notAfter选项设置为true，以便从Validity部分添加Not After信息.
数据取自以下证书部分:
```yaml
    Validity
        Not After : Dec  5 11:10:16 2020 GMT
```
它的格式在标题中如下:
```yaml
NA="1607166616"
```
##### info.notBefore
将`info.notBefore`选项设置为`true`，以便从`Validity`部分添加`Not Before`信息
数据取自以下证书部分:
```yaml
Validity
    Not Before: Dec  6 11:10:16 2018 GMT
```
它的格式在标题中如下:
```yaml
NB="1544094616"
```
##### info.sans
将 `info.sans` 选项设置为 `true`，以便从`Subject Alternative Name`部分添加`Subject Alternative Name`信息
数据取自以下证书部分:
```yaml
X509v3 Subject Alternative Name:
    DNS:*.example.org, DNS:*.example.net, DNS:*.example.com, IP Address:10.0.1.0, IP Address:10.0.1.2, email:test@example.org, email:test@example.net
```
它的格式在标题中如下:
```yaml
SAN="*.example.org,*.example.net,*.example.com,test@example.org,test@example.net,10.0.1.0,10.0.1.2"

```
多个值：`SAN`被","隔开

##### info.subject
`info.subject`选择你想添加到`X-Forwarded-Tls-Client-Cert-Info`头中的特定客户证书主题细节
数据取自以下证书部分:
```yaml
Subject: DC=org, DC=cheese, O=Cheese, O=Cheese 2, OU=Simple Signing Section, OU=Simple Signing Section 2, CN=*.example.org, CN=*.example.com, C=FR, C=US, L=TOULOUSE, L=LYON, ST=Cheese org state, ST=Cheese com state/emailAddress=cert@example.org/emailAddress=cert@sexample.com
```
###### INFO.SUBJECT.COUNTRY
将`info.subject.country`选项设置为`true`，以将国家信息添加到主题中.
这些数据是用`C`键从主题部分提取的.
并且在标题中的格式如下:
```yaml
C=FR,C=US
```
###### INFO.SUBJECT.PROVINCE
将info.subject.province选项设置为true，以将省份信息添加到主题中.
这些数据是用`ST`键从主题部分提取的.
并且在标题中的格式如下:
```yaml
ST=Cheese org state,ST=Cheese com state
```
###### INFO.SUBJECT.LOCALITY
将`info.subject.locality`选项设置为 `true`，以将位置信息添加到主题中.
这些数据是从带有`L`键的主题部分中提取的.
并且在标题中的格式如下:
```yaml
L=TOULOUSE,L=LYON
```
###### INFO.SUBJECT.ORGANIZATION
将`info.subject.organization`选项设置为`true`，以将组织信息添加到主题中.
这些数据是从带有`O`键的主题部分提取的.
并且在标题中的格式如下:
```yaml
O=Cheese,O=Cheese 2
```
###### INFO.SUBJECT.ORGANIZATIONALUNIT
将`info.subject.organizationalUnit`选项设置为`true`，以将`organizationalUnit`信息添加到主题中.
这些数据来自带有`OU`键的主题部分.
并且在标题中的格式如下:
```yaml
OU=Cheese Section,OU=Cheese Section 2
```
###### INFO.SUBJECT.COMMONNAME
将`info.subject.commonName`选项设置为`true`以将`commonName`信息添加到主题中.
这些数据是从带有`CN`键的主题部分中提取的.
并且在标题中的格式如下:
```yaml
CN=*.example.com
```
###### INFO.SUBJECT.SERIALNUMBER
将`info.subject.serialNumber`选项设置为`true`以将`serialNumber`信息添加到主题中.
该数据从带有`SN`键的主题部分提取.
并且在标题中的格式如下:
```yaml
SN=1234567890
```
###### INFO.SUBJECT.DOMAINCOMPONENT
将`info.subject.domainComponent`选项设置为`true`以将`domainComponent`信息添加到主题中.
这些数据是从带有`DC`键的主题部分中提取的.
并且在标题中的格式如下:
```yaml
DC=org,DC=cheese
```
##### info.issuer
`info.issuer`选择了你想添加到`X-Forwarded-Tls-Client-Cert-Info`头中的特定客户证书签发者细节.
这些数据取自以下证书部分:
```yaml
Issuer: DC=org, DC=cheese, O=Cheese, O=Cheese 2, OU=Simple Signing Section, OU=Simple Signing Section 2, CN=Simple Signing CA, CN=Simple Signing CA 2, C=FR, C=US, L=TOULOUSE, L=LYON, ST=Signing State, ST=Signing State 2/emailAddress=simple@signing.com/emailAddress=simple2@signing.com
```
###### INFO.ISSUER.COUNTRY
将`info.issuer.country`选项设置为`true`，将国家信息添加到发行人中.
这些数据是用`C`键从签发人部分提取的.
而它在标题中的格式如下:
```yaml
C=FR,C=US
```
###### INFO.ISSUER.PROVINCE
将`info.issuer.province`选项设置为`true`，以将省份信息添加到发行人中.
这些数据是从带有`ST`键的签发人部分提取的.
而它在标题中的格式如下:
```yaml
ST=Signing State,ST=Signing State 2
```
###### INFO.ISSUER.LOCALITY
将`info.issuer.locality`选项设置为 `true`，以将位置信息添加到签发人中.
该数据从带有`L`键的签发人部分提取.
而它在标题中的格式如下:
```yaml
L=TOULOUSE,L=LYON
```
###### INFO.ISSUER.ORGANIZATION
将`info.issuer.organization`选项设置为`true`，将组织信息添加到发行人中.
这些数据是从带有`O`键的签发人部分提取的.
而它在标题中的格式如下:
```yaml
O=Cheese,O=Cheese 2
```
###### INFO.ISSUER.COMMONNAME
将`info.issuer.commonName`选项设置为`true`，将`commonName`信息添加到签发人中.
该数据取自带有`CN`密钥的签发人部分.
而它在标题中的格式如下:
```yaml
CN=Simple Signing CA 2
```
###### INFO.ISSUER.SERIALNUMBER
将`info.issuer.serialNumber`选项设置为`true`，将`serialNumber`信息添加到签发人中.
该数据从带有`SN`密钥的发行者部分提取.
而它在标题中的格式如下:
```yaml
SN=1234567890
```
###### INFO.ISSUER.DOMAINCOMPONENT
将`info.issuer.domainComponent`选项设置为`true`，将`domainComponent`信息添加到签发人中.
这些数据是从带有`DC`密钥的发行者部分提取的.
而它在标题中的格式如下:
```yaml
DC=org,DC=cheese
```
## 6.RateLimit
控制进入一个服务的请求的数量,`RateLimit`中间件确保服务将收到公平数量的请求，并允许人们定义什么是公平.
### 6.1配置示例
```yaml
# 这里，允许每秒平均100个请求.
# 此外，允许有50个请求的突发.
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-ratelimit
spec:
  rateLimit:
    average: 100
    burst: 50
```
### 6.2配置选项
#### 6.2.1 average
`average`(平均速率)是指允许来自特定来源的最大速率，默认为每秒请求数。它的默认值是0，这意味着没有速率限制。
速率实际上是由`average`(平均数)除以`period`(周期)来定义的。因此，对于一个低于1 req/s的速率，需要定义一个大于一秒钟的`period`(周期)
```yaml
# 100 reqs/s
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-ratelimit
spec:
  rateLimit:
    average: 100
```
#### 6.2.2 period
`period`(周期)，与`average`(平均数)相结合，定义了实际的最高费率，如： 它的默认值是1秒.
```yaml
r = average / period
```
```yaml
# 6 reqs/minute
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-ratelimit
spec:
  rateLimit:
    period: 1m
    average: 6
```

#### 6.2.3 burst
`burst`是允许在同一任意小的时间段内通过的最大请求数.它的默认值是1.
```yaml
# 允许有100个请求的突发.
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-ratelimit
spec:
  rateLimit:
    burst: 100
```
#### 6.2.4 sourceCriterion
`sourceCriterion`选项定义了使用什么标准来将请求归类为来自一个共同的来源。如果同时定义了几个策略，将产生一个错误。如果没有设置，默认是使用请求的远程地址字段（作为`ipStrategy`）.

##### sourceCriterion.ipStrategy
`ipStrategy`选项定义了两个参数，用于配置Traefik如何确定客户的IP：`depth`和`excludedIPs`.
###### IPSTRATEGY.DEPTH
depth选项告诉Traefik使用X-Forwarded-For报头并获取位于depth位置的IP(从右侧开始).

- 如果深度大于X-Forwarded-For中的IP总数，则客户端IP为空
- 如果深度的值小于或等于0，则忽略深度

depth和X-Forwarded-For的例子：如果深度设置为2，请求X-Forwarded-For的头部是“10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1”，那么真正的客户端IP是“10.0.0.1”(深度4)，但是用于白名单的IP是“12.0.0.1”(深度2)

| X-Forwarded-For | depth | `clientIP` |
| --- | --- | --- |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | 1 | "13.0.0.1" |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | 3 | "11.0.0.1" |
| "10.0.0.1,11.0.0.1,12.0.0.1,13.0.0.1" | 5 | "" |

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-ratelimit
spec:
  rateLimit:
    sourceCriterion:
      ipStrategy:
        depth: 2
```
###### IPSTRATEGY.EXCLUDEDIPS
**提示：**

- 与名称所暗示的相反，这个选项并不是要把一个IP从速率限制器中排除，因此不能用来停用某些IP的速率限制.
- 如果指定了depth，则忽略excludedip.

excludedIPs是为了解决两类有点不同的使用情况:

1. 区分在同一个(一组)反向代理后面的IP，以便它们中的每一个都独立于其他的，为自己的限速"bucket"做出贡献(参考 [漏桶的比喻](https://wikipedia.org/wiki/Leaky_bucket))。在这种情况下，`excludedIPs`应该被设置为与要被排除的`X-Forwarded-For IPs`列表相匹配，以便找到实际的客户IP.

每个IP作为一个不同的来源:

| X-Forwarded-For | excludedIPs | clientIP |
| --- | --- | --- |
| "10.0.0.1,11.0.0.1,12.0.0.1" | "11.0.0.1,12.0.0.1" | "10.0.0.1" |
| "10.0.0.2,11.0.0.1,12.0.0.1" | "11.0.0.1,12.0.0.1" | "10.0.0.2" |

2. 将一组IP（也在一组共同的反向代理后面）组合在一起，使它们被认为是同一来源，并且都对同一速率限制桶作出贡献.

将同一源ip组合在一起:

| X-Forwarded-For | excludedIPs | clientIP |
| --- | --- | --- |
| "10.0.0.1,11.0.0.1,12.0.0.1" | "12.0.0.1" | "11.0.0.1" |
| "10.0.0.2,11.0.0.1,12.0.0.1" | "12.0.0.1" | "11.0.0.1" |
| "10.0.0.3,11.0.0.1,12.0.0.1" | "12.0.0.1" | "11.0.0.1" |

为了完整起见，下面是一些额外的例子来说明匹配是如何进行的。对于一个给定的请求，`X-Forwarded-For` IPs的列表被从最近到最远的排除IPs池检查，并且返回第一个不在池中的IP（如果有的话）。
客户端IP的匹配:

| X-Forwarded-For | excludedIPs | clientIP |
| --- | --- | --- |
| "10.0.0.1,11.0.0.1,13.0.0.1" | "11.0.0.1" | "13.0.0.1" |
| "10.0.0.1,11.0.0.1,13.0.0.1" | "15.0.0.1,16.0.0.1" | "13.0.0.1" |
| "10.0.0.1,11.0.0.1" | "10.0.0.1,11.0.0.1" | "" |

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-ratelimit
spec:
  rateLimit:
    sourceCriterion:
      ipStrategy:
        excludedIPs:
        - 127.0.0.1/32
        - 192.168.1.7
```
##### 
##### sourceCriterion.requestHeaderName
用于对传入请求进行分组的报头的名称
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-ratelimit
spec:
  rateLimit:
    sourceCriterion:
      requestHeaderName: username
```
##### sourceCriterion.requestHost
是否将请求主机视为源
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-ratelimit
spec:
  rateLimit:
    sourceCriterion:
      requestHost: true
```

## 7.RedirectRegex
将客户端重定向到不同的位置,RedirectRegex使用重码匹配和替换来重定向一个请求.
### 7.1配置示例
```yaml
# 用域名替换重定向
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-redirectregex
spec:
  redirectRegex:
    regex: ^http://localhost/(.*)
    replacement: http://mydomain/${1}
```
### 7.2配置选项
#### 7.2.1 permanent
将`permanent`选项设为 `true` 以应用永久重定向
#### 7.2.2 regex
`regex`选项是正则表达式，用于匹配和捕获请求URL中的元素。
**提示：**

- 正则表达式和替换可以使用在线工具，如[Go Playground](https://go.dev/play/)或[Regex101](https://regex101.com/)进行测试
- 在YAML中定义正则表达式时，任何转义字符都需要转义两次： example.com需要写成 example\\.com

#### 7.2.3 replacement
替换选项定义了如何修改URL以拥有新的目标URL
**注意：**在定义替换扩展变量时应注意：1x相当于${1x}，而不是${1}x（参见 [Regexp.Expand](https://golang.org/pkg/regexp/#Regexp.Expand)），因此应使用${1}语法

## 8.RedirectScheme
将客户端重定向到一个不同的Scheme/Port
RedirectScheme将请求从一个scheme/port重定向到另一个

### 8.1 配置示例
```yaml
# Redirect to https
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-redirectscheme
spec:
  redirectScheme:
    scheme: https
    permanent: true
```
### 8.2配置选项
#### 8.2.1 permanent
将permanent选项设置为true以应用永久重定向
```yaml
# Redirect to https
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-redirectscheme
spec:
  redirectScheme:
    # ...
    permanent: true
```
#### 8.2.2 scheme
scheme选项定义了新URL的模式
```yaml
# Redirect to https
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-redirectscheme
spec:
  redirectScheme:
    scheme: https
```
#### 8.2.3 port
`port`选项定义新URL的端口
```yaml
# Redirect to https
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-redirectscheme
spec:
  redirectScheme:
    # ...
    port: "443" # 此配置中的Port是字符串，而不是数值
```

## 9.ReplacePath
在转发请求前更新路径，替换请求URL的路径
### 9.1 配置示例
```yaml
# 用/foo替换路径
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-replacepath
spec:
  replacePath:
    path: /foo
```
### 9.2 配置选项
#### 9.2.1常规
ReplacePath中间件：

- 用指定的路径替换实际路径
- 在X-Replaced-Path头中存储原始路径
#### 9.2.2 path
path选项定义了在请求URL中作为替换的路径

## 10.ReplacePathRegex
在转发请求前更新路径（使用Regex），ReplaceRegex使用regex匹配和替换来替换一个URL的路径。
### 10.1 配置示例
```yaml
# 用regex替换路径
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-replacepathregex
spec:
  replacePathRegex:
    regex: ^/foo/(.*)
    replacement: /bar/$1
```
### 10.2 配置选项
#### 10.2.1 常规
ReplacePathRegex中间件：

- 用指定的路径替换实际路径
- 在X-Replaced-Path头中存储原始路径

**提示：**

- 正则表达式和替换可以使用在线工具，如[Go Playground](https://go.dev/play/)或[Regex101](https://regex101.com/)进行测试
- 在YAML中定义正则表达式时，任何转义字符都需要转义两次： example.com需要写成 example\\.com
#### 10.2.2 regex
regex选项是正则表达式，用于匹配和捕获请求URL中的路径

#### 10.2.3 replacement
替换选项定义了替换路径格式，其中可以包括捕获的变量。
**注意：**在定义替换扩展变量时应注意：1x相当于${1x}，而不是${1}x（参见 [Regexp.Expand](https://golang.org/pkg/regexp/#Regexp.Expand)），因此应使用${1}语法

## 11.Retry
重试直到成功，如果一个后端服务器没有回复，Retry中间件会向该服务器重新发出一定次数的请求。一旦服务器回复，中间件就会停止重试，不管响应状态如何。Retry中间件有一个可选的配置，以启用指数后退。
### 11.1 配置示例
```yaml
# 用指数退避法重试4次
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-retry
spec:
  retry:
    attempts: 4
    initialInterval: 100ms
```

### 11.2 配置选项
#### 11.2.1 attempts
强制性的，`attempts`选项定义了请求应该被重试的次数。

#### 11.2.2 initialInterval
`initialInterval`选项定义了指数退避系列中的第一个等待时间。最大的间隔时间被计算为`initialInterval`的两倍。如果没有指定，请求将被立即重试。
`initialInterval`的值应以秒为单位，或以有效的持续时间格式提供，见[time.ParseDuration](https://golang.org/pkg/time/#ParseDuration).
## 12.StripPrefix
在转发请求前从路径中删除前缀，从URL路径中删除指定的前缀。
### 12.1 配置示例
```yaml
# 剥去前缀 /foobar 和 /fiibar
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-stripprefix
spec:
  stripPrefix:
    prefixes:
      - /foobar
      - /fiibar
```
### 12.2 配置选项
#### 12.2.1 常规
StripPrefix中间件剥离匹配的路径前缀并将其存储在`X-Forwarded-Prefix`头中。
**提醒：**如果你的后端在根路径（`/`）上监听，但应该在一个特定的前缀上暴露，请使用`StripPrefix`中间件.
#### 12.2.2 prefixes
前缀选项定义了要从请求URL中剥离的前缀.
例如，`/products`也匹配`/products/shoes`和`/products/shirts`.
如果你的后端正在提供资产（例如，图像或JavaScript文件），它可以使用`X-Forwarded-Prefix`头来正确构建相对的URL。使用前面的例子，后端应该返回`/products/shoes/image.png`（而不是`/images.png`，Traefik可能无法将其与同一后端联系起来）
#### 12.2.3 forceSlash
`_可选， Default=true_`
`forceSlash`选项确保产生的剥离路径不是空字符串，必要时用`/`替换。添加这个选项是为了保持这个中间件最初的（非直观的）行为，以避免引入一个破坏性的变化。
建议明确将`forceSlash`设置为`false`。
**行为举例：**

- `forceSlash=true`
| Path | Prefix to strip | Result |
| --- | --- | --- |
| / | / | / |
| /foo | /foo | / |
| /foo/ | /foo | / |
| /foo/ | /foo/ | / |
| /bar | /foo | /bar |
| /foo/bar | /foo | /bar |

- `forceSlash=false`
| Path | Prefix to strip | Result |
| --- | --- | --- |
| / | / | 空的 |
| /foo | /foo | 空的 |
| /foo/ | /foo | / |
| /foo/ | /foo/ | 空的 |
| /bar | /foo | /bar |
| /foo/bar | /foo | /bar |



```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: example
spec:
  stripPrefix:
    prefixes:
      - "/foobar"
    forceSlash: false
```
## 13.StripPrefixRegex
在转发请求前从路径中移除前缀（使用Regex），从URL路径中删除匹配的前缀。
### 13.1 配置示例
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: test-stripprefixregex
spec:
  stripPrefixRegex:
    regex:
      - "/foo/[a-z0-9]+/[0-9]+/"
```

### 13.2 配置选项
#### 13.2.1 常规
`StripPrefixRegex`中间件剥离匹配的路径前缀并将其存储在`X-Forwarded-Prefix`头中。
**提醒：**如果您的后端在root路径（`/`）上聆听，但应在特定的前缀上曝光，请使用`StrippreFixRegex`中间件.
#### 13.2.2 regex
`regex`选项是正则表达式，用于匹配请求URL的路径前缀。例如，`/products`也匹配`/products/shoes`和`/products/shirts`。如果你的后端正在提供资产（例如，图像或JavaScript文件），它可以使用`X-Forwarded-Prefix`头来正确构建相对的URL。使用前面的例子，后端应该返回`/products/shoes/image.png`（而不是`/images.png`，Traefik可能无法将其与同一后端联系起来）。

**提示：**

- 正则表达式和替换可以使用在线工具，如[Go Playground](https://go.dev/play/)或[Regex101](https://regex101.com/)进行测试。
- 在YAML中定义正则表达式时，任何转义字符都需要转义两次： example.com需要写成 example\.com。

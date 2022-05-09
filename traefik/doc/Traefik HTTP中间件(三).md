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

## 2.Headers

管理请求/响应报头

![Headers](https://doc.traefik.io/traefik/assets/img/middleware/headers.png)

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
- 在YAML中定义正则表达式时，任何转义字符都需要转义两次： example\.com需要写成 example\\.com。

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

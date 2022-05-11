# 一、TCP Middlewares
控制连接
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/middleware.png#crop=0&crop=0&crop=1&crop=1&id=UfRgi&originHeight=1730&originWidth=3165&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=#crop=0&crop=0&crop=1&crop=1&id=jYsPu&originHeight=1730&originWidth=3165&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=#crop=0&crop=0&crop=1&crop=1&id=WTTRZ&originHeight=1730&originWidth=3165&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)
## 1.1 配置示例
```yaml
# As a Kubernetes Traefik IngressRoute
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: middlewaretcps.traefik.containo.us
spec:
  group: traefik.containo.us
  version: v1alpha1
  names:
    kind: MiddlewareTCP
    plural: middlewaretcps
    singular: middlewaretcp
  scope: Namespaced

---
apiVersion: traefik.containo.us/v1alpha1
kind: MiddlewareTCP
metadata:
  name: foo-ip-whitelist
spec:
  ipWhiteList:
    sourcerange:
      - 127.0.0.1/32
      - 192.168.1.7

---
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRouteTCP
metadata:
  name: ingressroute
spec:
# more fields...
  routes:
    # more fields...
    middlewares:
      - name: foo-ip-whitelist
```
## 1.2 TCP中间件清单
| 中间件 | 目的 | 区域 |
| --- | --- | --- |
| InFlightConn | 限制了同时连接的数量 | 安全性, 请求生命周期 |
| IPWhiteList | 限制允许的客户端IP | 安全性, 请求生命周期 |


# 二、InFlightConn
限制同时连接的数量。为了主动防止服务被高负荷所淹没，可以限制IP允许的同时连接数。
## 2.1 配置示例
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: MiddlewareTCP
metadata:
  name: test-inflightconn
spec:
  inFlightConn:
    amount: 10
```
## 2.2 配置选项
### 2.2.1 amount
`amount`选项定义允许同时连接的最大数量。如果已经有一定数量的连接打开，中间件将关闭连接。

# 三、IPWhiteList
将客户限制在特定的IP上，`IPWhitelist`接受/拒绝基于客户端IP的连接。
## 3.1 配置示例
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: MiddlewareTCP
metadata:
  name: test-ipwhitelist
spec:
  ipWhiteList:
    sourceRange:
      - 127.0.0.1/32
      - 192.168.1.7
```
## 3.2 配置选项
### 3.2.1 sourceRange
`sourceRange`选项设置允许的ip地址(或使用CIDR表示法允许的ip地址范围)。

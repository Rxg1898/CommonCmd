elasticsearch常用命令

## 1、curl请求认证的es

```shell
curl -XGET  'http://账号:密码@127.0.0.1:9200/_cluster/health?pretty'

示例：curl -XGET  'http://elastic:elastic@127.0.0.1:9200/_cluster/health?pretty'
```




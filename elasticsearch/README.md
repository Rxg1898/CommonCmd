elasticsearch常用命令

## 1、curl请求认证的es

```shell
curl -XGET  'http://账号:密码@127.0.0.1:9200/_cluster/health?pretty'

示例：curl -XGET  'http://elastic:elastic@127.0.0.1:9200/_cluster/health?pretty'
```

## 2、查看索引结构

```shell
curl -XGET http://elastic:elastic@127.0.0.1:9200/demo
```

## 3、删除索引

```shell
curl -XDELETE 'http://elastic:elastic@127.0.0.1:9200/demo?pretty'
curl -XDELETE 'http://elastic:elastic@127.0.0.1:9200/demo'
```

## 4、查询索引

```shell
# 查询单个
curl -XGET 'http://elastic:elastic@127.0.0.1:9200/_cat/indices/demo?v'
# 查询多个
curl -XGET 'http://elastic:elastic@127.0.0.1:9200/_cat/indices/demo-0*?v'
# 查询全部
curl -XGET 'http://elastic:elastic@127.0.0.1:9200/_cat/indices?v'
```

```shell
# 按照size从大到小存储排序
curl -XGET 'http://elastic:elastic@127.0.0.1:9200/_cat/indices?v&s=store.size:desc'

GET /_cat/indices?v&s=store.size:desc

# 按照索引名称反向排序
GET /_cat/indices?v&s=index:desc
```

## 4、备份索引

### elasticsearch-dump工具

```shell
# 仅备份mappings
docker run -v /tmp/es/:/opt/ --rm -ti elasticdump/elasticsearch-dump:v6.86.0 \
  --input=http://elastic:elastic@127.0.0.1:9200/demo \
  --output=/opt/demo.json \
  --type=mapping
```

```shell
# 恢复
# 先创建索引
curl -XPUT 'http://elastic:elastic@127.0.0.1:9200/demo?pretty'
# 恢复索引mappings
docker run -v /tmp/es/:/opt/ --rm -ti elasticdump/elasticsearch-dump:v6.86.0 \
  --input=/opt/demo.json \
  --output=http://elastic:elastic@127.0.0.1:9200/ \
  --type=mapping
```

## 5、索引别名

```shell
# 创建索引别名
curl -XPUT 'http://elastic:elastic@127.0.0.1:9200/demo-2022/_alias/demo_alias'
# 删除索引别名
POST /_aliases
{
    "actions": [
        { "remove": { "index": "demo-2022", "alias": "demo_alias" }}
    ]
}
```

## 6、硬盘限额

为了保护节点数据安全，ES 会定时(`cluster.info.update.interval`，默认 30 秒)检查一下各节点的数据目录磁盘使用情况。在达到 `cluster.routing.allocation.disk.watermark.low` (默认 85%)的时候，新索引分片就不会再分配到这个节点上了。在达到 `cluster.routing.allocation.disk.watermark.high` (默认 90%)的时候，就会触发该节点现存分片的数据均衡，把数据挪到其他节点上去。这两个值不但可以写百分比，还可以写具体的字节数。

```shell
PUT /_cluster/settings
{
    "transient" : {
        "cluster.routing.allocation.disk.watermark.low" : "85%",
        "cluster.routing.allocation.disk.watermark.high" : "90%",
        "cluster.routing.allocation.disk.watermark.flood_stage" : "95%",
        "cluster.info.update.interval" : "30s"
    }
}
```








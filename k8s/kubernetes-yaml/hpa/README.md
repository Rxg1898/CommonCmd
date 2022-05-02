## 介绍

kubectl autoscale ⾃动控制在k8s集群中运⾏的pod数量(⽔平⾃动伸缩)，需要提前设置pod范围及触发条件。k8s从1.1版本开始增加了名称为HPA(Horizontal Pod Autoscaler)的控制器，⽤于实现基于pod中资源(CPU/Memory)利⽤率进⾏对pod的⾃动扩缩容功能的实现，早期的版本只能基于Heapster组件实现对CPU利⽤率做为触发条件，但是在k8s 1.11版本开始使⽤Metrices Server完成数据采集，然后将采集到的数据通过API（Aggregated API，汇总API），例如metrics.k8s.io、custom.metrics.k8s.io、external.metrics.k8s.io，然后再把数据提供给HPA控制器进⾏查询，以实现基于某个资源利⽤率对pod进⾏扩缩容的⽬的。

```
控制管理器默认每隔15s（可以通过–horizontal-pod-autoscaler-sync-period修改）查询metrics的资源使⽤情况
⽀持以下三种metrics指标类型：
 ->预定义metrics（⽐如Pod的CPU）以利⽤率的⽅式计算
 ->⾃定义的Pod metrics，以原始值（raw value）的⽅式计算
 ->⾃定义的object metrics
⽀持两种metrics查询⽅式：
 ->Heapster
 ->⾃定义的REST API
⽀持多metrics
```

## 部署metrics-server

使用metrics-server作为HPA数据源

项目地址：https://github.com/kubernetes-sigs/metrics-server


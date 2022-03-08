## 启动顺序

1. prometheus-pv-pvc.yaml
2. prometheus-configmap.yaml
3. node-exporter.yaml
4. kube-state-metrics.yaml
5. prometheus-deployment.yaml

## 腾讯云-cfs-pv实例

```yaml
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: prometheus-pv
  namespace: ops
spec:
  capacity:
    storage: 100Gi
  accessModes:
    - ReadWriteMany
  persistentVolumeReclaimPolicy: Retain
  storageClassName: storageClass名             # cfs存储
  volumeMode: Filesystem
  csi:
    driver: com.tencent.cloud.csi.cfs
    volumeAttributes:
      fsid: ftcyuyer
      host: cfs地址                            # cfs地址
      path: /prometheus                        # cfs路径
      vers: "4"
    volumeHandle: prometheus-pv
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: prometheus-pvc
  namespace: ops
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 100Gi
  storageClassName: hk-cmp-prod
  volumeMode: Filesystem
  volumeName: prometheus-pv
```

## 阿里云-动态nas-pvc-实例

```yaml
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: prometheus-pvc
  namespace: ops
spec:
  accessModes:
    - ReadWriteMany
  storageClassName: alicloud-nas-subpath
  resources: 
    requests:
      storage: 100Gi
```






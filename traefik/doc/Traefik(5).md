# 一、介绍
插件支持是一项强大的功能，允许开发人员向 Traefik 添加新功能并定义新行为。例如，插件可以修改请求或标头、发出重定向、添加身份验证等，提供与 Traefik中间件类似的功能。<br />然而，与传统的中间件不同，插件是动态加载的，并由[嵌入式解释器yaegi](https://github.com/traefik/yaegi)执行。无需编译二进制文件，所有插件都是 100% 跨平台的，这使得它们易于开发并与更广泛的 Traefik 社区共享<br />`Traefik v2.3 及更高版本提供对插件的支持`<br />`插件可能会以不希望的方式修改 Traefik 的行为。向生产 Traefik 实例添加新插件时要小心。`
# 二、插件和 Traefik Pilot
Traefik 与 Traefik Pilot 一起启用插件生态系统。Traefik 操作员可以从在线目录中浏览和安装插件，该目录可从Traefik Pilot 仪表板的[插件选项卡中获得](https://pilot.traefik.io/plugins)<br />![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/%E6%8F%92%E4%BB%B601.png)<br />选择插件的磁贴会打开一个描述插件功能的页面，以及可选的可用配置选项.<br />![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/traefik/%E6%8F%92%E4%BB%B602.png)<br />在那里，选择安装插件将显示必要的代码，添加到Traefik代理的静态 and/or 动态配置中以完成安装过程
# 三、安装插件
对于一个特定的Traefik实例来说，一个插件要被激活，它必须在静态配置中被声明。当你选择安装插件时，要添加的代码是由Traefik Pilot UI提供的。

插件完全在启动过程中被解析和加载，这使得Traefik能够检查代码的完整性，并在早期捕获错误。如果在加载过程中发生错误，该插件将被禁用。

`需要重新启动：出于安全考虑，在Traefik运行时，无法启动一个新的插件或修改现有的插件`

一旦加载，中间件插件的行为就像静态编译的中间件。它们的实例化和行为是由动态配置驱动的。

## 3.1 静态配置
在下面的例子中，我们添加了`blockpath`和`rewritebody`插件:
```yaml
--entryPoints.web.address=:80
--pilot.token=xxxxxxxxx
--experimental.plugins.block.modulename=github.com/traefik/plugin-blockpath
--experimental.plugins.block.version=v0.2.0
--experimental.plugins.rewrite.modulename=github.com/traefik/plugin-rewritebody
--experimental.plugins.rewrite.version=v0.3.0
```
## 3.2 动态配置
一些插件需要通过添加动态配置来配置。对于`bodyrewrite`插件，例如:
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: my-rewritebody
spec:
  plugin:
    rewrite:
      rewrites:
        - regex: example
          replacement: test
```

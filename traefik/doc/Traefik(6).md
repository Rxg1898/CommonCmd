# 一、介绍
[Traefik](https://traefik.io/)插件架构使开发人员可以轻松创建新插件、修改现有插件以及与 Traefik 社区共享插件。<br />Traefik 插件是使用`Go 语言`开发的，而Traefik中间件插件只是一个`Go 包`，它提供了一个`http.Handler`执行特定请求和响应处理的包。<br />然而，插件并没有被预编译和链接，而是由[Yaegi 动态](https://github.com/traefik/yaegi)执行，Yaegi是一个嵌入在 Traefik 应用程序代理中的 Go 解释器。<br />这意味着插件不需要编译，也不需要复杂的工具链来开始。开发 Traefik 插件的过程可与 Web 浏览器扩展相媲美。<br />`插件可能会以不希望的方式修改 Traefik 的行为.向生产 Traefik 实例添加新插件时要小心.`
# 二、封装插件
如`Traefik 使用插件(五)`中所述，Traefik 管理员可以从 Traefik Pilot 仪表板的目录中浏览插件并将其添加到他们的实例中。<br />从技术角度来看，每个插件的实际代码都存储和托管在一个公共的GitHub仓库中。每隔30分钟，Traefik Pilot就会轮询GitHub，找到符合Traefik插件标准的仓库，并将它们添加到目录中。
## 2.1 前提条件
为了被 Traefik Pilot 识别，你的插件存储库必须满足以下条件：

- `traefik-plugin`必须设置主题
- `.traefik.yml`清单必须存在并具有有效的内容

此外，Traefik Pilot从Go模块代理处获取你的源代码，所以你的插件必须使用[git tag](https://git-scm.com/book/en/v2/Git-Basics-Tagging)进行版本控制！<br />如果你的资源库不能满足这些前提条件，Traefik Pilot将无法识别它，你的插件将不会被添加到目录中！

## 2.2 插件声明
清单也是必须的，它应该被命名为`.traefik.yml`并存储在项目的根目录下。<br />这个YAML文件为Traefik Pilot提供了关于你的插件的信息，例如描述、全名等。
```yaml
# 在Traefik Pilot网页用户界面上显示的你的插件的名称
displayName: Name of your plugin

# 目前"中间件"是唯一可用的类型
type: middleware

# 插件的导入路径
import: github.com/username/my-plugin

# 简要描述你的插件是做什么的
summary: Description of what my plugin is doing

# 与该插件相关的图标（可选）
iconPath: foo/icon.png
bannerPath: foo/banner.png

# 插件的配置数据
# 这是必须的
# Traefik Pilot将尝试使用你提供的配置执行该插件，作为其启动有效性测试的一部分
testData:
  Headers:
    Foo: Bar
```
### 2.2.1 配置说明

- `displayName`: Traefik Pilot web UI 中显示的插件名称
- `type`: 目前`middleware`是唯一可用的类型
- `import`：插件的导入路径
- `summary`：简要说明你的插件是做什么的
- `testData`：插件的配置数据。这是必填的，Traefik Pilot 将尝试使用你提供的配置执行该插件，作为其启动有效性测试的一部分
- `iconPath(可选)`：存储库中的本地路径，用于显示项目的图标
- `bannerPath(可选)`：存储库中的本地路径，当你在社交媒体上分享你的插件页面时，将使用该图片

项目的根目录中还应该有一个`go.mod`文件。Traefik Pilot 将使用此文件来验证项目的名称。
## 2.3 开发者模式
对于那些喜欢在将插件部署到 GitHub 之前私下开发插件的人，Traefik 还提供了一种可用于临时测试的开发人员模式。<br />要在开发模式下部署插件，需要同时更改静态和动态配置。静态配置必须定义模块名称（通常用于 Go 包）和[Go 工作区](https://golang.org/doc/gopath_code.html#Workspaces)的路径，该路径可以是存储在本地 GOPATH 环境变量或任何其他路径中的内容。动态配置必须引用标签`dev`。<br />`静态配置`：
```yaml
# 插件将从'/plugins/go/src/github.com/traefik/plugindemo'路径加载
pilot:
  token: xxxxx

experimental:
  devPlugin:
    goPath: /plugins/go
    moduleName: github.com/traefik/plugindemo
```
`动态配置`：
```yaml
http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - my-plugin

  services:
   service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000

  middlewares:
    my-plugin:
      plugin:
        dev:
          headers:
            Foo: Bar
```
`注意：一次只能在dev模式下测试一个插件，并且在使用dev模式时，Traefik将在30分钟后关闭`
# 三、自定义开发插件
插件包必须定义以下导出的Go对象：

- 一个`type Config struct { ... }` 结构体，里面字段任意
- 一个函数`func CreateConfig() *Config`
- 一个函数`func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error)`
```go
// 打包example（一个插件的例子）
package example

import (
    "context"
    "net/http"
)

// 插件的配置
type Config struct {
    // ...
}

// 创建默认的插件配置
func CreateConfig() *Config {
    return &Config{
        // ...
    }
}

// 示例一个插件
type Example struct {
    next     http.Handler
    name     string
    // ...
}

// New一个新的插件
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
    // ...
    return &Example{
        // ...
    }, nil
}

func (e *Example) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    // ...
    e.next.ServeHTTP(rw, req)
}
```

## 3.1 额外的依赖
如果你的插件有任何外部模块依赖，它们必须被[vendored](https://golang.org/ref/mod#vendoring)并包含在GitHub仓库中，不支持[Go modules](https://blog.golang.org/using-go-modules) 
# 四、疑难解答
如果你的插件在集成过程中出现问题，Traefik Pilot将在您的GitHub仓库中创建一个issue来解释这个问题，并将停止尝试添加你的插件，直到你关闭了此issue。<br />为了让 Traefik Pilot 成功导入你的插件，请查阅以下清单：

- 仓库必须设置`traefik-plugin`主题
- 项目的根目录中必须有一个`.traefik.yml`文件描述您的插件，并且它必须具有`testData`用于测试目的的有效属性
- 项目的根目录必须有一个有效的`go.mod`文件
- 项目必须有一个`git tag`版本
- 如果你的插件有任何外部模块依赖，它们必须被[vendored](https://golang.org/ref/mod#vendoring)并包含在GitHub仓库中，不支持[Go modules](https://blog.golang.org/using-go-modules)

# 五、示例代码
完整的`demo插件例子`：[plugindemo](https://github.com/traefik/plugindemo)



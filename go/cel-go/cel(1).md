# 一、cel简介
## 1.1 什么是CEL
CEL是一种非图灵的完整`表达式语言`，被设计为`快速、可移植和安全执行`。CEL可以单独使用，也可以嵌入到一个更大的产品中。

CEL被设计成一种可以安全执行用户代码的语言。虽然在用户的python代码上盲目地调用eval()是危险的，但你可以安全地执行用户的CEL代码。因为 CEL 防止了会使其性能降低的行为，它可以在纳秒到微秒的时间内安全地进行评估；它是性能关键型应用的理想选择。

CEL评估表达式，这类似于单行函数或lambda表达式。虽然CEL通常用于布尔决策，但它也可用于构建更复杂的对象，如JSON或protobuf消息。
## 1.2 CEL理想应用领域
由于CEL是在`纳秒到微秒`的时间内对AST中的表达式进行评估，所以CEL的理想使用场合是具有性能关键点的应用。将CEL代码编译到AST中不应该在关键路径中进行；理想的应用是配置`经常被执行`而`修改相对不频繁`的应用。

例如，在`对服务的每个HTTP请求中执行安全策略`是CEL的一个理想用例，因为安全策略很少改变，CEL对响应时间的影响可以忽略不计。在这种情况下，CEL返回一个布尔值，即请求是否应该被允许，但它可以返回一个更复杂的消息。

## 1.3 使用CEL的相关资料
下面会进行涵盖常见用例的编码练习。要想更深入地了解语言、语义和功能，请参见[GitHub上的CEL语言定义](https://github.com/google/cel-spec/blob/master/doc/langdef.md) 和[CEL Go文档](https://github.com/google/cel-go#common-expression-language).。

# 二、CEL的核心概念
## 2.1 应用领域
CEL是通用的，并已被用于不同的应用，从路由RPC到定义安全策略。CEL是可扩展的，与应用无关，并为一次编译、多次评估的工作流程而优化。

许多服务和应用程序评估声明式配置。例如，基于角色的访问控制（RBAC）是一个声明式的配置，它产生一个给定角色和一组用户的访问决定。如果声明式配置是80%的用例，那么当用户需要更多的表达能力时，CEL是一个有用的工具，可以完善剩下的20%。
## 2.2 编译
一个表达式被针对环境进行编译。编译步骤产生一个`protobuf`形式的抽象语法树（AST）。编译后的表达式通常会被存储起来供将来使用，以保持尽可能快的评估速度。一个已编译的表达式可以用许多不同的输入进行评估。
## 2.3 表达式
`用户定义表达式`；服务和应用程序定义它的运行环境。一个函数签名声明了输入，并写在CEL表达式的外面。CEL可用的函数库是自动导入的。

在下面的例子中，表达式接受了一个请求对象，并且请求包括一个索赔标记。该表达式返回一个布尔值，表明该索赔令牌是否仍然有效。
```go
// 通过检查"exp"要求，检查JSON网络令牌是否已经过期。
//
// Args:
//   claims - 认证要求.
//   now    - 表示当前系统时间的时间戳.
// 返回true表示令牌已经过期.
//
timestamp(claims["exp"]) < now
```
## 2.4 运行环境
`环境是由服务定义`。嵌入 CEL 的服务和应用程序声明表达式环境。环境是可以在表达式中使用的变量和函数的集合。

基于原语的声明被 CEL 类型检查器使用，以确保表达式中的所有标识符和函数引用被正确声明和使用。
## 2.5 解析表达式
处理一个表达式有三个阶段：`解析`、`检查`和`评估`。CEL最常见的模式是控制平面在配置时对表达式进行解析和检查，并存储AST。<br />![image.png](https://cdn.nlark.com/yuque/0/2022/png/2579230/1655376815539-90fe68f5-a3c6-441b-a564-ec3fcad92f8d.png)<br />在运行时，数据平面反复检索和评估AST。CEL对运行时的效率进行了优化，但`解析和检查不应该在延迟关键的代码路径中进行`。<br />![image.png](https://cdn.nlark.com/yuque/0/2022/png/2579230/1655376871943-425d8117-73e8-4bf3-bf95-60dd1ae12cf9.png)


使用 [ANTLR](https://www.antlr.org/)词典/解析器(Lexer/Parser)语法将CEL从人类可读的表达式解析为抽象的语法树。解析阶段发出一个基于`proto`的抽象语法树，其中AST中的每个Expr节点都包含一个整数ID，用于解析和检查期间产生的元数据。在解析过程中产生的  [syntax.proto](https://github.com/googleapis/googleapis/blob/master/google/api/expr/v1alpha1/syntax.proto) 忠实地代表了表达式的字符串形式中输入内容的抽象表示。

一旦表达式被解析，就可以根据环境对其进行检查，以确保表达式中的所有变量和函数标识符都被声明并且正确使用。类型检查器产生一个 [checked.proto](https://github.com/googleapis/googleapis/blob/master/google/api/expr/v1alpha1/checked.proto) ，其中包括类型、变量和函数解析元数据，可以极大地提高评估效率。

`最佳实践：执行类型检查以提高解析表达式的速度和安全性，即使是像JSON这样类型推理有限的动态数据。`

CEL评估器需要3样东西：

- 任何自定义扩展的函数绑定
- 变量绑定
- 要评估的AST

函数和变量绑定应该匹配用于编译AST时的内容一样。这些输入中的任何一个都可以在多次评估中重复使用，比如一个AST在多个变量绑定集上评估，或者相同的变量被用于多个AST，或者在一个进程的生命周期中使用函数绑定(这是一种常见的情况)。

# 三、官方示例
项目地址：[https://github.com/google/cel-go.git](https://github.com/google/cel-go.git)
```go
// 拉取项目
git clone https://github.com/google/cel-go.git

cd cel-go
// 下载依赖包
go mod tidy
// 运行
go run codelab/codelab.go
```
```go
// 打印内容
=== Exercise 1: Hello World ===

=== Exercise 2: Variables ===

=== Exercise 3: Logical AND/OR ===

=== Exercise 4: Customization ===

=== Exercise 5: Building JSON ===

=== Exercise 6: Building Protos ===

=== Exercise 7: Macros ===

=== Exercise 8: Tuning ===
```
## 3.1 codelab/codelab.go包介绍

- **compile function:  **根据环境解析、检查和输入表达式
- eval function：根据输入计算已编译程序的值
- report function： 打印出评估结果
- 此外，还提供了和**request和**auth ，以协助各种练习的输入构建。
- 还有一些辅助包如下所示：
| Package | 地址 | 描述 |
| --- | --- | --- |
| cel | [cel-go/cel/cel.go](https://github.com/google/cel-go/blob/master/cel/cel.go) | 顶层接口 |
| decls | [cel-go/checker/decls/decls.go](https://github.com/google/cel-go/blob/master/checker/decls/decls.go) | 变量和函数声明实用工具 |
| functions | [cel-go/interpreter/functions/functions.go](https://github.com/google/cel-go/blob/master/interpreter/functions/functions.go) | 运行时捆绑 |
| ref | [cel-go/common/types/ref/reference.go](https://github.com/google/cel-go/blob/master/common/types/ref/reference.go) | 参考文献接口 |


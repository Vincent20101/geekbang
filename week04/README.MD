
#wire的用法
##provider和injector是wire的两个核心概念。
通过提供provider函数，让wire知道如何产生这些依赖对象。wire根据我们定义的injector函数签名，生成完整的injector函数，injector函数是最终我们需要的函数，它将按依赖顺序调用provider。

##provider
provider就是普通的Go函数，可以把它看作是某对象的构造函数，我们通过provider告诉wire该对象的依赖情况

##injector
injector是wire生成的函数，我们通过调用injector来获取我们所需的对象或值，injector会按照依赖关系，按顺序调用provider函数
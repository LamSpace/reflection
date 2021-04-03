# reflection

#### 介绍（Introductions）
**reflection** 是基于Golang的反射机制进行依赖注入、实例方法调用等的简单实现。

1. **reflection.ForMethod(m interface{})** 是调用普通方法的入口，调用 **Call(params ...interface{}) ** 或者 **CallType(params ...interface{})** 可以进行普通方法调用并获得返回结果。其中前者普通方法的参数是严格按照方法定义中参数的类型依次传递，参数可存在相同类型的参数；而后者可不按照方法定义中参数的类型依次传递，前提是方法的所有参数类型均不同。
2. **reflection.ForInstance(v interface{})** 是注入结构体实例并调用实例方法的入口，调用 **Map(key string, value interface{})** 和 **Inject()** 可按照结构体实例的字段名注入依赖，允许存在相同类型的字段；调用 **MapType(value interface{})** 和 **InjectType()** 可按照结构体字段的类型进行注入，前提是结构体中每个字段的类型均不一致。调用 **Invoke(function string, params ...interface{})** 和 **InvokeType(function string, params ...interface{})** 可调用实例的方法，参数的传递与 **Call(params ...interface{})** 和 **CallType(params ...interface{})** 类似。

***

**reflection** is a simple implementation of reflect scheme of Golang to inject dependencies, invoke method and so on.

1. **reflection.ForMethod(m interface{})** is entrance to invoke normal method by using  **Call(params ...interface{}) ** or **CallType(params ...interface{})** and responses can be acquired where types of the former's parameters must be consistent with that in the defination of invoked method with the same order strictly, allowing same-type parameters, while the latter invokes method by types of parameters whose order may not be the same with that in the defination of invoked method if and only if parameters of invoked method differs from each others.
2. **reflection.ForInstance(v interface{})** is entrance to inject dependencies and invoke instance's method of a struct. By calling  **Map(key string, value interface{})** and **Inject()**, fields of specified struct can be injected by field name, allowing parameters own the same type while calling  **MapType(value interface{})** and **InjectType()** can also inject dependencies only when all types of fields are different from each other. After injecting dependencies, using  **Invoke(function string, params ...interface{})** and **InvokeType(function string, params ...interface{})** can invoke a specified method of struct instance and differences between both are consistent with that of  **Call(params ...interface{})** and **CallType(params ...interface{})**.


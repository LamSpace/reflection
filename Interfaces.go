package reflection

type Caller interface {
	Call(params ...interface{}) ([]interface{}, error)
	CallType(params ...interface{}) ([]interface{}, error)
}

type Injector interface {
	Inject() (*reflector, error)
	InjectType() (*reflector, error)
}

type Invoker interface {
	Invoke(function string, params ...interface{}) ([]interface{}, error)
	InvokeType(function string, params ...interface{}) ([]interface{}, error)
}

type Mapper interface {
	Map(key string, value interface{}) *reflector
	MapType(value interface{}) *reflector
}

type Reflector interface {
	Mapper
	Injector
	Invoker
}

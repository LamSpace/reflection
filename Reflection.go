package reflection

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Creates an instance of caller for methods.
func ForMethod(m interface{}) Caller {
	return &caller{value: m}
}

// Creates an instance of reflector for structs.
func ForInstance(value interface{}) Reflector {
	return &reflector{
		value:     value,
		entryName: make(map[string]reflect.Value),
		entryType: make(map[reflect.Type]reflect.Value),
	}
}

// Calls a method with given parameters in order corresponding to invoked method's
// and returns the value of called method and nil if and only if invocation is
// success; otherwise, returns nil and an error.
func (c *caller) Call(params ...interface{}) ([]interface{}, error) {
	t := reflect.TypeOf(c.value)
	if t.Kind() != reflect.Func {
		return nil, errors.New("Called 'method' is not a method. ")
	}
	// Checks number of method's parameters with given parameters.
	if parameterNum := t.NumIn(); parameterNum == len(params) {
		var in = make([]reflect.Value, parameterNum)
		for i := 0; i < parameterNum; i++ {
			parameterType := t.In(i)
			// Checks type of parameters with given parameters.
			if pt := reflect.TypeOf(params[i]); parameterType.Kind() != pt.Kind() {
				return nil, errors.New("Type {" + parameterType.Kind().String() +
					"} of invoked method does not match given type {" + pt.Kind().String() +
					"} of parameter {" + fmt.Sprint(params[i]) + "}.")
			} else {
				in[i] = reflect.ValueOf(params[i])
			}
		}
		// Invokes method and processes returned values.
		res := reflect.ValueOf(c.value).Call(in)
		out := convertResponse(res)
		return out, nil
	} else {
		return nil, errors.New("Number of invoked method's {" + strconv.Itoa(parameterNum) +
			"} does not match given parameters' {" + strconv.Itoa(len(params)) + "}.")
	}
}

// Calls a method with given parameters whose type differs from each other and
// returns the value of called method and nil if and only if invocation is
// success; otherwise, returns nil and an error.
func (c *caller) CallType(params ...interface{}) ([]interface{}, error) {
	t := reflect.TypeOf(c.value)
	if t.Kind() != reflect.Func {
		return nil, errors.New("Called 'method' is not a method. ")
	}
	// Checks number of method's parameters with given parameters.
	if parameterNum := t.NumIn(); parameterNum == len(params) {
		// Maps given parameters with key-value entry.
		var m = make(map[reflect.Type]reflect.Value)
		for i := 0; i < len(params); i++ {
			m[reflect.TypeOf(params[i])] = reflect.ValueOf(params[i])
		}
		var in = make([]reflect.Value, parameterNum)
		for i := 0; i < parameterNum; i++ {
			parameterType := t.In(i)
			if val, ok := m[parameterType]; ok {
				in[i] = val
			} else {
				return nil, errors.New("Type {" + parameterType.Kind().String() +
					"} of invoked method does not exists in given parameters.")
			}
		}
		res := reflect.ValueOf(c.value).Call(in)
		out := convertResponse(res)
		return out, nil
	} else {
		return nil, errors.New("Number of invoked method's {" + strconv.Itoa(parameterNum) +
			"} does not match given parameters' {" + strconv.Itoa(len(params)) + "}.")
	}
}

func (r *reflector) transform() (*reflect.Value, error) {
	v := reflect.ValueOf(r.value)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, errors.New("Parameter of invoked method ForInstance(interface{}) is not a pointer of struct. ")
	}
	return &v, nil
}

// Adds an entry representing a field of struct to modify that one
// and returns an pointer of injector.
func (r *reflector) Map(key string, value interface{}) *reflector {
	r.entryName[key] = reflect.ValueOf(value)
	return r
}

// Adds an entry representing a field of struct whose fields' type are
// different from each other.
func (r *reflector) MapType(value interface{}) *reflector {
	r.entryType[reflect.TypeOf(value)] = reflect.ValueOf(value)
	return r
}

// Injects all entries into the given struct and returns an pointer
// of injector and nil if and only if this method succeed; otherwise,
// returns nil and an error.
func (r *reflector) Inject() (*reflector, error) {
	if len(r.entryType) > 0 && len(r.entryName) == 0 {
		fmt.Println("Maybe method invoked is 'InjectType()', and this method injects nothing.")
	}
	v, err := r.transform()
	if err != nil {
		return nil, err
	}
	t := (*v).Type()
	for i := 0; i < v.NumField(); i++ {
		tf := t.Field(i)
		vf := v.Field(i)
		// Checks whether current field can be modified.
		if vf.CanSet() {
			// Lookups of value with corresponding field name.
			if val, ok := r.entryName[tf.Name]; ok {
				// Checks type of field with given parameter's
				if tf.Type.Kind() == val.Type().Kind() {
					vf.Set(val)
				} else {
					return nil, errors.New("Type {" + tf.Type.Kind().String() + "} of field {" + tf.Name +
						"} does not match given type {" + val.Type().Kind().String() + "} of parameter {" +
						val.String() + "}.")
				}
			} else {
				//fmt.Printf("{ %s } does not have a given value, using default zero-value.\n", tf.Name)
			}
		} else {
			//fmt.Printf("Field { %s } can not be changed\n", vf)
		}
	}
	return r, nil
}

// Injects all entries into the given struct and returns an pointer of
// injector and nil if and only if this method succeed; otherwise, returns
// nil and an error.
func (r *reflector) InjectType() (*reflector, error) {
	if len(r.entryName) > 0 && len(r.entryType) == 0 {
		fmt.Println("Maybe method invoked is 'Inject()', and this method injects nothing.")
	}
	v, err := r.transform()
	if err != nil {
		return nil, err
	}
	t := (*v).Type()
	for i := 0; i < v.NumField(); i++ {
		tf, vf := t.Field(i), v.Field(i)
		if vf.CanSet() {
			if val, ok := r.entryType[tf.Type]; ok {
				vf.Set(val)
			} else {
				//fmt.Printf("{ %s } does not have a given name, using default zero-value.\n", tf.Name)
			}
		} else {
			//fmt.Printf("Field { %s } can not be changed.\n", vf)
		}
	}
	return r, nil
}

// Invokes a method by function name with given parameters and returns responses of invoked
// method and nil if and only if invocation is success; otherwise, returns nil and an error.
func (r *reflector) Invoke(function string, params ...interface{}) ([]interface{}, error) {
	v, err := r.transform()
	if err != nil {
		return nil, err
	}
	// Checks existence of invoked method.
	if funcExist := v.MethodByName(function).Kind() == reflect.Invalid; !funcExist {
		invokedFunc, invokedFuncType := v.MethodByName(function), v.MethodByName(function).Type()
		// Checks number of invoked method's parameters with given one.
		if numParams := invokedFuncType.NumIn(); numParams == len(params) {
			var in = make([]reflect.Value, numParams)
			for i := 0; i < numParams; i++ {
				// Checks type of parameter of invoked method with given one.
				if ift, rft := invokedFuncType.In(i), reflect.TypeOf(params[i]); ift.Kind() != rft.Kind() {
					return nil, errors.New("Type {" + ift.Kind().String() + "} of invoked method's " +
						"parameter does not match given parameter {" + fmt.Sprint(params[i]) + "} type {" +
						rft.Kind().String() + "}.")
				} else {
					in[i] = reflect.ValueOf(params[i])
				}
			}
			// Invokes method and processes value of invoked method.
			res := invokedFunc.Call(in)
			out := convertResponse(res)
			return out, nil
		} else {
			return nil, errors.New("Number of invoked method's parameter {" + strconv.Itoa(numParams) +
				"} does not match given parameters' {" + strconv.Itoa(len(params)) + "}.")
		}
	} else {
		return nil, errors.New("Method Invoked {" + function + "} does not exists. ")
	}
}

// Invokes a method by function name with given parameters whose type are different from each other
// and returns responses of invoked method and nil if and only if invocation is success; otherwise,
// returns nil and an error.
func (r *reflector) InvokeType(function string, params ...interface{}) ([]interface{}, error) {
	v, err := r.transform()
	if err != nil {
		return nil, err
	}
	// Checks existence of invoked method.
	if funcExist := v.MethodByName(function).Kind() == reflect.Invalid; !funcExist {
		invokedFunc, invokedFuncType := v.MethodByName(function), v.MethodByName(function).Type()
		// Checks number of parameters of invoked method with given parameters.
		if numParams := invokedFuncType.NumIn(); numParams == len(params) {
			// Converts given parameters into a map whose key is reflect.Type and value is reflect.Value.
			var parameterMap = make(map[reflect.Type]reflect.Value)
			for i := 0; i < len(params); i++ {
				parameterMap[reflect.TypeOf(params[i])] = reflect.ValueOf(params[i])
			}
			var in = make([]reflect.Value, numParams)
			for i := 0; i < numParams; i++ {
				// Checks existence of current parameter by type of invoked method.
				if val, ok := parameterMap[invokedFuncType.In(i)]; ok {
					in[i] = val
				} else {
					return nil, errors.New("Type {" + invokedFuncType.In(i).String() + "} does not exists.")
				}
			}
			res := invokedFunc.Call(in)
			out := convertResponse(res)
			return out, nil
		} else {
			return nil, errors.New("Number of invoked method's parameter {" + strconv.Itoa(numParams) +
				"} does not match given parameters' {" + strconv.Itoa(len(params)) + "}.")
		}
	} else {
		return nil, errors.New("Method Invoked {" + function + "} does not exists. ")
	}
}

// Converts the responses of invoked method represented by slice of
// reflect.Value into slice of interface{}.
func convertResponse(res []reflect.Value) []interface{} {
	var out = make([]interface{}, len(res))
	for i := 0; i < len(out); i++ {
		out[i] = res[i].Interface()
	}
	return out
}

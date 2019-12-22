package core

type (

	// Argument is key value struct that should be passed in a service call
	Argument struct {
		Key   string
		Value interface{}
		// Func returns a dynamic Value instead of Argument.Values
		Func func() interface{}
	}
)

func (a *Argument) GetKey() string {
	return a.Key
}

func (a *Argument) GetValue() interface{} {
	if a.Func != nil {
		return a.Func()
	}
	return a.Value
}

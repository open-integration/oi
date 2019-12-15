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

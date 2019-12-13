package openc

type (

	// Argument is key value struct that should be passed in a service call
	Argument struct {
		Key   string
		Value string
		// Func returns a dynamic Value instead of Argument.Values
		Func func() string
	}
)

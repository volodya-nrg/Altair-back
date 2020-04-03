package response

type Result struct {
	Status int
	Err    error
	Data   interface{}
}

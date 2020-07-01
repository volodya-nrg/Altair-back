package response

// Result - структура ответа, общий результат
type Result struct {
	Status int
	Err    error
	Data   interface{}
}

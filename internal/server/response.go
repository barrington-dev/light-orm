package server

type JSONPayload[T any] struct {
	Status int    `json:"status,omitempty"`
	Data   T      `json:"data,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func JSONSuccess[T any](responseData JSONPayload[T]) JSONPayload[T] {
	return responseData
}

func JSONError[T any](responseData []JSONPayload[T]) map[string][]JSONPayload[T] {
	return map[string][]JSONPayload[T]{
		"errors": responseData,
	}
}

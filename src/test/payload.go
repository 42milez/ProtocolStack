package test

var PayloadBuilder payloadBuilder

type payloadBuilder struct{}

func (payloadBuilder) Default() []byte {
	return []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
}

func init() {
	PayloadBuilder = payloadBuilder{}
}

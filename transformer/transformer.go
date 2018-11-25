package transformer

type Message struct {
	Topic   string
	Payload []byte
}

type Transformer interface {
	Transform(from *Message) (*Message, error)
}

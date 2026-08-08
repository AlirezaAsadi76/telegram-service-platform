package message

type Message struct {
	Text string
}

func New(text string) Message {
	return Message{
		Text: text,
	}
}

package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct { //более общая сущность чем Update, работает и с другими мессенджерами
	Type Type
	Text string
	Meta interface{} //нужен для полей chatID и username которые могут отсутствовать в других мессенджерах
}

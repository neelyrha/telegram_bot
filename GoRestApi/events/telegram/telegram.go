package telegram

import (
	"GoRestApi/clients/telegram"
	"GoRestApi/events"
	"GoRestApi/lib/e"
	"GoRestApi/storage"
	"errors"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct { //прописываем его в этом файле, тк это имеет отношение только к телеграмму
	ChatID   int
	Username string
}

var (
	ErrUnknownEvent    = errors.New("unknown event type")
	ErrNoUpdates       = errors.New("no updates")
	ErrUnknownMetaType = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit) //получаем апдейты
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 { //если список апдейтов пуст то сразу заканчиваем работу функции
		return nil, e.Wrap("no updates", ErrNoUpdates)
	}

	res := make([]events.Event, 0, len(updates)) //готовим переменную для результата,
	// заранее алоцируя память под нее

	for _, u := range updates { //перебираем все апдейты и преобразуем их в тип ивент
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1 //обновляем параметр офсет чтобы в следующий раз
	// получить следующую пачку изменений

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEvent)
	}
}
func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message(1)", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process message(2)", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta) //type assertion. Если не будет мета, то в ok вернется false
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(update telegram.Update) events.Event {
	updType := fetchType(update)
	res := events.Event{
		Type: fetchType(update),
		Text: fetchText(update),
	}

	if updType == events.Message {
		res.Meta = Meta{
			update.Message.Chat.ID,
			update.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}

package telegram

import (
	"GoRestApi/lib/e"
	"GoRestApi/storage"
	"errors"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	// add page: http://...
	// rnd page: /rnd
	// help: /help
	//start: /start: начинается при добавлении бота. в ответ мы присылаем hi + help

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}

	return nil
}

func (p *Processor) savePage(chatID int, pageUrl string, username string) (err error) { //именованные ошибки нужны
	//потому будет несколько мест с их возвратом
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }() //значение err попадет туда перед ретёрном

	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}

	isExist, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}

	if isExist {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil { //здесь можно сделать замыкание
		return err
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.Wrap("can't do command: send random page", err) }()

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool { //оставляем название Add на тот случай если
	// в будущем мы расширим функционал
	return isURL(text)
}

func isURL(text string) bool { //этот метод будет работать только для ссылок с протоколом,
	// без него (ya.ru) работать не будет
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}

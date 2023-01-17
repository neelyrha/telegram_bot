package telegram

import (
	"GoRestApi/lib/e"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

//получение апдейтов, новых сообщений, отправка сообщений пользователям

type Client struct {
	host     string //host api-сервиса телеграма
	basePath string //префикс с которого начинаются все запросы, например //tg-bot.com/bot<token>
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q) //do request <- getUpdates - метод из документации
	if err != nil {
		return nil, err
	}

	var res UpdatesResponce

	if err := json.Unmarshal(data, &res); err != nil {
		//здесь обязательно ссылка на значение иначe Анмаршал
		//не сможет ничего добавить
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error { //неэкспортируемые методы лучше помещать ниже экспортируемых
	q := url.Values{} //пересмотреть
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("Can't send message", err)
	}
	return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   path.Join(c.basePath, method), //мы используем метод Join
		//чтобы решить проблему с повторяющимся '/'
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		//errors.Is
		//errors.As ПОЧИТАТЬ!
		return nil, err
	}

	req.URL.RawQuery = query.Encode() //приведение запроса к тому виду,
	// чтобы его можно было отправлять на сервер
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

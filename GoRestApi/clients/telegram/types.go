package telegram

type UpdatesResponce struct { //то что мы получаем в ответе
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct { //то что мы отправляем в запросе. Это понятие только телеграмма
	ID      int              `json:"update_id"` //update id - из документации телеграм бота
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

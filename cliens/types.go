package client

type UpdateRespouns struct {
	Ok     bool      `json:"ok"`
	Result []Updates `json:"result"`
}

type Updates struct {
	ID      int    `json:"Update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
Text string `json:"text"`

From From `json:"from"`

Chat Chat `json:"caht"`

}

type From struct {

	UserName string `json:"username"`

}

type Chat struct {

	ID int `json:"id"`

}
package websocket

func (hub *Hub) RegisterClientToHub(client *Client) {
	hub.register <- client
}

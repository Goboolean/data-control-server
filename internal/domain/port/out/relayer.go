package outport


type RelayerPort interface {
	SubscribeWebsocket(string) error
	UnsubscribeWebsocket(string) error
}
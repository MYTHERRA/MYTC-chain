package types

const (
	EventTypeRegisterRelay   = "register_relay"
	EventTypeUnregisterRelay = "unregister_relay"
	EventTypeHeartbeat       = "relay_heartbeat"

	AttributeKeyOperator = "operator"
	AttributeKeyWssUrl   = "wss_url"
	AttributeKeyVersion  = "version"
)

package call

type Call struct {
	ChannelID           string
	PeerID              string
	Caller              string
	CallerName          string
	Called              string
	BillingID           string
	ActiveComponentName string
	Data                map[string]interface{}
}

type CallManager struct {
	Calls []*Call
}

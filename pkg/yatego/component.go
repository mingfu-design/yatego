package yatego

// InstallDef Defines one message handler
type InstallDef struct {
	Priority    int
	FilterName  string
	FilterValue string
}

// Consts for results of callback
const (
	ResTransfer = "transfer"
	ResStop     = "stop"
	ResStay     = "stay"
	ResEnter    = "enter"
)

// CallbackResult holds return result from callback func
type CallbackResult struct {
	result            string
	transferComponent string
}

// NewCallbackResult generate CB result
func NewCallbackResult(result string, transferComponent string) *CallbackResult {
	return &CallbackResult{
		result:            result,
		transferComponent: transferComponent,
	}
}

// Callback is the type of msg handler function
type Callback func(call *Call, message *Message) *CallbackResult

// Component is the contract for a object to be a component
type Component interface {
	Enter(call *Call)
	Name() string
	MessagesToWatch() []string
	MessagesToInstall() map[string]InstallDef
	Callback(messageName string) Callback
	Listen(messageName string, callback Callback)
	OnEnter(callback Callback)
	Config(key string) map[string]string
	Logger() Logger
}

// ComponentFactory is a factory method type to build a component
type ComponentFactory func(class string, name string, config map[string]interface{}) Component

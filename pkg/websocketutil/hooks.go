package websocketutil

type HookCollection struct {
	hooks map[Event][]EventCallback
}

func NewHookCollection() (*HookCollection, error) {
	hc := &HookCollection{
		hooks: make(map[Event][]EventCallback),
	}

	return hc, nil
}

func (hc *HookCollection) Register(hookName Event, callback EventCallback) {
	hc.hooks[hookName] = append(hc.hooks[hookName], callback)
}

func (hc *HookCollection) Call(hookName Event, message ClientMessage) {
	eventCallbacks, exists := hc.hooks[hookName]
	if !exists {
		return
	}

	for _, fn := range eventCallbacks {
		fn(message)
	}
}

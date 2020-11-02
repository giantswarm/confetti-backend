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

// Register adds a event listener callback for a specific event type.
func (hc *HookCollection) Register(hookName Event, callback EventCallback) {
	hc.hooks[hookName] = append(hc.hooks[hookName], callback)
}

// Call will execute an event listener for a specific event type.s
func (hc *HookCollection) Call(hookName Event, message ClientMessage) {
	eventCallbacks, exists := hc.hooks[hookName]
	if !exists {
		return
	}

	for _, fn := range eventCallbacks {
		fn(message)
	}
}

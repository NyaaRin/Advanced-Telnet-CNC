package lufx

func NewExecutor(values map[string]any, startTag, endTag string) (*Executor, error) {
	executor := &Executor{
		commands: make(map[string]FxFunction),
		startTag: startTag,
		endTag:   endTag,
	}

	err := executor.register(values)
	if err != nil {
		return nil, err
	}

	return executor, nil
}

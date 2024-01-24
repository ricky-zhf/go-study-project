package local_cache

type Handler interface {
	Handle(sid []string, uuid string) error
}

type HandlerA struct{}

func (h *HandlerA) Handle(sid []string, uuid string) error {
	// 处理sid为a,b,c的情况
	return nil
}

type HandlerB struct{}

func (h *HandlerB) Handle(sid []string, uuid string) error {
	// 处理sid为d,e的情况
	return nil
}

func Post(sid []string, uuid string) error {
	handlers := map[string]Handler{
		"a": &HandlerA{},
		"b": &HandlerA{},
		"c": &HandlerA{},
		"d": &HandlerB{},
		"e": &HandlerB{},
	}

	groups := make(map[Handler][]string)
	for _, s := range sid {
		if handler, ok := handlers[s]; ok {
			groups[handler] = append(groups[handler], s)
		}
	}

	for handler, sids := range groups {
		if err := handler.Handle(sids, uuid); err != nil {
			return err
		}
	}

	return nil
}

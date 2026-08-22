package dispatcher

type MessageDispatcher struct {
	handlers []ConversationHandler
}

func New(handlers ...ConversationHandler) *MessageDispatcher {
	return &MessageDispatcher{
		handlers: handlers,
	}
}

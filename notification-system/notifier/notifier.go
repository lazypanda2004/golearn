package notifier

type Notification struct {
	UserID    string
	Type      string
	Recipient string
	Message   string
}

type Notifier interface {
	Notify(notification Notification) error
}

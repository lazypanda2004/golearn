package notifier

import "log"

type EmailNotifier struct{}

func (e *EmailNotifier) Notify(n Notification) error {
	log.Printf("[EMAIL] To: %s | Msg: %s\n", n.Recipient, n.Message)
	return nil
}

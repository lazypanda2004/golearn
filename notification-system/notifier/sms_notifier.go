package notifier

import "log"

type SMSNotifier struct{}

func (s *SMSNotifier) Notify(n Notification) error {
	log.Printf("[SMS] To: %s | Msg: %s\n", n.Recipient, n.Message)
	return nil
}

package notification

import "crypto/rsa"

// Notification must be implemented on messages passed through the notifier
type Notification interface {
	// Desitnation of the notification
	Destination() string
	// Body of the notification
	Body() []byte
}

// Message contians the notification to send and a retry counter for delayed retries
type Message struct {
	retry int
	note  Notification
}

// Config options for the notifier.
type Config struct {
	// MaxWorkers controls how many go routines should be dedicated to sending notifications
	MaxWorkers int
	// MaxDepth controls how many pending sends/results can be awaiting processing
	MaxDepth int
}

type Notifier interface {
	Start()
	Stop()
	Stopped() bool
	Send(note Notification)
}

func New(config *Config, key *rsa.PrivateKey) Notifier {
	return newBase(config, key)
}

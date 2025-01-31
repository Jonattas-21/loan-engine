package interfaces

import ()

type Queue interface {
	CreateQueue(name string) error
	PublishMessage(queueName string, bodyJson string)
}
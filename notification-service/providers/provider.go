package providers

import (
	"log"
)

type Provider interface {
	Send(email string, message string) error
}

type MockNotificationProvider struct{}

func NewMockNotificationProvider() *MockNotificationProvider {
	return &MockNotificationProvider{}
}

func (p *MockNotificationProvider) Send(email string, message string) error {
	log.Printf("[PROVIDER SEND] ──> Кому: %s | Сообщение: %s", email, message)
	return nil
}

package webhook

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/xmidt-org/urlegit"
)

type MockValidator struct {
	mock.Mock
	ValidateOneEventFunc    func() error
	ValidateEventRegexFunc  func() error
	ValidateDeviceIdFunc    func() error
	ValidateUntilFunc       func(time.Duration, time.Duration, func() time.Time) error
	ValidateNoUntilFunc     func() error
	ValidateDurationFunc    func(time.Duration) error
	ValidateFailureURLFunc  func(*urlegit.Checker) error
	ValidateReceiverURLFunc func(*urlegit.Checker) error
	ValidateAltURLFunc      func(*urlegit.Checker) error
	SetNow                  func(func() time.Time)
}

func (m *MockValidator) ValidateOneEvent() error {
	args := m.Called()
	return args.Error(1)
}

func (m *MockValidator) ValidateEventRegex() error {
	args := m.Called()
	return args.Error(1)
}

func (m *MockValidator) ValidateDeviceId() error {
	args := m.Called()
	return args.Error(1)
}

func (m *MockValidator) ValidateUntil(duration time.Duration, maxTTL time.Duration, now func() time.Time) error {
	args := m.Called()
	return args.Error(1)
}

func (m *MockValidator) ValidateNoUntil() error {
	args := m.Called()
	return args.Error(1)
}

func (m *MockValidator) ValidateDuration(ttl time.Duration) error {
	args := m.Called()
	return args.Error(1)
}

func (m *MockValidator) ValidateFailureURL(checker *urlegit.Checker) error {
	args := m.Called()
	return args.Error(1)
}

func (m *MockValidator) ValidateReceiverURL(checker *urlegit.Checker) error {
	args := m.Called()
	return args.Error(1)
}

func (m *MockValidator) ValidateAltURL(checker *urlegit.Checker) error {
	args := m.Called()
	return args.Error(1)
}

func (m *MockValidator) SetNowFunc(now func() time.Time) {
	args := m.Called()
	fmt.Print(args...)
}

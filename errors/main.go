package errors

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Exchange Related errors
var (
	ErrExchangeAlreadyExists = Error{
		Code:    "E001",
		Message: "Exchange already exists"}
	ErrExchangeDoesnotExist = Error{Code: "E002", Message: "Exchange doesn't exist"}
)

// Queue related errors
var (
	ErrQueueAlreadyExists = Error{Code: "Q001",
		Message: "Queue already exists"}
	ErrQueueDoesnotExist = Error{Code: "Q002", Message: "Queue doesn't exist"}
)

// Subscription related errors

var (
	ErrSubscriptionDoesnotExist = Error{Code: "S001", Message: "Subscription doesnot exist"}
)

func Handle(err error) error {
	logrus.Error(err)
	return err
}

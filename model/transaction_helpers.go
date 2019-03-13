package model

import (
	"errors"
)

type transactionFunc func(TransactionHelper *Helper)
func Transaction(tf transactionFunc, attempts uint) {
	if attempts <= 0 {
		attempts = 1
	}
	var currentAttempt uint
	currentAttempt = 1
	h := Helper{}
	h.SetDB(DB().Begin())
	defer func(_h *Helper) {
		if err := recover(); err != nil {
			var __err error
			if _err, ok := err.(error); ok {
				__err = _err
			} else {
				__err = errors.New(err.(string)) //@todo err.(string) may be down when `panic(123)`
			}
			handleTransactionException(_h, tf, __err, currentAttempt, attempts)
		}
	}(&h)
	tf(&h)
	h.DB().Commit()
}
func handleTransactionException(_h *Helper, f transactionFunc, err error, currentAttempt uint, maxAttempts uint) {
	_h.DB().Rollback()
	if currentAttempt < maxAttempts {
		Transaction(f, maxAttempts-currentAttempt)
	}

	panic(err)
}
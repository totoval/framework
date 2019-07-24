package helper

import (
	"errors"
	"fmt"

	"github.com/totoval/framework/database"
)

type transactionFunc func(TransactionHelper *Helper)

func Transaction(tf transactionFunc, attempts uint) {
	if attempts <= 0 {
		attempts = 1
	}
	var currentAttempt uint
	currentAttempt = 1
	h := Helper{}
	h.SetDB(database.DB().Begin())
	defer func(_h *Helper) {
		if err := recover(); err != nil {
			var __err error
			if _err, ok := err.(error); ok {
				__err = _err
			} else {
				__err = errors.New(fmt.Sprint(err))
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

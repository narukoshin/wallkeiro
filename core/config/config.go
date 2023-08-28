package config

import "wallkeiro/core/errors"

const ExpensesFile string = "user/expenses.json"

var MinimalBalanceAfterExpenses float64 = 150

func SetLevel(level int) error {
	if level > 3 {
		return errors.ErrLevelTooHigh
	}
	switch level {
	case 1:
		MinimalBalanceAfterExpenses = 150
	case 2:
		MinimalBalanceAfterExpenses = 130
	case 3:
		MinimalBalanceAfterExpenses = 110
	}
	return nil
}
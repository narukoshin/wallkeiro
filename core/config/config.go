package config

import "wallkeiro/core/errors"

const ExpensesFile string = "user/expenses.json"

var MinimalBalanceAfterExpenses float64 = 150

func SetLevel(level int) error {
	if level > 4 {
		return errors.ErrLevelTooHigh
	}
	switch level {
	case 1:
		MinimalBalanceAfterExpenses = 150
	case 2:
		MinimalBalanceAfterExpenses = 130
	case 3:
		MinimalBalanceAfterExpenses = 110
	case 4:
		MinimalBalanceAfterExpenses = 50
	}
	return nil
}
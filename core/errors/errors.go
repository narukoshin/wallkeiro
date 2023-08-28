package errors

import "errors"

var ErrSalaryRequired 			= errors.New("Salary is required to calculate the expenses")
var	ErrExpensesMoreThanSalary 	= errors.New("Nothing to save, expenses are more than salary")
var ErrWithdrawnAmountTooLow 	= errors.New("Sorry, for now it seems that your salary is too small to make additional savings.")
var ErrLevelTooHigh 			= errors.New("You can't set level higher than 3.")
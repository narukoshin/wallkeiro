package core

import (
	"wallkeiro/core/config"
	"wallkeiro/core/errors"
	"wallkeiro/core/expenses"

	"flag"
	"fmt"
)

var salary float64
var level int

func Start() error {
	// Initialize flags
	flag.Float64Var(&salary, "salary", 0.0, "Salary")
	flag.IntVar(&level, "level", 1, "Level")
	flag.Parse()
	
	// Parsing commands
	if command := command_parser(); command {
		return nil
	}

	// Setting saving level
	// There are 4 levels, 4 is the strongest in saving matter.
	err := config.SetLevel(level)
	if err != nil {
		return err
	}
	if salary == 0 {
		return errors.ErrSalaryRequired
	}
	// Calculating the savings.
	saving, err := expenses.CalculateSavings(salary)
	if err != nil {
		return err	
	}
	fmt.Println(saving)
	return nil
}

func command_parser() bool {
	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
			case "expenses":
				expenses.ShowExpenses()
			case "add":
				addCmd := flag.NewFlagSet("add", flag.PanicOnError)
				name := addCmd.String("name", "", "name")
				amount := addCmd.Float64("amount", 0.0, "Amount")
				err := addCmd.Parse(args[1:])
				if err != nil {
					panic(err)
				}
				err = expenses.Add(*name, *amount)
				if err != nil {
					panic(err)
				}
			case "remove":
				removeCmd := flag.NewFlagSet("remove", flag.PanicOnError)
                name := removeCmd.String("name", "", "name")
                err := removeCmd.Parse(args[1:])
                if err!= nil {
                    panic(err)
                }
                err = expenses.Remove(*name)
                if err!= nil {
                    panic(err)
                }
            default:
		}
		return true
	}
	return false
}
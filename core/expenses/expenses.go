package expenses

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"wallkeiro/core/config"
	"wallkeiro/core/errors"
)

// Monthly expenses
type ExpensesStuct struct {
	Name string `json:"name"`
	Amount float64 `json:"amount"`
}

func ReadExpenses() ([]ExpensesStuct, error) {
	// Checking if the file exists
	contents, err := os.ReadFile(config.ExpensesFile)
	if err != nil {
		return nil, err
	}

	var expenses []ExpensesStuct
	json.Unmarshal(contents, &expenses)
	return expenses, nil

}

func WriteExpenses (expenses []ExpensesStuct) error {
	data, err := json.Marshal(expenses)	
	if err != nil {
		return err
	}
	err = os.WriteFile(config.ExpensesFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Add(name string, amount float64) error {
	expenses, err := ReadExpenses()
	if err != nil {
		return err
	}

	e := ExpensesStuct {
		Name: name,
		Amount: amount,
	}

	expenses = append(expenses, e)

	err = WriteExpenses(expenses)
	if err != nil {
		return err
	}

	return nil
}

func Remove(name string) error {
	expenses, err := ReadExpenses()
    if err!= nil {
        return err
    }

    for i, e := range expenses {
		if strings.Contains(strings.ToLower(e.Name), strings.ToLower(name)) {
            expenses = append(expenses[:i], expenses[i+1:]...)
            break
        }
	}

    err = WriteExpenses(expenses)
    if err!= nil {
        return err
    }

    return nil
}

func Total() (float64, error){
	expenses, err := ReadExpenses()
    if err!= nil {
        return 0, err
    }

    var total float64
    for _, e := range expenses {
        total += e.Amount
    }

    return total, nil
}

func CalculateSavings(salary float64) (float64, error) {
	expensesTotal, err := Total()
	if err != nil {
		return 0.0, err
	}
	if expensesTotal > salary {
		return 0.0, errors.ErrExpensesMoreThanSalary
	}

	bills := expensesTotal
	desiredFinalBalance := config.MinimalBalanceAfterExpenses 

	remainingAmount := salary - bills - desiredFinalBalance
	withdrawnAmount := math.Max(5*math.Floor(remainingAmount/5), 0)

	if withdrawnAmount <= 10 {
		return 0.0, errors.ErrWithdrawnAmountTooLow
	}
	return withdrawnAmount, nil
}

func ShowExpenses() error {
	expenses, err := ReadExpenses()
    if err!= nil {
        return err
    }
	columns := []string{"Name", "Price"}
	note := "Note: This table displays various products and their prices for your convenience.\nFeel free to browse!"
	printFlexibleTable(note, columns, expenses)
	return nil
}

func printFlexibleTable(note string, columns []string, rows []ExpensesStuct) {
	// Find the maximum width for each column
	colWidths := make([]int, len(columns))
	for colIdx, colName := range columns {
		colWidths[colIdx] = len(colName)
	}
	total := 0.0
	for _, row := range rows {
		// Calculating the dynamic size of the name column
		if len(row.Name) > colWidths[0] {
			colWidths[0] = len(row.Name)
		}
		
		// Setting default column size
		if colWidths[0] < 20 {
			colWidths[0] = 20
		}

		amountStr := fmt.Sprintf("$%.2f", row.Amount)
		// Calculating the dynamic size of the amount column
		if len(amountStr) > colWidths[1] {
			colWidths[1] = len(amountStr)
		}
		// Calculating total expenses
		total += row.Amount
	}

	// Split the multi-line note into lines
	noteLines := strings.Split(note, "\n")

	// Find the maximum width among the note lines
	noteWidth := 0
	for _, line := range noteLines {
		if len(line) > noteWidth {
			noteWidth = len(line)
		}
	}

	// Print the note lines
	fmt.Println()
	for _, line := range noteLines {
		fmt.Printf("%-*s\n", noteWidth, line)
	}
	fmt.Println()

	// Print the table header
	fmt.Print("+")
	for _, width := range colWidths {
		fmt.Print(strings.Repeat("-", width+2), "+")
	}
	fmt.Println()
	for colIdx, colName := range columns {
		fmt.Printf("| %-*s ", colWidths[colIdx], colName)
	}
	fmt.Println("|")
	fmt.Print("+")
	for _, width := range colWidths {
		fmt.Print(strings.Repeat("-", width+2), "+")
	}
	fmt.Println()

	// Print the table rows
	for _, row := range rows {
		fmt.Printf("| %-*s | %-*s |\n", colWidths[0], row.Name, colWidths[1], fmt.Sprintf("%.2f€", row.Amount))
	}

	// Print the total row
	fmt.Print("+")
	for _, width := range colWidths {
		fmt.Print(strings.Repeat("-", width+2), "+")
	}
	fmt.Println()
	fmt.Printf("| %-*s | %-*s |\n", colWidths[0], "Total", colWidths[1], fmt.Sprintf("%.2f€", total))
	fmt.Print("+")
	for _, width := range colWidths {
		fmt.Print(strings.Repeat("-", width+2), "+")
	}
	fmt.Println()
}
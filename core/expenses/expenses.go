package expenses

import (
	"fmt"
	"math"
	"strings"
	"wallkeiro/core/config"
	"wallkeiro/core/errors"

	"github.com/manifoldco/promptui"
)

func Show(ProfileData config.ProfileData) {
	columns := []string{"Name", "Amount"}
	note := "Note: This table displays various products and their prices for your convenience.\nFeel free to browse!"
	printFlexibleTable(note, columns, ProfileData.Expenses)
}

// Add a new expense to the given ProfileData. It takes the name of the expense
// and its amount as arguments, appends the expense to the ProfileData's list
// of expenses and returns the updated ProfileData.
func Add(ProfileData config.ProfileData, name string, amount float64) config.ProfileData {
	ProfileData.Expenses = append(ProfileData.Expenses, config.ExpensesStuct{Name: name, Amount: amount})
	return ProfileData
}

// Edit allows the user to edit an existing expense in the given ProfileData.
// It presents a prompt to select the expense to edit, and then presents
// a prompt to select the action to take: change name, change value, or delete expense.
// After the user selects an action and provides any required information, the function
// updates the ProfileData accordingly and returns the updated ProfileData.
func Edit(ProfileData config.ProfileData) config.ProfileData {
	var expenseNames []string
	for _, expense := range ProfileData.Expenses {
		expenseNames = append(expenseNames, expense.Name)
	}
	prompt := promptui.Select{
		Label: "Select Expense to Edit",
		Items: expenseNames,
		Searcher: func(input string, index int) bool {
			expense := expenseNames[index]
			name := strings.Replace(strings.ToLower(expense), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)
			return strings.Contains(name, input)
		},
	}
	_, expenseSelector, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	var selectedExpense config.ExpensesStuct
	for _, expense := range ProfileData.Expenses {
		if expense.Name == expenseSelector {
			selectedExpense = expense
			break
		}
	}

	// actions-change name, change value, delete
	actionPrompt := promptui.Select{
		Label: "Select Action",
		Items: []string{"Change Name", "Change Amount", "Delete Expense"},
	}
	_, actionSelector, err := actionPrompt.Run()
	if err != nil {
		panic(err)
	}
	switch actionSelector {
	case "Change Name":
		namePrompt := promptui.Prompt{
			Label: "Enter new name",
			Default: selectedExpense.Name,
		}
		newName, err := namePrompt.Run()
		if err != nil {
			panic(err)
		}
		for i, expense := range ProfileData.Expenses {
			if expense.Name == selectedExpense.Name {
				ProfileData.Expenses[i].Name = newName
				break
			}
		}
	case "Change Amount":
		amountPrompt := promptui.Prompt{
			Label: "Enter new amount",
			Default: fmt.Sprintf("$%.2f", selectedExpense.Amount),
		}
		newAmountStr, err := amountPrompt.Run()
		if err != nil {
			panic(err)
		}
		var newAmount float64
		_, err = fmt.Sscanf(newAmountStr, "%f", &newAmount)
		if err != nil {
			panic(err)
		}
		for i, expense := range ProfileData.Expenses {
			if expense.Name == selectedExpense.Name {
				ProfileData.Expenses[i].Amount = newAmount
				break
			}
		}
	case "Delete Expense":
		for i, expense := range ProfileData.Expenses {
			if expense.Name == selectedExpense.Name {
				ProfileData.Expenses = append(ProfileData.Expenses[:i], ProfileData.Expenses[i+1:]...)
				break
			}
		}
	}
	return ProfileData
}

func Calculate(ProfileData config.ProfileData) {
	config.SetLevel(ProfileData.Config.SavingLevel)
	salary := ProfileData.Config.Salary
	expenses := ProfileData.Expenses

	totalExpenses := 0.0
	for _, expense := range expenses {
		totalExpenses += expense.Amount
	}
	if totalExpenses > salary {
		fmt.Errorf(errors.ErrExpensesMoreThanSalary.Error())
	}
	bills := totalExpenses
	desiredFinalBalance := config.MinimalBalanceAfterExpenses
	remainingAmount := salary - bills - desiredFinalBalance
	withdrawnAmount := math.Max(5*math.Floor(remainingAmount/5), 0)
	if withdrawnAmount <= 10 {
		fmt.Errorf(errors.ErrWithdrawnAmountTooLow.Error())
	}
	fmt.Printf("Salary: %.2f€\n", salary)
	fmt.Printf("Total Expenses: %.2f€\n", totalExpenses)
	fmt.Printf("Desired Final Balance: %.2f€\n", desiredFinalBalance)
	fmt.Printf("Remaining Amount after Expenses and Desired Balance: %.2f€\n", remainingAmount)
	fmt.Printf("Suggested Withdrawn Amount: %.2f€\n", withdrawnAmount)
}

func printFlexibleTable(note string, columns []string, rows []config.ExpensesStuct) {
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
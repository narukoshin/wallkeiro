package core

import (
	"wallkeiro/core/config"
	"wallkeiro/core/expenses"
	
	"wallkeiro/core/errors"
	"github.com/manifoldco/promptui"

	"fmt"
)

const CreateNewProfile string = "create new profile"

// Start is the main entry point for the application. It shows a menu of available profiles to the user, and allows them to select a profile to work with. If the user selects the "Create New Profile" option, they are prompted to enter a name for the new profile, and the new profile is created. After selecting or creating a profile, the user is shown a menu of available actions to take on their profile, and can select an action to take. If the user selects the "Go Back" option, they are returned to the main menu.
func Start() error {
	var profiles []string
	var selectedProfile string
	profiles = config.GetProfiles()
	prompt := promptui.Select{
		Label: "Select Profile",
		Items: append(profiles, CreateNewProfile),
	}
	_, profileSelector, err := prompt.Run()
	if err != nil {
		return err
	}

	if profileSelector == CreateNewProfile {
		prompt := promptui.Prompt{
			Label: "Enter New Profile Name",
		}
		profileName, err := prompt.Run()
		if err != nil {
			return err
		}
		err = config.CreateNewProfile(profileName)
		if err != nil {
			return err
		}
		selectedProfile = profileName
		fmt.Printf("Profile %s created successfully.\n", profileName)
	} else {
		selectedProfile = profileSelector
	}
	fmt.Printf("Profile %s selected.\n", selectedProfile)
	prompt = promptui.Select{
		Label: "Select Action",
		Items: []string{"Calculate Savings", "Edit Saving Level", "Edit Salary", "Show Expenses", "Add Expense", "Edit Expenses", "Edit Profile Name", "Delete Profile"},
	}
	ActionMenu:
	_, actionSelector, err := prompt.Run()
	if err != nil {
		return err
	}
	switch actionSelector {
	case "Calculate Savings":
		profileData, err := config.ReadProfile(selectedProfile)
		if err != nil {
			return err
		}
		expenses.Calculate(profileData)
	case "Edit Salary":
		// fixed or hourtly wage
		prompt := promptui.Select{
			Label: "Select Salary Type",
			Items: []config.SalaryType{config.Fixed, config.Hourly},
		}
		_, salaryType, err := prompt.Run()
		if err != nil {
			return err
		}
		var salaryValuePrompt promptui.Prompt
		if salaryType == config.Fixed.String() {
			salaryValuePrompt = promptui.Prompt{
				Label: "Enter Fixed Salary Amount",
			}
		} else {
			salaryValuePrompt = promptui.Prompt{
				Label: "Enter Hourly Wage Amount",
			}
		}
		salaryValueStr, err := salaryValuePrompt.Run()
		if err != nil {
			return err
		}
		var salaryValue float64
		_, err = fmt.Sscanf(salaryValueStr, "%f", &salaryValue)
		if err != nil {
			return err
		}
		err = config.SetSalary(selectedProfile, salaryValue, config.SalaryType(salaryType))
	case "Edit Saving Level":
		prompt := promptui.Prompt{
			Label: "Enter New Saving Level (1-4)",
			Default: "1",
			Validate: func(input string) error {
				var level int
				_, err := fmt.Sscanf(input, "%d", &level)
				if err != nil || (level < 1 || level > 4) {
					return errors.ErrLevelTooHigh
				}
				return nil
			},
		}
		levelStr, err := prompt.Run()
		if err != nil {
			return err
		}
		var level int
		_, err = fmt.Sscanf(levelStr, "%d", &level)
		if err != nil {
			return err
		}
		err = config.SetSavingLevel(selectedProfile, level)
		if err != nil {
			return err
		}
	case "Show Expenses":
		profileData, err := config.ReadProfile(selectedProfile)
		if err != nil {
			return err
		}
		expenses.Show(profileData)
	case "Add Expense":
		promptName := promptui.Prompt{
			Label: "Enter Expense Name",
		}
		expenseName, err := promptName.Run()
		if err != nil {
			return err
		}
		promptAmount := promptui.Prompt{
			Label: "Enter Expense Amount",
		}
		expenseAmountStr, err := promptAmount.Run()
		if err != nil {
			return err
		}
		var expenseAmount float64
		_, err = fmt.Sscanf(expenseAmountStr, "%f", &expenseAmount)
		if err != nil {
			return err
		}
		profileData, err := config.ReadProfile(selectedProfile)
		if err != nil {
			return err
		}
		profileData = expenses.Add(profileData, expenseName, expenseAmount)
		err = config.UpdateProfile(selectedProfile, &profileData)
		if err != nil {
			return err
		}
		fmt.Printf("Expense %s of amount %.2f€ added successfully.\n", expenseName, expenseAmount)
	case "Edit Expenses":
		profileData, err := config.ReadProfile(selectedProfile)
		if err != nil {
			return err
		}

		if len(profileData.Expenses) > 0 {
			profileData = expenses.Edit(profileData)
			err = config.UpdateProfile(selectedProfile, &profileData)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("No expenses to edit.")
		}
	case "Edit Profile Name":
		prompt := promptui.Prompt{
			Label: "Enter New Profile Name",
		}
		newProfileName, err := prompt.Run()
		if err != nil {
			return err
		}
		config.RenameProfile(selectedProfile, newProfileName)
		selectedProfile = newProfileName
	case "Delete Profile":
		return config.DeleteProfile(selectedProfile)
	}
	if goBack, err := GoBackPrompt(); err != nil {
		return err
	} else if goBack {
		goto ActionMenu
	}
	return nil
}

// GoBackPrompt asks the user if they want to go back to the main menu.
// It takes no arguments, and returns a boolean indicating whether the user
// wants to go back to the main menu, and an error if there was an error
// running the prompt.
func GoBackPrompt() (bool, error) {
	prompt := promptui.Select{
		Label: "Do you want to go back to the main menu?",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}
	return result == "Yes", nil
}

// CreateFolder creates the profiles folder if it does not exist.
// This function is called when the application starts and is used to
// ensure that the user's profile folder is created before attempting to
// save expenses data.
func CreateFolder() error {
	return config.CreateProfilesFolder()
}
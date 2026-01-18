package config

import (
	"wallkeiro/core/errors"
	"strings"
	"encoding/json"
	"fmt"
	"os"
)

const ExpensesFile string = "user/expenses.json"

const ProfilesFolder string = "profiles"

var MinimalBalanceAfterExpenses float64 = 150

type ProfileData struct {
	Config ConfigStruct `json:"config"`
	Expenses []ExpensesStuct `json:"expenses"`
}

type SalaryType string

func (s SalaryType) String() string {
	return string(s)
}

const (
	Hourly SalaryType = "hourly"
	Fixed SalaryType = "fixed"
)

type ConfigStruct struct {
	Salary     		float64    `json:"salary"`
	SalaryType 		SalaryType `json:"salary_type"`
	SavingLevel     int        `json:"level"`
}

type ExpensesStuct struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

// SetSalary sets the salary and salary type of a profile.
// It takes the name of the profile, the new salary, and the new salary type as arguments,
// unmarshals the JSON file containing the profile's configuration,
// updates the salary and salary type, marshals the updated configuration back to JSON,
// and writes it to the file.
// If there was an error reading the file, unmarshaling the JSON, or writing the file,
// this function returns that error.
func SetSalary(profile string, salary float64, salaryType SalaryType) error {
	var configData ProfileData
	configData, err := ReadProfile(profile)
	if err != nil {
		return err
	}

	configData.Config.Salary = salary
	configData.Config.SalaryType = salaryType

	err = UpdateProfile(profile, &configData)
	if err != nil {
		return err
	}

	return nil
}

// SetSavingLevel sets the saving level of a profile.
// It takes the name of the profile and the new level as arguments,
// unmarshals the JSON file containing the profile's configuration,
// updates the saving level, marshals the updated configuration back to JSON,
// and writes it to the file.
// If there was an error reading the file, unmarshaling the JSON, or writing the file,
// this function returns that error.
func SetSavingLevel(profile string, level int) error {
	filePath := fmt.Sprintf("%s/%s.json", ProfilesFolder, strings.ToLower(profile))
	profileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	var configData ProfileData
	err = json.Unmarshal(profileData, &configData)
	if err != nil {
		return err
	}
	configData.Config.SavingLevel = level
	data, err := json.Marshal(configData)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// CreateProfilesFolder creates the profiles folder if it does not exist.
func CreateProfilesFolder() error {
	if _, err := os.Stat(ProfilesFolder); os.IsNotExist(err) {
		err := os.Mkdir(ProfilesFolder, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetProfiles returns a list of all profiles in the profiles folder.
// Each profile is represented by the name of the JSON file containing
// the profile's configuration, without the ".json" extension.
func GetProfiles() []string {
	var profiles []string
	files, err := os.ReadDir(ProfilesFolder)
	if err != nil {
		return profiles
	}
	for _, file := range files {
		if !file.IsDir() {
			profiles = append(profiles, file.Name()[:len(file.Name())-5])
		}
	}
	return profiles
}

// ReadProfile reads a profile from the file with the given name in the profiles folder.
// It returns a ProfileData struct containing the profile's configuration, and an error if there was an error reading the file or unmarshaling the JSON.
// If the file does not exist, this function returns an error.
func ReadProfile(profileName string) (ProfileData, error) {
	filePath := fmt.Sprintf("%s/%s.json", ProfilesFolder, strings.ToLower(profileName))
	profileData, err := os.ReadFile(filePath)
	if err != nil {
		return ProfileData{}, err
	}
	var configData ProfileData
	err = json.Unmarshal(profileData, &configData)
	if err != nil {
		return ProfileData{}, err
	}
	return configData, nil
}

// UpdateProfile updates the profile with the given name.
// It takes a pointer to a ProfileData struct as an argument, marshals it to JSON,
// and writes it to the file with the given name in the profiles folder.
// If the file already exists, this function will overwrite it.
// If there was an error marshaling the JSON or writing the file, this function
// returns that error.
func UpdateProfile(profileName string, data *ProfileData) error {
	filePath := fmt.Sprintf("%s/%s.json", ProfilesFolder, strings.ToLower(profileName))
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// CreateNewProfile creates a new profile with the given name.
// It creates a new file in the profiles folder with the given name and
// returns an error if the file already exists or if there was an error
// creating the file.
// If the profile is successfully created, it returns nil.
func CreateNewProfile(profileName string) error {
	filePath := fmt.Sprintf("%s/%s.json", ProfilesFolder, strings.ToLower(profileName))
	_, err := os.Create(filePath)
	// initialize with default config
	defaultConfig := ProfileData{
		Config: ConfigStruct{
			Salary:      0,
			SalaryType:  Fixed,
			SavingLevel: 1,
		},
		Expenses: []ExpensesStuct{},
	}
	data, err := json.Marshal(defaultConfig)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return err
}

// DeleteProfile deletes the profile with the given name.
// It removes the corresponding file for the profile from the profiles folder.
// If the profile does not exist, it returns an error.
func DeleteProfile(profileName string) error {
	filePath := fmt.Sprintf("%s/%s.json", ProfilesFolder, strings.ToLower(profileName))
	err := os.Remove(filePath)
	return err
}

// RenameProfile renames a profile from oldName to newName.
// It does not perform any additional validation or checks,
// so it is up to the caller to ensure that the oldName exists and
// that the newName is valid.
func RenameProfile(oldName, newName string) error {
	oldPath := fmt.Sprintf("%s/%s.json", ProfilesFolder, strings.ToLower(oldName))
	newPath := fmt.Sprintf("%s/%s.json", ProfilesFolder, strings.ToLower(newName))
	err := os.Rename(oldPath, newPath)
	return err
}

func SetLevel(level int) error {
	if level > 4 {
		return errors.ErrLevelTooHigh
	}
	switch level {
	case 1:
		MinimalBalanceAfterExpenses = 190
	case 2:
		MinimalBalanceAfterExpenses = 170
	case 3:
		MinimalBalanceAfterExpenses = 150
	case 4:
		MinimalBalanceAfterExpenses = 100
	}
	return nil
}

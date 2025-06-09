package auth

import (
	"fmt"
	"os"
)

type AuthLogoutCommand struct{}

func NewAuthLogoutCommand() *AuthLogoutCommand {
	return &AuthLogoutCommand{}
}

func (c *AuthLogoutCommand) Execute() error {
	filePath := getAuthConfigPath()
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Failed to remove token file: %v", err)
	}
	fmt.Println("Logged out successfully. Token file removed.")
	return nil
} 
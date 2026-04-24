package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"

	"github.com/hhung06/digimap-backend/config"
	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/database"
	postgresrepo "github.com/hhung06/digimap-backend/internal/repository/postgres"
)

var createAdminCmd = &cobra.Command{
	Use:   "create-admin",
	Short: "Create a system admin user",
	RunE:  runCreateAdmin,
}

func init() {
	rootCmd.AddCommand(createAdminCmd)
}

func runCreateAdmin(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	pool, err := database.NewPool(context.Background(), cfg.Database)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("First name: ")
	firstName, _ := reader.ReadString('\n')
	firstName = strings.TrimSpace(firstName)

	fmt.Print("Last name: ")
	lastName, _ := reader.ReadString('\n')
	lastName = strings.TrimSpace(lastName)

	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return fmt.Errorf("read password: %w", err)
	}
	password := strings.TrimSpace(string(passwordBytes))

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	userRepo := postgresrepo.NewUserRepository(pool)
	user := &domain.User{
		Email:         email,
		PasswordHash:  string(hash),
		FirstName:     firstName,
		LastName:      lastName,
		IsActive:      true,
		IsSystemAdmin: true,
	}

	if err := userRepo.Create(context.Background(), user); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	fmt.Printf("System admin created: %s %s <%s>\n", user.FirstName, user.LastName, user.Email)
	return nil
}

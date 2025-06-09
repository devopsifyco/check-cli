package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	"github.com/devopsifyco/check-cli/checks"
	"github.com/devopsifyco/check-cli/checks/code"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	outputFormat string
	apiKey     string
)

// Path to store the backend token
func getAuthConfigPath() string {
	homeDir := ""
	if u, err := user.Current(); err == nil {
		homeDir = u.HomeDir
	} else {
		homeDir, _ = os.UserHomeDir()
	}
	dosDir := filepath.Join(homeDir, ".dos")
	if _, err := os.Stat(dosDir); os.IsNotExist(err) {
		_ = os.MkdirAll(dosDir, 0700)
	}
	return filepath.Join(dosDir, "checkcli.json")
}

type BackendTokenData struct {
	AccessToken string `json:"access_token"`
	Email       string `json:"email,omitempty"`
	FullName    string `json:"full_name,omitempty"`
	UserID      int    `json:"id,omitempty"`
	GoogleID    string `json:"google_id,omitempty"`
}

func saveBackendToken(data *BackendTokenData) error {
	filePath := getAuthConfigPath()
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(data)
}

func loadBackendToken() (*BackendTokenData, error) {
	filePath := getAuthConfigPath()
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var data BackendTokenData
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func main() {
	// Create root command
	rootCmd := &cobra.Command{
		Use:   "check",
		Short: "DevOpsify Check Tool for various system checks",
		Long:  "DevOpsify Check Tool for performing various system checks including version, OS, speed, and SSL certificate checks.",
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "Output format (json, yaml)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "apikey", "", "API key for version checks")

	// Create command registry
	registry := make(map[string]checks.CheckCommand)

	// Register commands with their names
	commands := map[string]checks.CheckCommand{
		"version":   checks.NewVersionCheckCommand(apiKey, "", false, false, false, false),
		"os":        checks.NewOSCheckCommand(),
		"speed":     checks.NewSpeedCheckCommand(),
		"ssl":       checks.NewSSLCheckCommand(),
	}

	// Add commands to root
	for name, checkCmd := range commands {
		// Store command in registry
		registry[name] = checkCmd

		// Create a new cobra command
		cmdName := name // Create a copy of the name for the closure
		cobraCmd := &cobra.Command{
			Use:   name + " [args]",
			Short: name + " check",
			Run: func(cobraCmd *cobra.Command, args []string) {
				// Get the check command from registry
				cmd := registry[cmdName]
				if cmd == nil {
					fmt.Printf("Error: Command %s not found\n", cmdName)
					os.Exit(1)
				}

				// For version command, create a new instance with the current flags
				if cmdName == "version" {
					// Get command-specific flags
					client, _ := cobraCmd.Flags().GetBool("client")
					full, _ := cobraCmd.Flags().GetBool("full")
					history, _ := cobraCmd.Flags().GetBool("history")
					output, _ := cobraCmd.Flags().GetString("output")
					cve, _ := cobraCmd.Flags().GetBool("cve")

					// Use command-specific output format if provided, otherwise use global output format
					if output != "" {
						outputFormat = output
					}

					cmd = checks.NewVersionCheckCommand(apiKey, outputFormat, full, history, client, cve)
					registry[cmdName] = cmd
				}

				// Execute the check command
				result, err := cmd.Execute(args)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}

				// Print the result
				result.Print(outputFormat)
			},
		}

		// Add version-specific flags
		if name == "version" {
			cobraCmd.Flags().BoolP("client", "c", false, "Use local client instead of remote API")
			cobraCmd.Flags().BoolP("full", "f", false, "Show full version information")
			cobraCmd.Flags().BoolP("history", "H", false, "Show version history")
			cobraCmd.Flags().StringP("output", "o", "", "Output format (json, yaml)")
			cobraCmd.Flags().Bool("cve", false, "Include CVEs in the response")
		}

		// Add the command to root
		rootCmd.AddCommand(cobraCmd)
	}

	// --- Add 'code' command with subcommands ---
	codeCmd := &cobra.Command{
		Use:   "code",
		Short: "Code analysis commands (deps, loc)",
		Long:  "Code analysis commands: dependencies and lines of code.",
	}

	// 'code deps' subcommand
	codeDepsCmd := &cobra.Command{
		Use:   "deps [path]",
		Short: "Check project dependencies",
		Run: func(cmd *cobra.Command, args []string) {
			cve, _ := cmd.Flags().GetBool("cve")
			output, _ := cmd.Flags().GetString("output")
			if output != "" {
				outputFormat = output
			}
			cmdObj := code.NewCodeDepsCheckCommand(outputFormat, cve)
			result, err := cmdObj.Execute(args)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			result.Print(outputFormat)
		},
	}
	codeDepsCmd.Flags().Bool("cve", false, "Include CVEs in the response")
	codeDepsCmd.Flags().StringP("output", "o", "", "Output format (json, yaml)")

	// 'code loc' subcommand
	codeLocCmd := &cobra.Command{
		Use:   "loc [path]",
		Short: "Count lines of code",
		Run: func(cmd *cobra.Command, args []string) {
			output, _ := cmd.Flags().GetString("output")
			if output != "" {
				outputFormat = output
			}
			cmdObj := code.NewLocCheckCommand(outputFormat)
			result, err := cmdObj.Execute(args)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			result.Print(outputFormat)
		},
	}
	codeLocCmd.Flags().StringP("output", "o", "", "Output format (json, yaml)")

	// Add subcommands to 'code'
	codeCmd.AddCommand(codeDepsCmd)
	codeCmd.AddCommand(codeLocCmd)

	// Add 'code' to root
	rootCmd.AddCommand(codeCmd)

	// --- Add 'auth' command with subcommands ---
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication commands (login, logout)",
		Long:  "Authentication commands for logging in and out of Google accounts.",
	}

	authLoginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Google account",
		Run: func(cmd *cobra.Command, args []string) {
			clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
			clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
			if clientID == "" || clientSecret == "" {
				fmt.Println("GOOGLE_OAUTH_CLIENT_ID and GOOGLE_OAUTH_CLIENT_SECRET environment variables must be set.")
				return
			}
			// Use a fixed redirect URL for local server
			redirectURL := "http://localhost:8085/auth/google/callback"
			oauthCfg := &oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				RedirectURL:  redirectURL,
				Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
				Endpoint:     google.Endpoint,
			}

			state := generateStateOauthCookie()
			url := oauthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)

			// Start local server to handle callback
			server := &http.Server{Addr: ":8085"}
			http.HandleFunc("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("state") != state {
					fmt.Fprintln(w, "State mismatch. Try again.")
					log.Println("State mismatch.")
					return
				}
				code := r.URL.Query().Get("code")
				token, err := oauthCfg.Exchange(context.Background(), code)
				if err != nil {
					fmt.Fprintln(w, "Failed to exchange token:", err)
					log.Println("Token exchange error:", err)
					return
				}

				// Exchange Google token for backend token
				apiEndpoint := "https://api.opsify.dev/checks" // or configurable
				loginUrl := apiEndpoint + "/user/google/login"
				payload := map[string]string{"token": token.AccessToken}
				jsonPayload, _ := json.Marshal(payload)

				// Get API key from flag or env
				apiKey := os.Getenv("CHECK_API_KEY")
				if apiKey == "" {
					apiKey = "SPK1HgBWcxO5EmLsCSP6aIRNhX6wXMYa"
				}
				if apiKey == "" {
					apiKey = cmd.Flag("apikey").Value.String()
				}

				req, err := http.NewRequest("POST", loginUrl, bytes.NewReader(jsonPayload))
				if err != nil {
					fmt.Fprintln(w, "Failed to create backend request:", err)
					log.Println("Backend request error:", err)
					return
				}
				req.Header.Set("Content-Type", "application/json")
				if apiKey != "" {
					req.Header.Set("apikey", apiKey)
				}
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					fmt.Fprintln(w, "Failed to authenticate with backend:", err)
					log.Println("Backend auth error:", err)
					return
				}
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					fmt.Fprintln(w, "Backend login failed with status:", resp.Status)
					log.Println("Backend login failed with status:", resp.Status)
					return
				}
				var backendResp struct {
					AccessToken string `json:"access_token"`
					TokenType   string `json:"token_type"`
					User struct {
						Email     string `json:"email"`
						FullName  string `json:"full_name"`
						ID        int    `json:"id"`
						GoogleID  string `json:"google_id"`
					} `json:"user"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&backendResp); err != nil {
					fmt.Fprintln(w, "Failed to parse backend response:", err)
					log.Println("Backend response parse error:", err)
					return
				}
				data := BackendTokenData{
					AccessToken: backendResp.AccessToken,
					Email:       backendResp.User.Email,
					FullName:    backendResp.User.FullName,
					UserID:      backendResp.User.ID,
					GoogleID:    backendResp.User.GoogleID,
				}
				if err := saveBackendToken(&data); err != nil {
					fmt.Fprintln(w, "Failed to save token:", err)
					log.Println("Token save error:", err)
					return
				}
				fmt.Fprintf(w, "Login successful! You can close this window.\n")
				fmt.Printf("User info: %s\n", toJsonString(data))
				go func() { time.Sleep(1 * time.Second); server.Shutdown(context.Background()) }()
			})

			// Open browser
			fmt.Println("Opening browser for Google login...")
			openBrowser(url)
			fmt.Println("Waiting for Google login...")
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				log.Println("Server error:", err)
			}
		},
	}

	authLogoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from Google account",
		Run: func(cmd *cobra.Command, args []string) {
			filePath := getAuthConfigPath()
			err := os.Remove(filePath)
			if err != nil && !os.IsNotExist(err) {
				fmt.Printf("Failed to remove token file: %v\n", err)
				return
			}
			fmt.Println("Logged out successfully. Token file removed.")
		},
	}

	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authLogoutCmd)
	rootCmd.AddCommand(authCmd)

	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// Helper functions for OAuth2
func generateStateOauthCookie() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func openBrowser(url string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Start()
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		fmt.Printf("Please open the following URL in your browser:\n%s\n", url)
	}
}

// Helper to pretty-print user info
func toJsonString(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
} 
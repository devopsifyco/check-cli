package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"runtime"
)

var (
	googleOAuthClientID     string
	googleOAuthClientSecret string
	CheckApiKeyDemo         string
)

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

func LoadBackendToken() (*BackendTokenData, error) {
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

func getAuthConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".check_auth.json"
	}
	return home + string(os.PathSeparator) + ".check_auth.json"
}

type AuthLoginCommand struct{}

func NewAuthLoginCommand() *AuthLoginCommand {
	return &AuthLoginCommand{}
}

func (c *AuthLoginCommand) Execute() error {
	clientID := googleOAuthClientID
	clientSecret := googleOAuthClientSecret
	if clientID == "" {
		clientID = os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	}
	if clientSecret == "" {
		clientSecret = os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	}
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("googleOAuthClientID and googleOAuthClientSecret must be set at build time using -ldflags or via environment variables.")
	}
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
		apiEndpoint := "https://api.opsify.dev/checks"
		loginUrl := apiEndpoint + "/user/google/login"
		payload := map[string]string{"token": token.AccessToken}
		jsonPayload, _ := json.Marshal(payload)
		apiKey := os.Getenv("CHECK_API_KEY")
		if apiKey == "" {
			apiKey = os.Getenv("CHECK_API_KEY_DEMO")
		}
		if apiKey == "" {
			apiKey = CheckApiKeyDemo
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
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `Login successful! <a href="#" onclick="window.close()">Click here to close this window</a>
<script>
setTimeout(function() {
    window.close();
}, 3000);
</script>`)
		log.Println("Login successful!")
		// Wait longer before shutting down to ensure the page is fully rendered
		go func() { 
			time.Sleep(4 * time.Second)
			server.Shutdown(context.Background())
		}()
	})
	fmt.Println("Opening browser for Google login...")
	openBrowser(url)
	fmt.Println("Waiting for Google login...")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Println("Server error:", err)
	}
	return nil
}

func generateStateOauthCookie() string {
	return fmt.Sprintf("st%d", time.Now().UnixNano())
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default: // linux, freebsd, openbsd, netbsd
		err = exec.Command("xdg-open", url).Start()
	}
	if err != nil {
		fmt.Printf("Please open the following URL in your browser:\n%s\n", url)
	}
} 
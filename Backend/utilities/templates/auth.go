package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var oauthConfig = &oauth2.Config{
	ClientID:     "YOUR_CLIENT_ID",
	ClientSecret: "YOUR_CLIENT_SECRET",
	RedirectURL:  "http://localhost:8080/callback",
	Scopes:       []string{"https://graph.microsoft.com/User.Read"},
	Endpoint:     microsoft.AzureADEndpoint("common"),
}

var tokenMap = make(map[string]*oauth2.Token)

// Login Handler
func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

// Callback Handler
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}

	client := oauthConfig.Client(context.Background(), token)
	userInfo, err := fetchUserInfo(client)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// Send welcome email
	err = sendEmail(userInfo.Email, userInfo.Name)
	if err != nil {
		log.Println("Failed to send email:", err)
	}

	// Display welcome page
	tmpl, _ := template.New("welcome").Parse(`
		<html>
		<head><title>Welcome</title></head>
		<body>
			<h1>Hello, {{ .Name }} 👋</h1>
			<p>Welcome to Virtual Classroom!</p>
		</body>
		</html>
	`)
	tmpl.Execute(w, userInfo)
}

// Fetch user info from Microsoft Graph API
type User struct {
	Name  string `json:"displayName"`
	Email string `json:"mail"`
}

func fetchUserInfo(client *http.Client) (*User, error) {
	req, _ := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user User
	json.NewDecoder(resp.Body).Decode(&user)
	return &user, nil
}

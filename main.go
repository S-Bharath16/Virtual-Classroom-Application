package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var (
	clientID     = "a40116bf-ccaf-409a-8798-7802d57c1f0f"
	clientSecret = "i-p8Q~9vVH8uP6qWdBRWLmA44erfojbIr3XgPaKJ"
	redirectURL  = "http://localhost:3000/auth/callback"
	providerURL  = "https://login.microsoftonline.com/common/v2.0"

	oauth2Config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     microsoft.AzureADEndpoint("common"),
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "email", "profile"},
	}
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	code := r.URL.Query().Get("code")

	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		log.Printf("Token exchange error: %+v", err) // Log full error details
		http.Error(w, "OAuth exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	const tenantID = "7bc5987c-2e79-42ff-bfec-cf5bbacacea6"

	provider, err := oidc.NewProvider(ctx, "https://login.microsoftonline.com/00f9cda3-075e-44e5-aa0b-aba3add6539f/v2.0")
	if err != nil {
		log.Fatal(err)
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	idToken, err := verifier.Verify(ctx, token.Extra("id_token").(string))
    if err != nil {

        log.Printf("Token verification error: %+v", err) // Log full error details
        http.Error(w, "Token verification failed: "+err.Error(), http.StatusInternalServerError)
        return
    }

	var claims struct {
		Email string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
		return
	}

	// Restrict access to college domain
	if !isAllowedDomain(claims.Email) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	fmt.Fprintf(w, "Welcome, %s!", claims.Email)
}

func isAllowedDomain(email string) bool {
	allowedDomain := "amrita.edu"
	return len(email) > len(allowedDomain) && email[len(email)-len(allowedDomain):] == allowedDomain
}

func main() {
	http.HandleFunc("/auth/login", loginHandler)
	http.HandleFunc("/auth/callback", callbackHandler)

	fmt.Println("Server started at :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

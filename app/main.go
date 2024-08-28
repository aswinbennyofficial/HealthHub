package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	fitbitAuthURL  = "https://www.fitbit.com/oauth2/authorize"
	fitbitTokenURL = "https://api.fitbit.com/oauth2/token"
	clientID       = "23PKSG"
	clientSecret   = "0ba8328d12924c03b89bd0cb23b90e1d"
	redirectURI    = "http://localhost:8080/callback"
	codeVerifier   = "01234567890123456789012345678901234567890123456789"
)

func main() {
	authURL := fmt.Sprintf("%s?client_id=%s&response_type=code&code_challenge=-4cf-Mzo_qg9-uq0F4QwWhRh4AjcAqNx7SbYVsdmyQM&code_challenge_method=S256&scope=activity%%20heartrate%%20location%%20nutrition%%20oxygen_saturation%%20profile%%20respiratory_rate%%20settings%%20sleep%%20social%%20temperature%%20weight", fitbitAuthURL, clientID)

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, authURL, http.StatusSeeOther)
	})

	http.HandleFunc("/callback", handleCallback)

	fmt.Println("Server is running on http://localhost:8080")
	fmt.Println("Visit http://localhost:8080/auth to start the OAuth2 flow")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")
	fmt.Println("Auth Code:", authCode)

	tokenResp, err := exchangeToken(authCode)
	if err != nil {
		fmt.Fprintf(w, "Error exchanging token: %v", err)
		return
	}

	fmt.Fprintf(w, "Access Token: %s\nRefresh Token: %s", tokenResp.AccessToken, tokenResp.RefreshToken)

	// set cookie as jwt token
	http.SetCookie(w, &http.Cookie{
		Name:  "jwt",
		Value: tokenResp.AccessToken,
		Path: "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func exchangeToken(authCode string) (*FitbitTokenResponse, error) {
	reqBody := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {authCode},
		"redirect_uri":  {redirectURI},
		"client_id":     {clientID},
		"code_verifier": {codeVerifier},
	}

	req, err := http.NewRequest("POST", fitbitTokenURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = reqBody.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var tokenResp FitbitTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &tokenResp, nil
}

type FitbitTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}



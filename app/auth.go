package main

import(
	"fmt"
	"net/http"
	"encoding/json"
	"net/url"
	
)


func handleCallback(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")
	fmt.Println("Auth Code:", authCode)

	tokenResp, err := exchangeToken(authCode)
	if err != nil {
		fmt.Fprintf(w, "Error exchanging token: %v", err)
		return
	}

	

	// set cookie as jwt token
	http.SetCookie(w, &http.Cookie{
		Name:  "jwt",
		Value: tokenResp.AccessToken,
		Path: "/",
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
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



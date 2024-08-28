package main

import (
	"fmt"
	"net/http"
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

	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("jwt")
		if err != nil {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}

		// serve dashboard.html from the static directory
		http.ServeFile(w, r, "static/dashboard.html")
	})

	// make a route to get index.html on root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "static/index.html")
	})

	// make a route to get summary data
	http.HandleFunc("/api/summary", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken := cookie.Value
		getSummary(accessToken,w)
	})

	// make a route to get Lifetime stats
	http.HandleFunc("/api/lifetime", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken := cookie.Value
		getLifetimeStats(accessToken,w)
	})


	// route to fetch cardio fitness score for the past 30 days
	http.HandleFunc("/api/cardio", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken := cookie.Value
		getCardioFitnessScore(accessToken,w)
	})


	// get hrv data
	http.HandleFunc("/api/hrv", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken := cookie.Value
		getHRVSummary(accessToken,w)
	})

	fmt.Println("Server is running on http://localhost:8080")
	fmt.Println("Visit http://localhost:8080/auth to start the OAuth2 flow")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}


	

}

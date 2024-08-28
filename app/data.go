package main

import (
	"encoding/json"
	"fmt"
	"time"

	// "log"
	"net/http"
)

func getProfile(accessToken string) {
    profileURL := "https://api.fitbit.com/1/user/-/profile.json"
    req, err := http.NewRequest("GET", profileURL, nil)
    if err != nil {
        fmt.Printf("error creating request: %v\n", err)
        return
    }
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        fmt.Printf("error sending request: %v\n", err)
        return
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        fmt.Printf("unexpected status code: %d\n", resp.StatusCode)
        return
    }
    var profile map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
        fmt.Printf("error decoding response: %v\n", err)
        return
    }
    fmt.Printf("Profile: %v\n", profile)
}

func getSummary(accessToken string, w http.ResponseWriter) {
    summaryURL := "https://api.fitbit.com/1/user/-/activities/date/today.json"
    req, err := http.NewRequest("GET", summaryURL, nil)
    if err != nil {
        http.Error(w, fmt.Sprintf("error creating request: %v", err), http.StatusInternalServerError)
        return
    }
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        http.Error(w, fmt.Sprintf("error sending request: %v", err), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        http.Error(w, fmt.Sprintf("unexpected status code: %d", resp.StatusCode), resp.StatusCode)
        return
    }
    var summary map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
        http.Error(w, fmt.Sprintf("error decoding response: %v", err), http.StatusInternalServerError)
        return
    }

    // Set the response header to JSON
    w.Header().Set("Content-Type", "application/json")
    
    // Encode the summary as JSON and send it as the response
    if err := json.NewEncoder(w).Encode(summary); err != nil {
        http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
        return
    }
}

func getLifetimeStats(accessToken string, w http.ResponseWriter) {
	lifetimeStatsURL := "https://api.fitbit.com/1/user/-/activities.json"
	req, err := http.NewRequest("GET", lifetimeStatsURL, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating request: %v", err), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("error sending request: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("unexpected status code: %d", resp.StatusCode), resp.StatusCode)
		return
	}
	var lifetimeStats map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&lifetimeStats); err != nil {
		http.Error(w, fmt.Sprintf("error decoding response: %v", err), http.StatusInternalServerError)
		return
	}

	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")
	
	// Encode the lifetime stats as JSON and send it as the response
	if err := json.NewEncoder(w).Encode(lifetimeStats); err != nil {
		http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func getCardioFitnessScore(accessToken string, w http.ResponseWriter) {
    // Calculate the date 30 days ago
    thirtyDaysAgo := time.Now().AddDate(0, 0, -15).Format("2006-01-02")
    today := time.Now().Format("2006-01-02")

    cardioFitnessScoreURL := fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/cardio/minutesVeryActive/date/%s/%s.json", thirtyDaysAgo, today)
    
    req, err := http.NewRequest("GET", cardioFitnessScoreURL, nil)
    if err != nil {
        http.Error(w, fmt.Sprintf("error creating request: %v", err), http.StatusInternalServerError)
        return
    }
    
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        http.Error(w, fmt.Sprintf("error sending request: %v", err), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        http.Error(w, fmt.Sprintf("%s", resp.Body), resp.StatusCode)
        return
    }
    
    var cardioFitnessScore map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&cardioFitnessScore); err != nil {
        http.Error(w, fmt.Sprintf("error decoding response: %v", err), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(cardioFitnessScore); err != nil {
        http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
        return
    }
}

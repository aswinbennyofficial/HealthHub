package main

import (
	"fmt"
	"net/http"
	"encoding/json"
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
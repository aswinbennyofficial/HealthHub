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



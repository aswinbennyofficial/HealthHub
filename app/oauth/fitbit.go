package oauth

import (
	"context"

	// "encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/aswinbennyofficial/SuSHi/models"
	"golang.org/x/oauth2"
	"github.com/rs/zerolog/log"
)



func GenerateFitbitAuthURL(config models.Config) string {
    fitbitConfig := config.OAuthConfig.FitBit

    // Prepare the query parameters
    params := url.Values{}
    params.Add("client_id", fitbitConfig.ClientID)
    params.Add("redirect_uri", fitbitConfig.RedirectURL)
    params.Add("scope", strings.Join(fitbitConfig.Scopes, " "))
    params.Add("state", fitbitConfig.State)

    // Construct the final URL
	authURL := "https://www.fitbit.com/oauth2/authorize?client_id=ABC123&response_type=code" +
		"&code_challenge=<code_challenge>&code_challenge_method=S256" +
		"&scope=activity%20heartrate%20location%20nutrition%20oxygen_saturation%20profile" +
		"%20respiratory_rate%20settings%20sleep%20social%20temperature%20weight"

    return authURL
}


func HandleFitbitCallback(config models.Config, code string) (string, string, error){
	fitbitConfig:= config.OAuthConfig.FitBit

	oauth2Config := &oauth2.Config{
		ClientID: fitbitConfig.ClientID,
		ClientSecret: fitbitConfig.ClientSecret,
		RedirectURL: fitbitConfig.RedirectURL,
		Scopes: fitbitConfig.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.fitbit.com/oauth2/authorize",
			TokenURL: "https://api.fitbit.com/oauth2/token",
		},
	}

	token,err := oauth2Config.Exchange(context.Background(),code)
	if err!=nil {
		return "","",fmt.Errorf("failed to exchange code for token: %w", err)
    }


	// client := oauth2Config.Client(context.Background(), token)

	log.Debug().Msg(token.AccessToken)

	return "","",nil



	

}
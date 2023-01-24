package gcloud

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stocks-sync/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type RedirectURL struct{
    State   string      `form:"state"`
    Code    string      `form:"code"`
    Scope   []string    `form:"scope"`
}

func GetGoogleConfig(env utils.Settings)(*oauth2.Config){
    credentialJson, err := os.ReadFile(env.GoogleOAuthCredentials)
    utils.CheckError(err)
    config, err := google.ConfigFromJSON(
        credentialJson, 
        sheets.SpreadsheetsScope,
        sheets.SpreadsheetsReadonlyScope,
    )
    utils.CheckError(err)
    return config
}

func GetTokenFromWeb(ctx context.Context){
    env := ctx.Value("env").(utils.Settings)
    config := GetGoogleConfig(env)
    authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    fmt.Printf("Go to the following link in your browser then type the "+ "authorization code: \n%v\n", authURL)
}

func GetOauthClient(env utils.Settings) (*http.Client, error) {
    config := GetGoogleConfig(env)
    tok, err := tokenFromFile(env.GoogleTokenFile)
    if err != nil{
        return nil, err
    }
    return config.Client(context.Background(), tok), nil
}

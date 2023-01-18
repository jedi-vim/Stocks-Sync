package gcloud

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"stocks-sync/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
        authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
        fmt.Printf("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var authCode string
        if _, err := fmt.Scan(&authCode); err != nil {
                log.Fatalf("Unable to read authorization code %v", err)
        }

        tok, err := config.Exchange(context.TODO(), authCode)
        utils.CheckError(err)
        return tok
}

func GetOauthClient(env utils.Settings) *http.Client {
    credentialJson, err := os.ReadFile(env.GoogleOAuthCredentials)
    utils.CheckError(err)

    config, err := google.ConfigFromJSON(
        credentialJson, 
        sheets.SpreadsheetsScope,
        sheets.SpreadsheetsReadonlyScope,
    )
    utils.CheckError(err)
    tok, err := tokenFromFile(env.GoogleTokenFile)
    if err != nil {
            tok = getTokenFromWeb(config)
            saveToken(env.GoogleTokenFile, tok)
    }
    return config.Client(context.Background(), tok)
}



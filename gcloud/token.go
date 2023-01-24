package gcloud

import (
	"fmt"
	"log"
	"os"

	"encoding/json"

	"golang.org/x/oauth2"
)

type ExpiredToken struct{
    accessToken string
}

func (e *ExpiredToken) Error() string {
    return fmt.Sprintf("Token Expirado. Gere um novo token")
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    if err != nil {
            return nil, err
    }
    defer f.Close()
    tok := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(tok); if err != nil{
        return tok, err
    }
    if !tok.Valid(){
        return tok, &ExpiredToken{
            accessToken: tok.AccessToken,
        }
    }
    return tok, nil
}

// Saves a token to a file path.
func SaveToken(path string, token *oauth2.Token) {
    fmt.Printf("Saving credential file to: %s\n", path)
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
            log.Fatalf("Unable to cache oauth token: %v", err)
    }
    defer f.Close()
    json.NewEncoder(f).Encode(token)
}

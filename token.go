package main 

import (
	"log"
	"os"
	"encoding/json"

	"golang.org/x/oauth2"
)

// Retrieves a token from a local file.
func TokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    if err != nil {
        log.Println(err)
        return nil, err
    }
    defer f.Close()
    tok := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(tok)
    return tok, err
}

// Saves a token to a file path.
func SaveToken(path string, token *oauth2.Token) {
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        log.Fatalf("Unable to cache oauth token: %v", err)
    }
    defer f.Close()
    log.Printf("Salvando novo token:\n%v", token)
    json.NewEncoder(f).Encode(token)
}



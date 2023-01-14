package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func checkError(err error){
    if err != nil{
        log.Fatal(err)
    }
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokFile string) *http.Client {
        // The file token.json stores the user's access and refresh tokens, and is
        // created automatically when the authorization flow completes for the first
        // time.
        tok, err := tokenFromFile(tokFile)
        if err != nil {
                tok = getTokenFromWeb(config)
                saveToken(tokFile, tok)
        }
        return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
        authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
        fmt.Printf("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var authCode string
        if _, err := fmt.Scan(&authCode); err != nil {
                log.Fatalf("Unable to read authorization code %v", err)
        }

        tok, err := config.Exchange(context.TODO(), authCode)
        if err != nil {
                log.Fatalf("Unable to retrieve token from web %v", err)
        }
        return tok
}

func main() {
    environment := Env()
    cxt := context.Background()
    
    credentialJson, err := os.ReadFile(environment.GoogleOAuthCredentials)
    checkError(err)

    config, err := google.ConfigFromJSON(
        credentialJson, 
        sheets.SpreadsheetsScope,
        sheets.SpreadsheetsReadonlyScope,
    )
    checkError(err)
    client := getClient(config, environment.GoogleTokenFile)

    sheetsService, err := sheets.NewService(cxt, option.WithHTTPClient(client))
    checkError(err)
    
    stockCodesRange, err := sheetsService.Spreadsheets.Values.Get(environment.DocumentID, "Acoes!B3:C17").Do()
    checkError(err)
    for idx, stockData := range stockCodesRange.Values{
        bmfCode := fmt.Sprintf("%v", stockData[0])
        stockPrice := GetStockPrice(bmfCode)
        stockCodesRange.Values[idx][1] = stockPrice
        time.Sleep(2 * time.Second)
    }
    response, err := sheetsService.Spreadsheets.Values.Update(
        environment.DocumentID, 
        stockCodesRange.Range, 
        stockCodesRange).ValueInputOption("USER_ENTERED").Do()
    checkError(err)
    fmt.Printf("%d Valores alterados com sucesso", response.UpdatedRows)
}


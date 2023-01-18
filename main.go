package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"stocks-sync/gcloud"
	"stocks-sync/utils"
)

func main() {
    environment := utils.Env()
    client := gcloud.GetOauthClient(environment)

    cxt := context.Background()
    sheetsService, err := sheets.NewService(cxt, option.WithHTTPClient(client))
    utils.CheckError(err)
    
    stockCodesRange, err := sheetsService.Spreadsheets.Values.Get(environment.DocumentID, "Acoes!B3:C17").Do()
    utils.CheckError(err)
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
    utils.CheckError(err)
    fmt.Printf("%d Valores alterados com sucesso", response.UpdatedRows)
}


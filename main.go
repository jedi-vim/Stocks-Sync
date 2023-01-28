package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
    env := Env()
    ctx := context.WithValue(context.Background(), "env", env)

    // service, err := GetDriveService(ctx)
    // CheckError(err)
    // files, err := service.Files.List().Do()
    // CheckError(err)
    // if files == nil{
    //     log.Fatalln("Nenhum arquivo encontrado")
    // }
    // log.Printf("Total de Arquivos: %d\n", len(files.Files))
    // for _, f := range files.Files{
        // log.Printf("%s %s\n", f.Id, f.Name)
    // }

    sheetsService, err := GetSheetsService(ctx)
    CheckError(err)

    stockCodesRange, err := sheetsService.Spreadsheets.Values.Get(env.DocumentID, env.SheetCellRange).Do()
    CheckError(err)
    for idx, stockData := range stockCodesRange.Values{
        bmfCode := fmt.Sprintf("%v", stockData[0])
        stockPrice := GetStockPrice(bmfCode)
        stockCodesRange.Values[idx][1] = stockPrice
        log.Printf("%s atualizada\n", bmfCode)
        time.Sleep(2 * time.Second)
    }
    response, err := sheetsService.Spreadsheets.Values.Update(
        env.DocumentID, 
        stockCodesRange.Range, 
        stockCodesRange).ValueInputOption("USER_ENTERED").Do()
    CheckError(err)
    log.Printf("%d Valores alterados com sucesso", response.UpdatedRows)
}

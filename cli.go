package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alecthomas/kong"
        "github.com/gin-gonic/gin"

	"stocks-sync/gcloud"
	"stocks-sync/utils"
)

type UpdateSheet struct{}

func (u *UpdateSheet) Run(kongCtx *kong.Context)error{
    environment := utils.Env()
    ctx := context.WithValue(context.Background(), "env", environment)
    sheetsService, err := gcloud.GetSheetsService(ctx)
    utils.CheckError(err)
    
    stockCodesRange, err := sheetsService.Spreadsheets.Values.Get(environment.DocumentID, environment.SheetCellRange).Do()
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
    log.Printf("%d Valores alterados com sucesso", response.UpdatedRows)
    return nil
}

type GenerateToken struct{}

func (g *GenerateToken) Run(kongCtx *kong.Context) error{
    environment := utils.Env()
    ctx := context.WithValue(context.Background(), "env", environment)
    gcloud.GetTokenFromWeb(ctx)

    r := gin.Default()
    r.GET("/", func(c *gin.Context) {
        var urlQueryData gcloud.RedirectURL
        if err := c.ShouldBind(&urlQueryData); err != nil{
            c.JSON(400, gin.H{"msg": err})
            return
        }
        config := gcloud.GetGoogleConfig(environment)
        tok, err := config.Exchange(context.TODO(), urlQueryData.Code)
        if err != nil{
            c.JSON(400, gin.H{"msg": err})
            return
        }
        gcloud.SaveToken(environment.GoogleTokenFile, tok)
        c.JSON(200, gin.H{"message": "Token Salvo com sucesso"})
    })
    r.Run()
    return nil
}

type CLI struct{
    Update UpdateSheet   `cmd:"" help:"Atualizar a Planilha de Acoes com precos atuais"`
    Token GenerateToken  `cmd:"generate-token" help:"Gerar novo access token"`
}

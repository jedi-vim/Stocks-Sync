package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func EnvMiddleware() gin.HandlerFunc{
    return func(c *gin.Context){
        c.Set("env", Env())
        c.Next()
    }
}

func UpdateStocks(c *gin.Context) {
        env := c.MustGet("env").(Settings)
        ctx := context.WithValue(context.Background(), "env", env)

        sheetsService, err := GetSheetsService(ctx)
        if err != nil{
            c.JSON(500, gin.H{"message": err})
        }

        stockCodesRange, err := sheetsService.Spreadsheets.Values.Get(env.DocumentID, env.SheetCellRange).Do()
        if err != nil{
            c.JSON(500, gin.H{"message": err})
        }
        for idx, stockData := range stockCodesRange.Values{
            bmfCode := fmt.Sprintf("%v", stockData[0])
            stockPrice := GetStockPrice(bmfCode)
            stockCodesRange.Values[idx][1] = stockPrice
            log.Printf("%s atualizada\n", bmfCode)
            // time.Sleep(1 * time.Second)
        }
        response, err := sheetsService.Spreadsheets.Values.Update(
            env.DocumentID, 
            stockCodesRange.Range, 
            stockCodesRange).ValueInputOption("USER_ENTERED").Do()
        if err != nil{
            c.JSON(500, gin.H{"message": err})
        }
        c.JSON(200, gin.H{"message": fmt.Sprintf("%d Valores alterados com sucesso", response.UpdatedRows)})
}

func GrantOauthPermission(c * gin.Context){
    grantPermissionUrl := GetGrantUrl(c)
    c.Redirect(http.StatusFound, grantPermissionUrl)
}

func GrantOauthPostback(c *gin.Context){
    var urlQueryData RedirectURL
    if err := c.ShouldBind(&urlQueryData); err != nil{
        log.Println(err)
        c.JSON(400, gin.H{"msg": "Olhe os logs"})
        return
    }
    env := c.Value("env").(Settings)
    config := GetGoogleConfig(env)
    tok, err := config.Exchange(context.Background(), urlQueryData.Code)
    if err != nil{
        log.Println(err)
        c.JSON(400, gin.H{"msg": "Olhe os logs"})
        return
    }
    SaveToken(env.GoogleTokenFile, tok)
    c.JSON(200, gin.H{"message": "Tudo Certo :)"})
}

func main() {
    r := gin.Default()
    r.Use(EnvMiddleware()) 
    r.GET("/", GrantOauthPostback)
    r.GET("/update", UpdateStocks)
    r.GET("/grant-permissions", GrantOauthPermission)
    r.Run()
}

package main

import (
	"context"
	"fmt"
	"net/http"

	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type RedirectURL struct{
    State   string      `form:"state"`
    Code    string      `form:"code"`
    Scope   []string    `form:"scope"`
}

func GetTokenFromWeb(config *oauth2.Config){
    url := fmt.Sprintf("%v", config.AuthCodeURL("state-token", oauth2.AccessTypeOffline))
    err := exec.Command("python", "-m", "webbrowser", url).Run()
    CheckError(err)
}

func RunGinServer(ctx context.Context, ch chan *oauth2.Token){
    r := gin.Default()
    r.GET("/", func(c *gin.Context) {
        var urlQueryData RedirectURL
        if err := c.ShouldBind(&urlQueryData); err != nil{
            c.JSON(400, gin.H{"msg": err})
            return
        }
        config := ctx.Value("config").(*oauth2.Config)
        tok, err := config.Exchange(ctx, urlQueryData.Code)
        if err != nil{
            c.JSON(400, gin.H{"msg": err})
            return
        }
        ch <- tok
        close(ch)
        c.JSON(200, gin.H{"message": "Tudo Certo :)"})
    })
    r.Run()
}

func GetGoogleConfig(env Settings)(*oauth2.Config){
    credentialJson, err := os.ReadFile(env.GoogleOAuthCredentials)
    CheckError(err)
    config, err := google.ConfigFromJSON(
        credentialJson, 
        sheets.SpreadsheetsScope,
        sheets.SpreadsheetsReadonlyScope,
    )
    CheckError(err)
    return config
}

func GetHttpClient(ctx context.Context)(*http.Client){
    env := ctx.Value("env").(Settings)
    config := GetGoogleConfig(env)

    tok, err := TokenFromFile(env.GoogleTokenFile)
    if os.IsNotExist(err){
        ctx = context.WithValue(ctx, "config", config)
        ch := make(chan *oauth2.Token)
        go RunGinServer(ctx, ch)
        GetTokenFromWeb(config)
        tok = <-ch
        SaveToken(env.GoogleTokenFile, tok)
    }
    return config.Client(ctx, tok)
}

func GetSheetsService(ctx context.Context)(*sheets.Service, error){
    client := GetHttpClient(ctx)
    return sheets.NewService(ctx, option.WithHTTPClient(client))
}

func GetDriveService(ctx context.Context)(*drive.Service, error){
    client := GetHttpClient(ctx)
    return drive.NewService(ctx, option.WithHTTPClient(client))
}

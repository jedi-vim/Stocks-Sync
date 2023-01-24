package gcloud

import (
    "context"

    "google.golang.org/api/option"
    "google.golang.org/api/sheets/v4"

    "stocks-sync/utils"
)

func GetSheetsService(ctx context.Context)(*sheets.Service, error){
    env := ctx.Value("env").(utils.Settings)
    client, err := GetOauthClient(env)
    if err != nil{
        return nil, err
    }
    return sheets.NewService(ctx, option.WithHTTPClient(client))
}

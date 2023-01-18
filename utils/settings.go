package utils

import (
    "log"

    "github.com/caarlos0/env/v6"
    "github.com/joho/godotenv"
)

type Settings struct{
    DocumentID                     string  `env:"DOCUMENT_ID"`
    GoogleApplicationCredentials   string  `env:"GOOGLE_APPLICATION_CREDENTIALS"`
    GoogleOAuthCredentials         string  `env:"GOOGLE_OAUTH_CREDENTIALS"`
    GoogleTokenFile                string  `env:"GOOGLE_TOKEN_FILE"`
}

func Env() Settings{
    godotenv.Load()
    settings := Settings{}
    if err := env.Parse(&settings);err != nil{
        log.Fatal(err)
    }
    return settings
}


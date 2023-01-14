package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetStockPrice(stockCode string) string{
    url := fmt.Sprintf("https://www.google.com/finance/quote/%s:BVMF", stockCode)

    request, _ := http.NewRequest(http.MethodGet, url, nil)
    request.Header.Add(
        "User-Agent", 
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.60 Safari/537.36",
        )
    response, err := http.DefaultClient.Do(request)

    if err != nil || response.StatusCode != 200{
        log.Fatal("Fail to fecth url")
    }
    defer response.Body.Close()

    doc, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil{
        log.Fatal("Cant decode response")
    }

    priceDiv := doc.Find("div.fxKbKc").First()
    stockPrice := priceDiv.Text()
    stockPrice = strings.ReplaceAll(stockPrice, ".", ",")
    return strings.ReplaceAll(stockPrice, "R$", "")
}

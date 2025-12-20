package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Kullanım: go run scraper.go <url>")
		return
	}
	targetURL := os.Args[1]

	c := colly.NewCollector()
	var htmlContent []byte

	c.OnResponse(func(r *colly.Response) {
		htmlContent = r.Body
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Bağlantı hatası: %d - %v\n", r.StatusCode, err)
	})

	fmt.Println("Veri çekiliyor...")
	c.Visit(targetURL)

	err := ioutil.WriteFile("site_data.html", htmlContent, 0644)
	if err != nil {
		log.Printf("Dosya yazma hatası: %v\n", err)
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	fmt.Println("Tam ekran görüntüsü alınıyor...")
	err = chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.Sleep(2*time.Second),
		chromedp.FullScreenshot(&buf, 100),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("screenshot.png", buf, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("İşlem başarıyla tamamlandı.")
}
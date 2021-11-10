package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"github.com/valyala/fasthttp"
)

var notifier *notify.Notify

var storeNumbers = [2]string{"R733", "R728"}
var parts = []string{"0=MK7M3TH/A"}

const partURL = "https://www.apple.com/th/shop/fulfillment-messages?pl=true&mt=compact&parts.%s&searchNearby=true&store=R733"

func main() {
	// Register notification
	notifier = notify.New()

	telegramService, _ := telegram.New("")
	telegramService.AddReceivers(1332508222)

	notifier.UseServices(telegramService)

	var wg sync.WaitGroup

	for _, part := range parts {
		wg.Add(1)
		go watch(part, &wg)
	}

	wg.Wait()

	log.Println("shutdown")
}

func watch(part string, wg *sync.WaitGroup) {
	c := &fasthttp.Client{}
	defer wg.Done()

	delay := time.Duration(0)
	for {
		if delay != 0 {
			log.Println("Sleep", delay)
			time.Sleep(delay)
		}

		url := fmt.Sprintf(partURL, part)

		log.Println("Checking...")
		statusCode, body, err := c.Get(nil, url)
		if err != nil {
			log.Printf("Error when loading google page through local proxy: %s", err)

			delay = time.Minute
			continue
		}

		if statusCode != fasthttp.StatusOK {
			log.Printf("Unexpected status code: %d. Expecting %d", statusCode, fasthttp.StatusOK)

			delay = time.Minute
			continue
		}

		res := new(Response)
		if err := json.Unmarshal(body, &res); err != nil {
			log.Printf("Unable to parse json response")

			delay = time.Minute
			continue
		}

		for _, store := range res.Body.Content.PickupMessage.Stores {
			for _, part := range store.PartsAvailability {
				if strings.Contains(part.PickupSearchQuote, "ขณะนี้ยังไม่มีจำหน่าย") {
					log.Printf("%s (Pickup %s)", store.StoreName, part.PickupDisplay)
					continue
				}

				go notifier.Send(context.Background(), fmt.Sprintf("[Pick-Up] %s", store.StoreName), part.StorePickupProductTitle)
				delay = time.Nanosecond
			}
		}

		if delay == 0 {
			delay = time.Second * 10
		}
	}
}

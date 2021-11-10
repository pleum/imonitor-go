package main

type Response struct {
	Head struct {
		Status string `json:"status"`
		Data   struct {
		} `json:"data"`
	} `json:"head"`
	Body struct {
		Content struct {
			PickupMessage struct {
				Stores []struct {
					StoreName         string          `json:"storeName"`
					Country           string          `json:"country"`
					StoreNumber       string          `json:"storeNumber"`
					PartsAvailability map[string]Part `json:"partsAvailability"`
				} `json:"stores"`
			} `json:"pickupMessage"`
		} `json:"content"`
	} `json:"body"`
}

type Part struct {
	PickupSearchQuote       string `json:"pickupSearchQuote"`
	PartNumber              string `json:"partNumber"`
	StorePickupProductTitle string `json:"storePickupProductTitle"`
	StorePickupQuote        string `json:"storePickupQuote"`
	PickupDisplay           string `json:"pickupDisplay"`
}

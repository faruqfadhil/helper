package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type SubmitSyncInventory struct {
	Data SubmitSyncInventoryData `json:"data"`
}

type SubmitSyncInventoryData struct {
	SubmitSyncJob SubmitSyncInventoryDataJob `json:"submit_sync_job"`
}

type SubmitSyncInventoryDataJob struct {
	ShopID     string `json:"shop_id"`
	MerchantID int64  `json:"merchant_id"`
}

type ProductDownload struct {
	Data ProductDownloadData `json:"data"`
}

type ProductDownloadData struct {
	ProductDownload ProductDownloadDataP `json:"product_download"`
}
type ProductDownloadDataP struct {
	// "auto_generate": true,
	// "channel": "tokopedia",
	// "force": true,
	// "marketplace_product_id": "string",
	// "marketplace_shop_id": "string",
	// "merchant_id": 0,
	// "shop_id": "string"

	ShopID            string `json:"shop_id"`
	Channel           string `json:"channel"`
	MarketplaceShopID string `json:"marketplace_shop_id"`
	MerchantID        int64  `json:"merchant_id"`
}

func main() {
	csvFile, err := os.Open("shopee.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	var i int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("failed update shop id = ", record[0])
			continue
		}
		if i == 0 {
			i++
			continue
		}

		shopId := strings.TrimSpace(record[0])
		// merchant := strings.TrimSpace(record[0])
		// merchantyInt, err := strconv.Atoi(merchant)
		// if err != nil {
		// 	fmt.Printf("failed update shop id = %s, err: %s", record[1], err.Error())
		// }

		// // payload := SubmitSyncInventory{
		// // 	Data: SubmitSyncInventoryData{
		// // 		SubmitSyncJob: SubmitSyncInventoryDataJob{
		// // 			ShopID:     shopId,
		// // 			MerchantID: int64(merchantyInt),
		// // 		},
		// // 	},
		// // }

		// productsvc := "https://ocmsproductsvc.production-0.shipper.id/v1/product/download"
		// product := ProductDownload{
		// 	Data: ProductDownloadData{
		// 		ProductDownload: ProductDownloadDataP{
		// 			MerchantID:        int64(merchantyInt),
		// 			ShopID:            strings.TrimSpace(record[1]),
		// 			MarketplaceShopID: strings.TrimSpace(record[2]),
		// 			Channel:           "shopee",
		// 		},
		// 	},
		// }

		// payloadJson, err := json.Marshal(product)
		// if err != nil {
		// 	fmt.Println("failed update shop id = ", record[0])
		// 	continue
		// }
		fmt.Println("processing shop id = ", shopId)

		_, err = doPostHTTPCall("", nil, "", true)

		if err != nil {
			fmt.Printf("failed update shop id = %s| %s\n", record[0], err.Error())
			continue
		}

		fmt.Println("tidur dulu")
		time.Sleep(60 * time.Second)
		i++
	}
}

func doPostHTTPCall(targetURL string, body []byte, token string, isPost bool) ([]byte, error) {
	method := "PUT"
	if isPost {
		method = "POST"
	}
	client := http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), method, targetURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	if err != nil {
		return nil, fmt.Errorf("error when initiate request http, err: %v", err)
	}

	fmt.Printf("call to %s.... ", targetURL)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when try to http call, err : %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error when call api, status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error read resp, err: %v", err)
	}
	fmt.Printf("success \n")
	return out, nil
}

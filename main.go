package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/i-mora/notifier/notifiers/chat/messenger"
	"github.com/spf13/viper"
)

const (
	viperDebug = "debug"

	viperSiteURL              = "site.url"
	viperSiteResponseExpected = "site.response.expected"

	viperMailTemplate = "mail.template"
)

func init() {
	viper.BindEnv(viperDebug, "DEBUG")

	viper.BindEnv(viperSiteURL, "SITE_URL")
	viper.BindEnv(viperSiteResponseExpected, "SITE_RESPONSE_EXPECTED")

	viper.BindEnv(viperMailTemplate, "MAIL_TEMPLATE")
}

func main() {
	ctx := context.Background()
	//
	rawURL := viper.GetString(viperSiteURL)

	url, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		panic(err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	bts, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	//
	expected := viper.GetString(viperSiteResponseExpected)
	gotten := string(bts)
	client := messenger.NewClient()

	//
	if expected == gotten {
		// err = client.SendBulkMessages(ctx, "❌ AUN NO TIENEN HABITACIONES DISPONIBLES")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		log.Println("Messages ❌ sent successfully")
		return
	}

	//
	err = client.SendBulkMessages(ctx, "✅ YA HAY HABITACIONES DISPONIBLES!!!")
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendBulkMessages(ctx, "https://direct-book.com/properties/AntarisCintermexDirect?locale=en&_gads_gcid=731761819&_gads_gclabel=4ZV2CPGWq7MBEJuZ99wC&_gha_gcid=100233201&_gha_phid=b9f2abdd-644f-46f1-bb1a-3d21ec66ed86&_src=DemandPlus&booking_source=organic&campaign_id=&checkInDate=2023-03-31&checkOutDate=2023-03-03&country=MX&currency=MXN&device=desktop&meta=Google&room_rate=&utm_source=GoogleHotelAds&items[0][adults]=2&items[0][children]=0&items[0][infants]=0&trackPage=no")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Messages ✅ sent successfully")
}

package main

import (
	"bytes"
	"context"
	// "fmt"
	"io"
	"log"
	"net/http"
	// "net/http/httputil"
	"net/url"
	// "reflect"

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

  reader := bytes.NewReader([]byte(`{"codigoPelicula":"spider-man-a-traves-del-spider-verso","codigoCiudad":"aguascalientes","tipos":"41783,41784,41785,41786,41787,41788,41789,41790,41799,41866,41867,41868,41869,41870,41875,41876,41877,41878,41881,41882"}`))
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), reader)
	if err != nil {
		panic(err)
	}

  request.Header.Set("Content-Type", "application/json")

  // dumped, err := httputil.DumpRequest(request, true)
  // if err != nil {
    // panic(err)
  // }
  // fmt.Println(string(dumped))

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

  // fmt.Println(reflect.DeepEqual(expected, gotten))
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
	err = client.SendBulkMessages(ctx, "✅ YA HAY BOLETOS EN EL VIP!!!")
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendBulkMessages(ctx, "https://cinepolis.com/preventas")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Messages ✅ sent successfully")
}

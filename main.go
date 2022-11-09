package main

import (
	"io"
	"net/http"
	"net/url"

	"github.com/i-mora/notifier/mail/gmail"
	"github.com/spf13/viper"
)

const (
	viperSiteURL              = "site.url"
	viperSiteResponseExpected = "site.response.expected"

	viperMailTemplate = "mail.template"
)

func init() {
	viper.BindEnv(viperSiteURL, "SITE_URL")
	viper.BindEnv(viperSiteResponseExpected, "SITE_RESPONSE_EXPECTED")

	viper.BindEnv(viperMailTemplate, "MAIL_TEMPLATE")
}

func main() {
	//
	rawURL := viper.GetString(viperSiteURL)

	url, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
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

	equal := string(bts) == viper.GetString(viperSiteResponseExpected)
	if equal {
		return
	}

	//
	template := viper.GetString(viperMailTemplate)

	_, err = gmail.SendMail(template)
	if err != nil {
		panic(err)
	}
}

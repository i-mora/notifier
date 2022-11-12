package messenger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

const (
	viperFacebookAPIBase        = "facebook.api.base"
	viperFacebookAPIVersion     = "facebook.api.version"
	viperFacebookAPIPageID      = "facebook.api.page_id"
	viperFacebookAPIAccessToken = "facebook.api.access_token"
	viperFacebookAPIPSIDS       = "facebook.api.psids"
)

func init() {
	viper.BindEnv(viperFacebookAPIBase, "FACEBOOK_API_BASE")
	viper.BindEnv(viperFacebookAPIVersion, "FACEBOOK_API_VERSION")
	viper.BindEnv(viperFacebookAPIPageID, "FACEBOOK_API_PAGE_ID")
	viper.BindEnv(viperFacebookAPIAccessToken, "FACEBOOK_API_ACCESS_TOKEN")
	viper.BindEnv(viperFacebookAPIPSIDS, "FACEBOOK_API_PSIDS")
}

type Client struct {
	*http.Client
}

func NewClient() Client {
	return Client{
		http.DefaultClient,
	}
}

func (c Client) SendBulkMessages(ctx context.Context, message string) error {
	psids := strings.Split(viper.GetString(viperFacebookAPIPSIDS), ",")
	errors := make([]string, 0, len(psids))

	for _, psid := range psids {
		if err := c.sendMessage(ctx, message, psid); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}

	return nil
}

func (c Client) sendMessage(ctx context.Context, message string, psid string) error {
	path, err := url.JoinPath(
		viper.GetString(viperFacebookAPIBase),
		viper.GetString(viperFacebookAPIVersion),
		viper.GetString(viperFacebookAPIPageID),
		"messages",
	)
	if err != nil {
		return err
	}

	builder := map[string]interface{}{
		"recipient": map[string]string{
			"id": psid,
		},
		"messaging_type": "MESSAGE_TAG",
		"tag":            "ACCOUNT_UPDATE",
		"message": map[string]string{
			"text": message,
		},
		"access_token": viper.GetString(viperFacebookAPIAccessToken),
	}

	bts, err := json.Marshal(builder)
	if err != nil {
		return err
	}

	log.Println(string(bts))
	body := bytes.NewReader(bts)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := c.do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		bts, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("couldn't read failed response's body: %w", err)
		}

		return fmt.Errorf(string(bts))
	}

	return nil
}

func (c Client) do(req *http.Request) (*http.Response, error) {
	res, err := c.Do(req)

	if viper.GetBool("debug") {
		bts, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(string(bts))

		bts, err = httputil.DumpResponse(res, true)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(string(bts))
	}

	return res, err
}

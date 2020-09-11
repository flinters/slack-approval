package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// used for flags
	slackUrl       string
	serverUrl      string
	confirmMessage string

	rootCmd = &cobra.Command{
		Use: "slack-approval",
		// Args:
		Short: `slack-approval enables group approval on slack`,
		Long: `slack-approval enables group approval on slack.
Please specify either --slack_url option or SLACK_APPROVAL_SLACK_URL environment variable.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			if slackUrl = viper.GetString("slack_url"); slackUrl == "" {
				cmd.Help()
				os.Exit(255)
			}
			serverUrl += "/events"

			id := createEvent(serverUrl)
			requestToSlack(slackUrl, id, confirmMessage)
			waitEventAndFinish(serverUrl, id)
		},
	}
)

func Execute() {

	rootCmd.Execute()
}

func init() {

	viper.SetEnvPrefix("slack_approval")
	viper.AutomaticEnv()
	flags := rootCmd.Flags()
	flags.StringVarP(&slackUrl, "slack_url", "", "", `slack incoming webhook URL.
If you want to use envirionment variable, use SLACK_APPROVAL_SLACK_URL.`)
	viper.BindPFlag("slack_url", flags.Lookup("slack_url"))

	flags.StringVarP(&confirmMessage, "confirm_message", "", "", "(REQUIRED) confirm message on Slack")
	viper.BindPFlag("confirm_message", flags.Lookup("confirm_message"))
	rootCmd.MarkFlagRequired("confirm_message")

	flags.StringVarP(&serverUrl, "server_url", "", "", "(REQUIRED) slack-approval server URL")
	viper.BindPFlag("server_url", flags.Lookup("server_url"))
	rootCmd.MarkFlagRequired("server_url")
}

type checkEventRespBody struct {
	Status string
}

func waitEventAndFinish(serverUrl string, id string) string {

	for {
		status := checkEventStatus(serverUrl, id)
		fmt.Println(status)

		if status == "timeout" || status == "rejected" {
			os.Exit(1)
		} else if status == "approved" {
			os.Exit(0)
		}
		time.Sleep(time.Second * 10)
	}
}

func checkEventStatus(serverUrl string, id string) string {

	res, err := http.Get(serverUrl + "/" + id)
	// res, err := http.Get(debugUrl + "/" + id)
	if err != nil {
		panic(err)
	}
	var data checkEventRespBody
	err = json.NewDecoder(res.Body).Decode(&data)

	defer res.Body.Close()

	if err != nil {
		panic(err)
	}

	return data.Status
}

type createEventRespBody struct {
	Id string
}

func createEvent(serverUrl string) string {
	fiveMinutesLater := (time.Now().Unix()) + 60*5

	body := `
	{
		"timeout_epoch": ` + strconv.FormatInt(fiveMinutesLater, 10) + `
	}`

	res, err := http.Post(serverUrl, "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var data createEventRespBody
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		panic(err)
	}

	fmt.Println("id: " + data.Id)
	return data.Id
}

func requestToSlack(slackUrl, id string, message string) {
	res, err := http.Post(slackUrl, "application/json", bytes.NewBuffer(buildMessageJson(id, message)))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
}

func buildMessageJson(id, message string) []byte {

	type obj = map[string]interface{}

	data := obj{
		"text": message,
		"blocks": []obj{
			obj{
				"type": "header",
				"text": obj{
					"type":  "plain_text",
					"text":  message,
					"emoji": true,
				},
			},
			obj{
				"type": "divider",
			},
			obj{
				"type": "section",
				"text": obj{
					"type": "mrkdwn",
					"text": "承認しますか",
				},
			},
			obj{
				"type":     "actions",
				"block_id": id,
				"elements": []obj{
					obj{
						"type":  "button",
						"style": "primary",
						"text": obj{
							"type":  "plain_text",
							"text":  "承認する",
							"emoji": true,
						},
						"value": "1",
					},
					obj{
						"type": "button",
						"text": obj{
							"type":  "plain_text",
							"text":  "承認しない",
							"emoji": true,
						},
						"value": "0",
					},
				},
			},
		},
	}
	b, err := json.Marshal(data)
	if err != nil {
		// shouldn't cause error
		panic(err)
	}
	return b
}

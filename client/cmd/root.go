package cmd


import (
	"fmt"
	// "strings"
	"encoding/json"

	"net/http"
	"bytes"
	"time"
	"strconv"



	"github.com/spf13/cobra"
)

func Execute() {
	// var echoTimes int

	// var cmdPrint = &cobra.Command{
	// 	Use:   "print [string to print]",
	// 	Short: "Print anything to the screen",
	// 	Long: `print is for printing anything back to the screen.
	// For many years people have printed back to the screen.`,
	// 	Args: cobra.MinimumNArgs(1),
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		fmt.Println("Print: " + strings.Join(args, " "))
	// 	},
	// }

	// var cmdEcho = &cobra.Command{
	// 	Use:   "echo [string to echo]",
	// 	Short: "Echo anything to the screen",
	// 	Long: `echo is for echoing anything back.
	// Echo works a lot like print, except it has a child command.`,
	// 	Args: cobra.MinimumNArgs(1),
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		fmt.Println("Echo: " + strings.Join(args, " "))
	// 	},
	// }

	// var cmdTimes = &cobra.Command{
	// 	Use:   "times [string to echo]",
	// 	Short: "Echo anything to the screen more times",
	// 	Long: `echo things multiple times back to the user by providing
	// a count and a string.`,
	// 	Args: cobra.MinimumNArgs(1),
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		for i := 0; i < echoTimes; i++ {
	// 		fmt.Println("Echo: " + strings.Join(args, " "))
	// 		}
	// 	},
	// }

	// cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	var rootCmd = &cobra.Command{
		Use: "slack-approval",
		Run: func(cmd *cobra.Command, args []string) {
			confirmMessage := args[0]
			fmt.Println(confirmMessage)
			id := createEvent()
			requestToSlack(id, confirmMessage)
		},
	}
	// rootCmd.AddCommand(cmdPrint, cmdEcho)
	// cmdEcho.AddCommand(cmdTimes)
	rootCmd.Execute()
}

var debugUrl = "https://412c6cfaf8b6.ngrok.io"
var serverUrl = "https://8eiq8vncn4.execute-api.ap-northeast-1.amazonaws.com/dev/events"

type createEventRespBody struct {
    Id string
}

func createEvent() string {
	fiveMinutesLater := (time.Now().Unix()) + 60 * 5

	body := `
	{
		"timeout_epoch": `+ strconv.FormatInt(fiveMinutesLater, 10) +`
	}`

	res, err := http.Post(serverUrl, "application/json", bytes.NewBuffer([]byte(body)))
    if err != nil {
        panic(err)
	}
	var data createEventRespBody
    err = json.NewDecoder(res.Body).Decode(&data)

	defer res.Body.Close()

    if err != nil {
        panic(err)
	}

	fmt.Println("id: " + data.Id)
	return data.Id
}

func requestToSlack(id string, message string) {
	url := "https://hooks.slack.com/services/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"

	// jsonStr := []byte(`{"token":"aaaa"}`)
	// req, err := http.Post(url, bytes.NewBuffer(jsonStr))

	// client := &http.Client{
		// CheckRedirect: redirectPolicyFunc,
	// }

	// http.Get(url)
	res, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(`
	{
		"text": "` + message + `",
		"blocks": [
			{
				"type": "header",
				"text": {
					"type": "plain_text",
					"text": "` + message + `",
					"emoji": true
				}
			},
			{
				"type": "divider"
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "承認しますか"
				}
			},
			{
				"type": "actions",
				"block_id": "` + id + `",
				"elements": [
					{
						"type": "button",
						"style": "primary",
						"text": {
							"type": "plain_text",
							"text": "承認する",
							"emoji": true
						},
						"value": "1"
					},
					{
						"type": "button",
						"text": {
							"type": "plain_text",
							"text": "承認しない",
							"emoji": true
						},
						"value": "0"
					}
				]
			}
		]
	}
	`)))
    defer res.Body.Close()

    if err != nil {
        panic(err)
    }

}

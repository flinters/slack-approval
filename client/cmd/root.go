package cmd


// import (
// 	"fmt"
// 	"os"

// 	homedir "github.com/mitchellh/go-homedir"
// 	"github.com/spf13/cobra"
// 	"github.com/spf13/viper"
// )

// var (
// 	// Used for flags.
// 	cfgFile     string
// 	userLicense string

// 	rootCmd = &cobra.Command{
// 		Use:   "slack-approval",
// 		Short: "Enables group approval on Slack",
// // 		Long: `Cobra is a CLI library for Go that empowers applications.
// // This application is a tool to generate the needed files
// // to quickly create a Cobra application.`,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			fmt.Println("hogeeeeeeee")
// 		},
// 	}
// )

// // Execute executes the root command.
// func Execute() error {
// 	return rootCmd.Execute()
// }

// func init() {
// 	cobra.OnInitialize(initConfig)

// 	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
// 	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
// 	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
// 	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
// 	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
// 	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
// 	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
// 	viper.SetDefault("license", "apache")
// }

// func er(msg interface{}) {
// 	fmt.Println("Error:", msg)
// 	os.Exit(1)
// }

// func initConfig() {
// 	if cfgFile != "" {
// 		// Use config file from the flag.
// 		viper.SetConfigFile(cfgFile)
// 	} else {
// 		// Find home directory.
// 		home, err := homedir.Dir()
// 		if err != nil {
// 			er(err)
// 		}

// 		// Search config in home directory with name ".cobra" (without extension).
// 		viper.AddConfigPath(home)
// 		viper.SetConfigName(".cobra")
// 	}

// 	viper.AutomaticEnv()

// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	}
// }

import (
	// "fmt"
	// "strings"
	"net/http"
	"bytes"
	"time"



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
			requestToSlack("ididid", "どうよ")
			// createEvent()
		},
	}
	// rootCmd.AddCommand(cmdPrint, cmdEcho)
	// cmdEcho.AddCommand(cmdTimes)
	rootCmd.Execute()
}

func createEvent()  {
    now := time.Now()
    secs := now.Unix()

	fmt.Println(secs)
	fmt.Println(secs + 1)
}
func requestToSlack(id string, message string) {
	// url := "https://412c6cfaf8b6.ngrok.io"
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

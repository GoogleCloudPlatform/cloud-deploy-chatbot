/*
Copyright 2021 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package deploybot

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/deploybot/bot"
	"github.com/GoogleCloudPlatform/deploybot/gcpclouddeploy"
)

var (
	chatToken string
	channel   string
	chatApp   string
	theBot    bot.Bot
)

// init is used to make it easier to access secrets and to adapt the code
// to use other chat apps.
func init() {

	var found, found2, found3 bool

	chatToken, found = os.LookupEnv("TOKEN")
	channel, found2 = os.LookupEnv("CHANNEL")
	if !found || !found2 {
		log.Fatalf("please define the TOKEN and CHANNEL env vars")
	}

	// Customise this with your own chat app implementation if not using Slack or Google Chat.
	chatApp, found3 = os.LookupEnv("CHATAPP")

	if !found3 || chatApp == "slack" { // Slack by default
		theBot = &bot.SlackAdapter{BotToken: chatToken}
	}
	if chatApp == "google" {
		theBot = &bot.GChatAdapter{BotToken: chatToken}
	}
}

// CloudFuncPubSubCDOps is an entry point function for Google Cloud Functions
// which is triggered by a PubSub notification using Cloud Deploy's "clouddeploy-operations" topic
func CloudFuncPubSubCDOps(ctx context.Context, m gcpclouddeploy.OpsMessage) error {

	fmt.Printf("{\"message\": \"received: %s | status: %s\", \"severity\":\"info\"}\n", m.Attributes["ResourceType"], m.Attributes["Action"])

	resp, err := theBot.SendMessage(channel, m.Attributes)
	resp = strings.ReplaceAll(resp, "\"", "'")
	if err != nil {
		fmt.Printf("{\"message\":\"error posting to Chat App: %s\", \"severity\":\"error\"}\n", err)
	} else {
		fmt.Printf("{\"message\": \"success posting to Chat App: %s\", \"severity\": \"info\"}\n", resp)
	}

	// no need to ack as per comment box at
	// https://cloud.google.com/functions/docs/calling/pubsub#sample_code
	return nil
}

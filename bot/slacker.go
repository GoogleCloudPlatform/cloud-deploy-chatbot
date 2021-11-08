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

package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const slackApiPostMessage = "https://slack.com/api/chat.postMessage"

type SlackAdapter struct {
	BotToken    string
	URLEndpoint string
}

func (slacker *SlackAdapter) SendMessage(channel string, message map[string]string) (string, error) {

	resource, ok := message["ResourceType"]
	if !ok {
		return "", fmt.Errorf("could not find ResourceType key")
	}

	var msgBlocks []Block

	if resource == "Release" || resource == "Rollout" {
		msgBlocks = GetSlackMsg(message)
	} else {
		return "", fmt.Errorf("resourceType not a Release or a Rollout")
	}

	// To aid in testing
	if slacker.URLEndpoint != "" {
		return chatPostMessage(slacker.BotToken, channel, msgBlocks, slacker.URLEndpoint)
	}

	return chatPostMessage(slacker.BotToken, channel, msgBlocks, slackApiPostMessage)
}

func chatPostMessage(token string, channel string, blockMessage []Block, url string) (string, error) {
	theMsg := SlackMessageWrapper{
		Token:   token,
		Channel: channel,
		Unfurl:  false,
		Blocks:  blockMessage,
	}

	marshalled, err := json.Marshal(theMsg)
	if err != nil {
		return "", fmt.Errorf("while marshalling SlackMessageWrapper we got: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, slackApiPostMessage, bytes.NewBuffer(marshalled))
	if err != nil {
		return "", fmt.Errorf("failed calling NewRequestWithContext: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-type", "application/json; charset=utf-8")

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("couldnt do request: %v", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bod, _ := ioutil.ReadAll(resp.Body)
		return string(bod), nil
	}

	return "", fmt.Errorf("request was not ok: %v", resp.StatusCode)
}

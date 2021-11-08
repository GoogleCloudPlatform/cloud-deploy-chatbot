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

import "fmt"

type SlackMessageWrapper struct {
	Token   string  `json:"token,omitempty"`
	Channel string  `json:"channel,omitempty"`
	Unfurl  bool    `json:"unfurl_links,omitempty"`
	Text    string  `json:"text,omitempty"`
	Blocks  []Block `json:"blocks,omitempty"`
}

type Block struct {
	TypeSectionBlock string        `json:"type,omitempty"`
	Text             *TextBlock    `json:"text,omitempty"`
	Fields           []TextBlock   `json:"fields,omitempty"`
	Elements         []ButtonBlock `json:"elements,omitempty"`
}

type TextBlock struct {
	TypeTextBlock string `json:"type,omitempty"`
	Text          string `json:"text,omitempty"`
	Emoji         bool   `json:"emoji,omitempty"`
}

type ButtonBlock struct {
	TypeButtonBlock string     `json:"type,omitempty"`
	Text            ButtonText `json:"text,omitempty"`
	Style           string     `json:"style,omitempty"`
	Value           string     `json:"value,omitempty"`
}

type ButtonText struct {
	TypeButton string `json:"type,omitempty"`
	Emoji      bool   `json:"emoji,omitempty"`
	Text       string `json:"text,omitempty"`
}

// GetSlackMsgRelease returns a struct representing a "Block Kit" formatted Slack message
// with information about a Release
func GetSlackMsgRelease(atts map[string]string) []Block {
	consoleUrl := fmt.Sprintf("https://console.cloud.google.com/deploy/delivery-pipelines/%s/%s/", atts["Location"], atts["DeliveryPipelineId"])
	deliveryPipe := fmt.Sprintf("%s?project=%s", consoleUrl, atts["ProjectNumber"])
	release := fmt.Sprintf("%sreleases/%s/rollouts?project=%s", consoleUrl, atts["ReleaseId"], atts["ProjectNumber"])

	theHeader := headerHelper(atts)
	statusEmoji := "⏳"

	if atts["Action"] == "Succeed" {
		statusEmoji = "✅"
	} else if atts["Action"] != "Succeed" && atts["Action"] != "Start" {
		statusEmoji = "⚠️"
	}

	return []Block{
		{
			TypeSectionBlock: "header",
			Text: &TextBlock{
				TypeTextBlock: "plain_text",
				Text:          theHeader,
				Emoji:         true,
			},
		},
		{
			TypeSectionBlock: "section",
			Text: &TextBlock{
				TypeTextBlock: "mrkdwn",
				Text:          fmt.Sprintf("*Release: <%s|%s>*", release, atts["ReleaseId"]),
			},
		},
		{
			TypeSectionBlock: "section",
			Text: &TextBlock{
				TypeTextBlock: "mrkdwn",
				Text:          fmt.Sprintf("*Status:* %s %s \n*Where:* <%s|%s>", atts["Action"], statusEmoji, deliveryPipe, atts["DeliveryPipelineId"]),
			},
		},
	}
}

// GetSlackMsgRelease returns a struct representing a "Block Kit" formatted Slack message
// with information about a Rollout
func GetSlackMsgRollout(atts map[string]string) []Block {
	consoleUrl := fmt.Sprintf("https://console.cloud.google.com/deploy/delivery-pipelines/%s/%s/", atts["Location"], atts["DeliveryPipelineId"])
	deliveryPipe := fmt.Sprintf("%s?project=%s", consoleUrl, atts["ProjectNumber"])
	release := fmt.Sprintf("%sreleases/%s/rollouts?project=%s", consoleUrl, atts["ReleaseId"], atts["ProjectNumber"])
	target := fmt.Sprintf("%stargets/%s?project=%s", consoleUrl, atts["TargetId"], atts["ProjectNumber"])

	theHeader := headerHelper(atts)
	statusEmoji := "⏳"

	if atts["Action"] == "Succeed" {
		statusEmoji = "✅"
	} else if atts["Action"] != "Succeed" && atts["Action"] != "Start" {
		statusEmoji = "⚠️"
	}

	return []Block{
		{
			TypeSectionBlock: "header",
			Text: &TextBlock{
				TypeTextBlock: "plain_text",
				Text:          theHeader,
				Emoji:         true,
			},
		},
		{
			TypeSectionBlock: "section",
			Text: &TextBlock{
				TypeTextBlock: "mrkdwn",
				Text:          fmt.Sprintf("*Rollout: <%s|%s>* \n*Target:* <%s|%s>", release, atts["RolloutId"], target, atts["TargetId"]),
			},
		},
		{
			TypeSectionBlock: "section",
			Text: &TextBlock{
				TypeTextBlock: "mrkdwn",
				Text:          fmt.Sprintf("*Status:* %s %s \n*Release:* <%s|%s> \n*Pipeline:* <%s|%s>", atts["Action"], statusEmoji, release, atts["ReleaseId"], deliveryPipe, atts["DeliveryPipelineId"]),
			},
		},
	}
}

func GetSlackMsg(atts map[string]string) []Block {

	if atts["ResourceType"] == "Release" {
		return GetSlackMsgRelease(atts)
	}
	return GetSlackMsgRollout(atts)

}

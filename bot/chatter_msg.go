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
	"fmt"

	"google.golang.org/api/chat/v1"
)

// GetChatMsg returns a struct representing a Message formatted with Google Chat "Cards"
// with information about a Release or Rollout depending on the ResourceType key in atts.
func GetChatMsg(atts map[string]string) *chat.Message {
	consoleUrl := fmt.Sprintf("https://console.cloud.google.com/deploy/delivery-pipelines/%s/%s/", atts["Location"], atts["DeliveryPipelineId"])
	target := fmt.Sprintf("%stargets/%s?project=%s", consoleUrl, atts["TargetId"], atts["ProjectNumber"])
	release := fmt.Sprintf("%sreleases/%s/rollouts?project=%s", consoleUrl, atts["ReleaseId"], atts["ProjectNumber"])

	link := release
	buttonText := "Release"

	theHeader := headerHelper(atts)
	sections := make([]*chat.Section, 0)

	section := &chat.Section{
		Widgets: []*chat.WidgetMarkup{
			{
				KeyValue: &chat.KeyValue{
					TopLabel: "Release",
					Content:  atts["ReleaseId"],
				},
			},
			{
				KeyValue: &chat.KeyValue{
					TopLabel: "Status",
					Content:  atts["Action"],
				},
			},
			{
				KeyValue: &chat.KeyValue{
					TopLabel: "Pipeline",
					Content:  atts["DeliveryPipelineId"],
				},
			},
		},
	}

	// Add a few fields and change the link and button text if this is a Rollout.
	if atts["ResourceType"] == "Rollout" {
		link = target
		buttonText = "Target"

		moreWidgets := []*chat.WidgetMarkup{
			{
				KeyValue: &chat.KeyValue{
					TopLabel: "Rollout",
					Content:  atts["RolloutId"],
				},
			},
			{
				KeyValue: &chat.KeyValue{
					TopLabel: "Target",
					Content:  atts["TargetId"],
				},
			},
		}
		section.Widgets = append(moreWidgets, section.Widgets...)
	}

	buttonSection := &chat.Section{
		Widgets: []*chat.WidgetMarkup{
			{
				Buttons: []*chat.Button{
					{
						TextButton: &chat.TextButton{
							Text: "View " + buttonText,
							OnClick: &chat.OnClick{
								OpenLink: &chat.OpenLink{
									Url: link,
								},
							},
						},
					},
				},
			},
		},
	}

	sections = append(sections, section)
	sections = append(sections, buttonSection)

	card := &chat.Card{
		Header: &chat.CardHeader{
			Title: theHeader,
		},
		Sections: sections,
	}
	cards := []*chat.Card{card}

	return &chat.Message{Cards: cards}

}

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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"google.golang.org/api/chat/v1"
)

type testStruct struct {
	atts          map[string]string
	shouldContain []string
	hasError      bool
}

var testTable = []testStruct{
	{
		map[string]string{"ResourceType": "Release", "Action": "Start", "ReleaseId": "rel-20", "DeliveryPipelineId": "pipe-1", "Location": "us-central1"},
		[]string{"Release", "started"},
		false,
	},
	{
		map[string]string{"ResourceType": "Release", "Action": "Succeed", "ReleaseId": "rel-20", "DeliveryPipelineId": "pipe-1", "Location": "us-central1"},
		[]string{"Release", "completed"},
		false,
	},
	{
		map[string]string{"ResourceType": "Rollout", "Action": "Start", "ReleaseId": "rel-20", "DeliveryPipelineId": "pipe-1", "Location": "us-central1"},
		[]string{"Rollout", "started"},
		false,
	},
	{
		map[string]string{"ResourceType": "Rollout", "Action": "Succeed", "ReleaseId": "rel-20", "DeliveryPipelineId": "pipe-1", "Location": "us-central1"},
		[]string{"Rollout", "completed"},
		false,
	},
	{
		map[string]string{"ResourceType": "Crash", "Action": "Succeed", "ReleaseId": "rel-20", "DeliveryPipelineId": "pipe-1", "Location": "us-central1"},
		[]string{"Rollout", "completed"},
		true,
	},
}

func TestSlackMessageConstructors(t *testing.T) {

	for _, item := range testTable {
		slackMsg := GetSlackMsg(item.atts)

		for _, value := range item.shouldContain {

			if !strings.Contains(slackMsg[0].Text.Text, value) {
				t.Errorf("wanted: %s in: %s", value, slackMsg[0].Text.Text)
			}
		}
	}
}

func TestChatMessageConstructors(t *testing.T) {

	for _, item := range testTable {
		chatMsg := GetChatMsg(item.atts)

		for _, value := range item.shouldContain {

			if !strings.Contains(chatMsg.Cards[0].Header.Title, value) {
				t.Errorf("wanted: %s in: %s", value, chatMsg.Cards[0].Header.Title)
			}
		}
	}
}

func testServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &chat.Message{
			Text: "All Good",
		}
		b, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "unable to marshal request: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(b)
	}))
	return ts
}

func TestPostingMessages(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	token := "dummy"
	gchatBot := &GChatAdapter{BotToken: token, URLEndpoint: ts.URL}
	slackBot := &SlackAdapter{BotToken: token, URLEndpoint: ts.URL}

	// Google Chat
	for _, value := range testTable {
		_, err := gchatBot.SendMessage("some channel", value.atts)

		if value.hasError && err == nil {
			t.Errorf("Expected error with attributes: %v", value.atts)
		}
		if !value.hasError && err != nil {
			t.Errorf("UNexpected error with attributes: %v", value.atts)
		}
	}

	// Slack
	for _, value := range testTable {
		_, err := slackBot.SendMessage("some channel", value.atts)

		if value.hasError && err == nil {
			t.Errorf("Expected error with attributes: %v", value.atts)
		}
		if !value.hasError && err != nil {
			t.Errorf("UNexpected error with attributes: %v", value.atts)
		}
	}
}

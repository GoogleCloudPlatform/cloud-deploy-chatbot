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
	"context"
	"fmt"

	"google.golang.org/api/chat/v1"
	"google.golang.org/api/option"
)

type GChatAdapter struct {
	BotToken    string
	URLEndpoint string
}

func (chatter *GChatAdapter) SendMessage(channel string, message map[string]string) (string, error) {

	resource, ok := message["ResourceType"]
	if !ok {
		return "", fmt.Errorf("could not find ResourceType key")
	}

	msg := &chat.Message{Text: "some other resource"}

	if resource == "Release" || resource == "Rollout" {
		msg = GetChatMsg(message)
	} else {
		return "", fmt.Errorf("resourceType not a Release or a Rollout")
	}

	ctx := context.Background()
	opts := option.WithCredentialsJSON([]byte([]byte(chatter.BotToken)))
	optsScope := option.WithScopes("https://www.googleapis.com/auth/chat.bot")

	var chatService *chat.Service
	var err error

	// To aid in testing
	if chatter.URLEndpoint != "" {
		testURL := option.WithEndpoint(chatter.URLEndpoint)
		chatService, err = chat.NewService(ctx, opts, optsScope, testURL, option.WithoutAuthentication())
	} else {
		chatService, err = chat.NewService(ctx, opts, optsScope)
	}

	if err != nil {
		return "", fmt.Errorf("could not create service: %v", err)
	}

	space := fmt.Sprintf("spaces/%s", channel)
	created := chatService.Spaces.Messages.Create(space, msg)
	messageCreated, err := created.Do()

	if err != nil {
		return "", fmt.Errorf("request was not ok: %v", err)
	}

	return fmt.Sprintf("%v", messageCreated), nil
}

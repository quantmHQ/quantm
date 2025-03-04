// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package defs

import (
	"encoding/json"
)

// Kind constants.
const (
	KindBot  = "bot"
	KindUser = "user"
)

type (
	MessageProviderSlackData struct {
		BotToken      string `json:"bot_token"`
		ChannelID     string `json:"channel_id"`
		ChannelName   string `json:"channel_name"`
		WorkspaceID   string `json:"workspace_id"`
		WorkspaceName string `json:"workspace_name"`
	}

	MessageProviderSlackUserInfo struct {
		BotToken       string `json:"bot_token"`
		UserToken      string `json:"user_token"`
		ProviderUserID string `json:"provider_user_id"`
		ProviderTeamID string `json:"provider_team_id"`
	}

	MessageProviderData interface {
		Marshal() ([]byte, error)
		Unmarshal(data []byte) error
	}
)

// Implement Marshal and Unmarshal for MessageProviderSlackData.
func (m *MessageProviderSlackData) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MessageProviderSlackData) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}

// Implement Marshal and Unmarshal for MessageProviderSlackUserInfo.
func (m *MessageProviderSlackUserInfo) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MessageProviderSlackUserInfo) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}

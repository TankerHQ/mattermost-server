// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"encoding/json"
)

type UserToken struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

func (o *UserToken) ToJson() string {
	b, _ := json.Marshal(o)
	return string(b)
}

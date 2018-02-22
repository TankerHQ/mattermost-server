// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"github.com/TankerHQ/tanker-go/usertoken"
	"github.com/mattermost/mattermost-server/model"
)

func (a *App) GetUserToken(UserID string) (*model.UserToken, error) {

	var tankerConfig = usertoken.Config{
		TrustchainID:         a.Config().TankerSettings.TrustchainId,
		TrustchainPrivateKey: a.Config().TankerSettings.TrustchainPrivateKey,
	}
	// Load existing userToken from DB
	user, err := a.GetUser(UserID)
	if err != nil {
		return nil, err
	}

	// If user doesn't have a userToken, generate a new one
	userToken := user.Props["usertoken"]
	if userToken == "" {
		userToken, err := usertoken.Generate(tankerConfig, UserID)
		if err != nil {
			return nil, err
		}

		user.Props["usertoken"] = userToken
		_, err = a.UpdateUser(user, true)
		if err != nil {
			return nil, err
		}
	}

	return &model.UserToken{UserId: UserID, Token: userToken}, nil
}

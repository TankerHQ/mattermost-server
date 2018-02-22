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
	user, err := a.GetUser(UserID)
	if err != nil {
		return nil, err
	}

	// If user doesn't have a userToken, generate a new one
	userToken, ok := user.Props["usertoken"]
	if !ok || len(userToken) == 0 {
		generatedToken, err2 := usertoken.Generate(tankerConfig, UserID)
		if err2 != nil {
			return nil, err2
		}

		user.Props["usertoken"] = generatedToken
		_, err3 := a.UpdateUser(user, true)
		if err3 != nil {
			return nil, err3
		}
		res := model.UserToken{UserId: UserID, Token: generatedToken}
		return &res, nil
	}
	return &model.UserToken{UserId: UserID, Token: userToken}, nil
}

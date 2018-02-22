// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"log"

	"github.com/TankerHQ/identity-go/identity"
)

func (a *App) GetTankerIdentity(userID string) (*string, error) {
	var tankerConfig = identity.Config{
		TrustchainID:         a.Config().TankerSettings.TrustchainId,
		TrustchainPrivateKey: a.Config().TankerSettings.TrustchainPrivateKey,
	}
	user, err := a.GetUser(userID)
	if err != nil {
		return nil, err
	}

	// If user doesn't have an identity
	tankerIdentity, ok := user.Props["tankeridentity"]
	userToken, ok2 := user.Props["usertoken"]

	log.Println("tankeridentity", tankerIdentity)
	log.Println("usertoken", userToken)

	if !ok || len(tankerIdentity) == 0 {
		// If user doesn't have a user token
		var generatedIdentity *string
		var err2 error
		if !ok2 || len(userToken) == 0 {
			generatedIdentity, err2 = identity.Create(tankerConfig, userID)
		} else {
			generatedIdentity, err2 = identity.UpgradeUserToken(tankerConfig, userID, userToken)
		}
		if err2 != nil {
			return nil, err2
		}
		user.Props["tankeridentity"] = *generatedIdentity
		_, err3 := a.UpdateUser(user, true)
		if err3 != nil {
			return nil, err3
		}
		return generatedIdentity, nil
	}
	return &tankerIdentity, nil
}

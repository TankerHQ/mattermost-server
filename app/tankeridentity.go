// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"github.com/TankerHQ/identity-go/identity"
	"github.com/pkg/errors"
)

type TankerIdentities struct {
	Identity            string `json:"tanker_identity"`
	ProvisionalIdentity string `json:"provisional_identity"`
}

func (a *App) GetTankerIdentity(userID string) (*TankerIdentities, error) {
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
	provisionalIdentity := user.Props["provisionalidentity"]

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
		tankerIdentity = *generatedIdentity
	}
	return &TankerIdentities{
		Identity:            tankerIdentity,
		ProvisionalIdentity: provisionalIdentity,
	}, nil
}

func (a *App) GetTankerPublicIdentities(userIDs []string) ([]string, error) {
	var tankerConfig = identity.Config{
		TrustchainID:         a.Config().TankerSettings.TrustchainId,
		TrustchainPrivateKey: a.Config().TankerSettings.TrustchainPrivateKey,
	}

	var publicIdentities []string
	for _, userID := range userIDs {
		user, err := a.GetUser(userID)
		if err != nil {
			return nil, err
		}
		tankerIdentity, ok := user.Props["tankeridentity"]
		provisionalIdentity, ok2 := user.Props["provisionalidentity"]

		if !ok && !ok2 {
			userToken, ok3 := user.Props["usertoken"]
			if !ok3 {
				return nil, errors.New("no identity for user" + userID)
			}
			generatedIdentity, err2 := identity.UpgradeUserToken(tankerConfig, userID, userToken)
			if err2 != nil {
				return nil, err2
			}
			user.Props["tankeridentity"] = *generatedIdentity
			_, err3 := a.UpdateUser(user, true)
			if err3 != nil {
				return nil, err3
			}
			tankerIdentity = *generatedIdentity
			ok = true
		}
		if ok && len(tankerIdentity) != 0 {
			publicIdentity, err := identity.GetPublicIdentity(tankerIdentity)
			if err != nil {
				return nil, err
			}
			publicIdentities = append(publicIdentities, *publicIdentity)
		} else if ok2 && len(provisionalIdentity) != 0 {
			publicIdentity, err := identity.GetPublicIdentity(provisionalIdentity)
			if err != nil {
				return nil, err
			}
			publicIdentities = append(publicIdentities, *publicIdentity)
		}
	}

	return publicIdentities, nil
}

func (a *App) GetTankerProvisionalIdentity(email string) (string, string) {
	var tankerConfig = identity.Config{
		TrustchainID:         a.Config().TankerSettings.TrustchainId,
		TrustchainPrivateKey: a.Config().TankerSettings.TrustchainPrivateKey,
	}

	provisionalIdentity, err := identity.CreateProvisional(tankerConfig, email)
	if err != nil {
		return "", ""
	}
	publicProvisionalID, err := identity.GetPublicIdentity(*provisionalIdentity)
	if err != nil {
		return "", ""
	}

	return *provisionalIdentity, *publicProvisionalID
}

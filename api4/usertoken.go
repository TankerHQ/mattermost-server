// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"net/http"

	l4g "github.com/alecthomas/log4go"
	"github.com/mattermost/mattermost-server/model"
)

func (api *API) InitUserToken() {
	l4g.Debug("Initializing usertoken api")
	api.BaseRoutes.UserToken.Handle("", api.ApiSessionRequired(getUserToken)).Methods("GET")
}

func getUserToken(c *Context, w http.ResponseWriter, r *http.Request) {
	userToken, err := c.App.GetUserToken(c.App.Session.UserId)
	if err != nil {
		c.Err = model.NewAppError("UserToken", "api.usertoken.cannot_generate.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(userToken.ToJson()))
}

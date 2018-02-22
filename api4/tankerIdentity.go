// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"

	l4g "github.com/alecthomas/log4go"
	"github.com/mattermost/mattermost-server/model"
)

func (api *API) InitTankerIdentity() {
	l4g.Debug("Initializing tankerIdentity api")
	api.BaseRoutes.TankerIdentity.Handle("", api.ApiSessionRequired(getTankerIdentity)).Methods("GET")
}

func getTankerIdentity(c *Context, w http.ResponseWriter, r *http.Request) {
	tankerIdentity, err := c.App.GetTankerIdentity(c.App.Session.UserId)
	if err != nil {
		c.Err = model.NewAppError("TankerIdentity", "api.tankeridentity.cannot_generate.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(*tankerIdentity)
	w.Write([]byte(res))
}

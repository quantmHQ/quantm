// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the Breu  
// Community License Agreement, Version 1.0 located at  
// http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  
// ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  
// OF SUCH LICENSE AGREEMENT. 

package client

import (
	"go.breu.io/ctrlplane/internal/entities"
)

func (c *Client) AppList() ([]entities.App, error) {
	url := "/apps"
	reply := make([]entities.App, 0)

	if err := c.request("GET", url, &reply, nil); err != nil {
		return reply, err
	}

	return reply, nil
}

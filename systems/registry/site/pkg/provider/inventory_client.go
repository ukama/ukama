/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package providers

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	ic "github.com/ukama/ukama/systems/common/initclient"
	"github.com/ukama/ukama/systems/common/rest"
)
  
  const Version = "/v1"
  const SystemName = "inventory"
  
  type InventoryClientProvider interface {
	 ValidateComponent(orgName string,Id string) error
  }
  
  type inventoryProvider struct {
	  R      *rest.RestClient
	  debug  bool
	  icHost string
  }
  
  type ValidateComponentReq struct {
	 Id string
 }
 
 
  func (r *inventoryProvider) GetRestyClient(org string) (*rest.RestClient, error) {
	  url, err := ic.GetHostUrl(ic.CreateHostString(org, SystemName), r.icHost, &org, r.debug)
	  
	  if err != nil {
		  log.Errorf("Failed to resolve inventory address to inventory/component: %v", err)
		  return nil, fmt.Errorf("failed to resolve inventory component address. Error: %v", err)
	  }
  
	  rc := rest.NewRestyClient(url, r.debug)
  
	  return rc, nil
  }

  
  func NewInventoryProvider(Host string, debug bool) *inventoryProvider {
  
	  r := &inventoryProvider{
		  debug:  debug,
		  icHost: Host,
	  }
  
	  return r
  }
  
  func (r *inventoryProvider) ValidateComponent(orgName string, componentId string) error {
    var err error

    /* Get Provider */
    r.R, err = r.GetRestyClient(orgName)
    if err != nil {
        return err
    }

    resp, err := r.R.C.R().
        Get(r.R.URL.String() + Version + "/components/" + componentId)

    if err != nil {
        log.Errorf("Failed to send api request to inventory at %s . Error %s", r.R.URL.String(), err.Error())
        return fmt.Errorf("api request to inventory at %s failure: %v", r.R.URL.String(), err)
    }

    if resp.StatusCode() >= 400 {
        errorMessage := string(resp.Body())
        log.Errorf("Failed to get component from inventory at %s. HTTP resp code %d and Error message is %s", r.R.URL.String(), resp.StatusCode(), errorMessage)
        if resp.StatusCode() == 404 {
            return fmt.Errorf("%s: component not found", componentId)
        }
        return fmt.Errorf("%s: %s", componentId, errorMessage)
    }

    return nil
}

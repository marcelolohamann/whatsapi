// Copyright (c) 2021 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package multidevice

import (
	"fmt"

	waBinary "go.mau.fi/whatsmeow/binary"
)

type nodeHandler func(cli *Client, node *waBinary.Node) bool

func handleStreamError(cli *Client, node *waBinary.Node) bool {
	if node.Tag != "stream:error" {
		return false
	}
	code, _ := node.Attrs["code"].(string)
	switch code {
	case "515":
		cli.Log.Debugln("Got 515 code, reconnecting")
		go func() {
			cli.Disconnect()
			err := cli.Connect()
			if err != nil {
				cli.Log.Errorln("Failed to reconnect after 515 code:", err)
			}
		}()
	}
	return true
}

func handleConnectSuccess(cli *Client, node *waBinary.Node) bool {
	if node.Tag != "success" {
		return false
	}
	cli.Log.Infoln("Successfully authenticated")
	// TODO upload prekeys
	//err := cli.sendPassiveIQ(false)
	//if err != nil {
	//	cli.Log.Warnln("Failed to send post-connect passive IQ:", err)
	//}
	return true
}

func (cli *Client) sendPassiveIQ(passive bool) error {
	tag := "active"
	if passive {
		tag = "passive"
	}
	res, err := cli.sendRequest(waBinary.Node{
		Tag: "iq",
		Attrs: map[string]interface{}{
			"to":    waBinary.ServerJID,
			"xmlns": "passive",
			"type":  "set",
		},
		Content: []waBinary.Node{{Tag: tag}},
	})
	if err != nil {
		return err
	}
	fmt.Println("passive iq response:", <-res)
	return nil
}

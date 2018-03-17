package main

import (
	"encoding/json"

	"github.com/Azraid/pasque/app"
	. "github.com/Azraid/pasque/core"
	n "github.com/Azraid/pasque/core/net"
	. "github.com/Azraid/pasque/services/julivonoblitz"
)

func OnCreateRoom(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body CreateRoomMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	roomID := GenerateGuid().String()
	if r, err := cli.SendReq("JuliWorld", "JoinRoom", JoinRoomMsg{
		RoomID: roomID,
		UserID: body.UserID,
		Mode:   body.Mode,
	}); err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInternal), nil)
		return gridData
	} else if r.Header.ErrCode != n.NErrorSucess {
		cli.SendResWithError(req, r.Header.GetError(), nil)
		return gridData
	}

	gd := CreateGridData(req.Header.Key, gridData)
	gd.RoomID = roomID

	cli.SendRes(req, CreateRoomMsgR{RoomID: roomID})
	return gd
}

func OnJoinRoom(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body JoinRoomMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	if r, err := cli.SendReq("JuliWorld", "JoinRoom", JoinRoomMsg{
		RoomID: body.RoomID,
		UserID: body.UserID,
		Mode:   body.Mode,
	}); err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInternal), nil)
		return gridData
	} else if r.Header.ErrCode != n.NErrorSucess {
		cli.SendResWithError(req, r.Header.GetError(), nil)
		return gridData
	}

	gd := CreateGridData(req.Header.Key, gridData)
	gd.RoomID = body.RoomID

	cli.SendRes(req, JoinRoomMsgR{})
	return gd
}

func OnLeaveRoom(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body LeaveRoomMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	if gridData != nil {
		gd := gridData.(*GridData)
		body.RoomID = gd.RoomID

		cli.SendReq("JuliWorld", "LeaveRoom", body)
	}

	rbody := LeaveRoomMsgR{}
	cli.SendRes(req, rbody)

	return gridData
}

func OnPlayRead(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body PlayReadyMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	if gridData == nil {
		cli.SendResWithError(req, RaiseNError(NErrorJulivonoblitzNotFoundRoomID, "not join room yet"), nil)
		return gridData
	}

	gd := gridData.(*GridData)
	body.RoomID = gd.RoomID

	if r, err := cli.SendReq("JuliWorld", "PlayReady", body); err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInternal), nil)
		return gd
	} else if r.Header.ErrCode != n.NErrorSucess {
		cli.SendResWithError(req, r.Header.GetError(), nil)
		return gd
	} else {
		var rbody PlayReadyMsgR
		if err := json.Unmarshal(r.Body, &rbody); err != nil {
			app.ErrorLog(err.Error())
			cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
			return gd
		}

		cli.SendRes(req, rbody)
	}

	return gd
}

func OnDrawGroup(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body DrawGroupMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	if gridData == nil {
		cli.SendResWithError(req, RaiseNError(NErrorJulivonoblitzNotFoundRoomID, "not join room yet"), nil)
		return gridData
	}

	gd := gridData.(*GridData)
	body.RoomID = gd.RoomID

	if r, err := cli.SendReq("JuliWorld", "DrawGroup", body); err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInternal), nil)
		return gd
	} else if r.Header.ErrCode != n.NErrorSucess {
		cli.SendResWithError(req, r.Header.GetError(), nil)
		return gd
	} else {
		var rbody DrawGroupMsgR
		if err := json.Unmarshal(r.Body, &rbody); err != nil {
			app.ErrorLog(err.Error())
			cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
			return gd
		}

		cli.SendRes(req, rbody)
	}

	return gd
}

func OnDrawSingle(cli n.Client, req *n.RequestMsg, gridData interface{}) interface{} {
	var body DrawSingleMsg

	if err := json.Unmarshal(req.Body, &body); err != nil {
		app.ErrorLog(err.Error())
		cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
		return gridData
	}

	if gridData == nil {
		cli.SendResWithError(req, RaiseNError(NErrorJulivonoblitzNotFoundRoomID, "not join room yet"), nil)
		return gridData
	}

	gd := gridData.(*GridData)
	body.RoomID = gd.RoomID

	if r, err := cli.SendReq("JuliWorld", "DrawSingle", body); err != nil {
		cli.SendResWithError(req, RaiseNError(n.NErrorInternal), nil)
		return gd
	} else if r.Header.ErrCode != n.NErrorSucess {
		cli.SendResWithError(req, r.Header.GetError(), nil)
		return gd
	} else {
		var rbody DrawSingleMsgR
		if err := json.Unmarshal(r.Body, &rbody); err != nil {
			app.ErrorLog(err.Error())
			cli.SendResWithError(req, RaiseNError(n.NErrorParsingError), nil)
			return gd
		}

		cli.SendRes(req, rbody)
	}

	return gd
}

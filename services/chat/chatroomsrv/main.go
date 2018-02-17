package main

import (
	"fmt"
	"os"

	"github.com/Azraid/pasque/app"
	co "github.com/Azraid/pasque/core"
	. "github.com/Azraid/pasque/services/chat"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ex) chatroomsrv.exe [eid]")
		os.Exit(1)
	}

	eid := os.Args[1]

	workPath := "./"
	if len(os.Args) == 3 {
		workPath = os.Args[2]
	}

	app.InitApp(eid, "", workPath)

	cli := co.NewClient(eid)
	cli.RegisterGridHandler(co.GetNameOfApiMsg(GetRoomMsg{}), OnGetRoom)
	cli.RegisterGridHandler(co.GetNameOfApiMsg(JoinRoomMsg{}), OnJoinRoom)
	cli.RegisterGridHandler(co.GetNameOfApiMsg(SendChatMsg{}), OnSendChat)

	toplgy := co.Topology{
		Spn:           app.Config.Spn,
		FederatedKey:  "RoomID",
		FederatedApis: cli.ListGridApis()}

	cli.Dial(toplgy)

	app.WaitForShutdown()
	return
}

#!/bin/bash
nohup ./router Router.1 > router.1.out & 
nohup ./svcgate Hello.Gate.1 > hello.svcgate.1.out &
nohup ./apigate HelloGame.ApiGate.1 > hellogame.apigate.1.out &
nohup ./hellosrv Hello.1 > hello.1.out &
nohup ./hellocli HelloGame.1 HelloGame > hellogame.1.out & 

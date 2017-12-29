# Gojiko  ( Golang gijiko ) の todo

- ~~ECHO Request/Response の Marshal/Unmarshal を実装~~ @ 2017/09/03
- CSreq/res の Marshal/Unmarshal を実装
  + ~~IMSI~~
  + ~~Cause~~
  + ~~APN~~
  + ~~MSISDN~~
  + ~~MEI~~
  + ~~ULI~~
  + ~~Serving Network~~
  + ~~RAT Type~~
  + ~~Indication Flag~~
  + ~~F-TEID~~
  + ~~Selection Mode~~
  + ~~PDN Type~~
  + ~~PDN Address Allocation (PAA)~~
  + ~~APN Restriction~~
  + ~~APN-AMBR~~
  + ~~EBI~~
  + ~~PCO - 0003H (DNS Server IPv6 Address Request)~~
  + ~~PCO - 000DH (DNS Server IPv4 Address Request)~~
  + ~~PCO - 000AH (IP address allocation via NAS signalling)~~
  + ~~PCO - 8021H (IPCP)~~
  + ~~Bearer Context - Bearer QoS~~
  + ~~Bearer Context - Charging ID~~
  + ~~Bearer Context~~
    + ~~Bearer Context to be created within Create Session Request~~
    + ~~Bearer Context created within Create Session Response~~
  + ~~CSreq~~
  + ~~CSres~~ @ 2017/10/15

- ~~gtpv2cのTEIDをgtp.Teidに変更~~

- SgwCtrl.CreateSession()を実装
  + ~~APN repo を実装~~
  + SgwCtrl.CreateSession() を実装
    + ~~Sender/Receiverとの連携チャンネルに宛先アドレスも入れるようにする~~
    + ~~CreateSessionでMakeCSReqArg()利用を止める~~
    + SgwCtrlSender を実装
      + ~~OpSPgw関連のFindOrCreate等をSgwCtrlではなくAbsSPgwに実装する~~
      + ~~connの保有者をOpSPgwからabsSPgwに移す~~
    + SgwCtrlReceiver を実装
      + ~~receiverで宛先TEIDを見てパケットを振り分ける~~ @ 2017/12/10
      + ~~受け取ったパケットを解析してセッション情報をアップデートする~~ @ 2017/12/18
      + ~~CreateSessionResのタイムアウトが機能していないので修正する~~ @ 2017/12/21
    + リトライする

- Config管理
  + ~~viper利用で進める~~ @ 2017/12/28
  + ~~S5のtimeout, retry数をコンフィグで指定できるようにする~~ @ 2017/12/28

- logを別レイヤに切り出し


- SgwCtrl.DeleteSession() を実装

- ECHO-Cを実装する
  + ECHOパケットを適切なECHOMgrへ振り分ける。

- UDP通信を実装する
  + SgwData.SendUDP() を実装
  + UdpResponser を実装


- Webサービス部の作成
  + goa利用で進める

- SPgwのgoroutine終了。UDP read timeoutもさせる必要がある。

- Monitor を実装

# Gojiko  ( Golang gijiko ) の todo

- ~~ECHO Request/Response の Marshal/Unmarshal を実装~~ @ 2017/09/03
- ~~CSreq/res の Marshal/Unmarshal を実装~~
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

- ~~SgwCtrl.CreateSession()を実装~~
  + ~~APN repo を実装~~
  + ~~SgwCtrl.CreateSession() を実装~~
    + ~~Sender/Receiverとの連携チャンネルに宛先アドレスも入れるようにする~~
    + ~~CreateSessionでMakeCSReqArg()利用を止める~~
    + ~~SgwCtrlSender を実装~~
      + ~~OpSPgw関連のFindOrCreate等をSgwCtrlではなくAbsSPgwに実装する~~
      + ~~connの保有者をOpSPgwからabsSPgwに移す~~
    + ~~SgwCtrlReceiver を実装~~
      + ~~receiverで宛先TEIDを見てパケットを振り分ける~~ @ 2017/12/10n
      + ~~受け取ったパケットを解析してセッション情報をアップデートする~~ @ 2017/12/18
      + ~~CreateSessionResのタイムアウトが機能していないので修正する~~ @ 2017/12/21
    + ~~リトライする~~ @2017/12/29
    + ~~CSresの戻り値に応じた動作にする~~ @2018/01/01
    + ~~NG時にsessionを削除する~~ @2018/01/01

- ~~Config管理~~
  + ~~viper利用で進める~~ @ 2017/12/28
  + ~~S5のtimeout, retry数をコンフィグで指定できるようにする~~ @ 2017/12/28
  + ~~SGW, APNを定義できるようにする。~~ @ 2018/02/18

- ~~Webサービス部の作成~~
  + ~~goa利用で進める~~ @ 2018/01/13
  + ~~createSession部作成~~ @ 2018/02/15
  + ~~Error時にJSONで返せるようにする~~ @ 2018/02/23
    + 中身はJSON だが Media Type が application/vnd.goa.error になっていて、
      それを+jsonに変えるのは generator を変える必要があり面倒。そのままとする。

- ~~バグ修正~~
  + ~~APN IEが間違っていたので修正~~ @ 2018/02/18
  + ~~APIでsgwが存在しない場合に落ちてしまう~~ @ 2018/02/23

- ~~SgwCtrl.DeleteBearer() を実装~~ @ 2018/02/25

- ~~ECHO-Cを実装する~~
  + ~~GTPv2C ECHOパケットに応答する。~~ @ 2018/02/24
  + ~~定期的にECHO-Cパケットを送出し、死活監視を行う。→要らないな。。~~ @ 2018/03/03

- ~~UDP通信を実装する~~
  + ~~gtpSession.SendUDP() を実装~~ @ 2018/03/01
  + ~~受信機能を実装する~~ @ 2018/03/04
  + ~~statsを実装する~~ @ 2018/03/04
  + ~~delete flow を実装する~~ @ 2018/03/11
  + ~~GETを実装する~~ @ 2018/3/12
  + ~~UdpResponder を実装~~ @ 2018/03/17

- ~~ログレベルを設定可能とする~~ @ 2018/03/03
- ~~logを別レイヤに切り出し。gtpv2c や ie パッケージをどうlogrus対応するか？~~ @ 2018/03/03
- ~~UdpEchoFlowSenderで、チャンネルが詰まったり遅延したりに対応する~~ @ 2018/03/04
- ~~UdpEchoFlowのAPIでエラー時にちゃんと返ってこないのを修正~~ @ 2018/03/11

- ~~sessionのステータス遷移をちゃんとして受け付けられるときだけ受け付けるようにする。channelが詰まったりするのを防ぐため。~~ @ 2018/03/21
- ~~CreateSessionをリファクタリング。コマンドchanでやる必要がない。~~ @ 2018/03/21
- ~~SgwCtrl.DeleteSession() を実装~~ @ 2018/03/25
- ~~UDPflow終了時にend_timeが入らなくなっているのを修正する。~~ @ 2018/03/29
- ~~CreateSessionResponseでMandarotyの扱いが誤っているので修正する。~~ @ 2018/04/01
- ~~UDPResponderのログが出まくるので-v以外は1行だけ出すように修正する。ヘッダ行も出力する。~~ @ 2018/04/01
- ~~UDP Responderで、ログにbitrateを出すようにする ~~ @ 2018/04/01

- GTPv1U ECHOパケットに応答する。
- POST gtpsessionで、応答コードを返却する。
- POST gtpsessionのエラー応答が全部500なのが良くないので修正する。

- Error Indicatorに対応
- contextを使ってキャンセル→goroutine停止を確実に実施する



- SPgwのgoroutine終了。UDP read timeoutもさせる必要がある。

- Monitor を実装


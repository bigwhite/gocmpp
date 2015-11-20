package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bigwhite/gocmpp"
	"github.com/bigwhite/gocmpp/utils"
)

const (
	user           string        = "900001"
	password       string        = "888888"
	connectTimeout time.Duration = time.Second * 2
)

func main() {
	c := cmpp.NewClient(cmpp.V30)
	defer c.Disconnect()
	err := c.Connect(":8888", user, password, connectTimeout)
	if err != nil {
		log.Println("client connect error:", err)
		return
	}
	log.Println("client connect and auth ok")

	t := time.NewTicker(time.Second * 5)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			//submit a message
			cont, err := cmpputils.Utf8ToUcs2("测试gocmpp submit")
			if err != nil {
				fmt.Println("utf8 to ucs2 transform err:", err)
				return
			}
			p := &cmpp.Cmpp3SubmitReqPkt{
				PkTotal:            1,
				PkNumber:           1,
				RegisteredDelivery: 1,
				MsgLevel:           1,
				ServiceId:          "test",
				FeeUserType:        2,
				FeeTerminalId:      "13500002696",
				FeeTerminalType:    0,
				MsgFmt:             8,
				MsgSrc:             "900001",
				FeeType:            "02",
				FeeCode:            "10",
				ValidTime:          "151105131555101+",
				AtTime:             "",
				SrcId:              "900001",
				DestUsrTl:          1,
				DestTerminalId:     []string{"13500002696"},
				DestTerminalType:   0,
				MsgLength:          uint8(len(cont)),
				MsgContent:         cont,
			}

			err = c.SendReqPkt(p)
			if err != nil {
				log.Println("send a cmpp3 submit request error:", err)
			} else {
				log.Println("send a cmpp3 submit request")
			}
			break
		default:
		}

		// recv packets
		i, err := c.RecvAndUnpackPkt(0)
		if err != nil {
			log.Println("client read and unpack pkt error:", err)
			break
		}

		switch p := i.(type) {
		case *cmpp.Cmpp3SubmitRspPkt:
			log.Println("receive a cmpp3 submit response:", p)

		case *cmpp.CmppActiveTestReqPkt:
			log.Println("receive a cmpp active request:", p)
			rsp := &cmpp.CmppActiveTestRspPkt{}
			err := c.SendRspPkt(rsp, p.SeqId)
			if err != nil {
				log.Println("send cmpp active response error:", err)
				break
			}
		case *cmpp.CmppActiveTestRspPkt:
			log.Println("receive a cmpp activetest response:", p)

		case *cmpp.CmppTerminateReqPkt:
			log.Println("receive a cmpp terminate request:", p)
			rsp := &cmpp.CmppTerminateRspPkt{}
			err := c.SendRspPkt(rsp, p.SeqId)
			if err != nil {
				log.Println("send cmpp terminate response error:", err)
				break
			}
		case *cmpp.CmppTerminateRspPkt:
			log.Println("receive a cmpp terminate response:", p)
		}
	}
}

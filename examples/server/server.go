package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"log"
	"time"

	"github.com/bigwhite/gocmpp/packet"
	cmppserver "github.com/bigwhite/gocmpp/server"
)

const (
	userS     string = "900001"
	passwordS string = "888888"
)

func handleLogin(r *cmppserver.Response, p *cmppserver.Packet, l *log.Logger) (bool, error) {
	req, ok := p.Packer.(*cmpppacket.CmppConnReqPkt)
	if !ok {
		// not a connect request, ignore it,
		// go on to next handler
		return true, nil
	}

	resp := r.Packer.(*cmpppacket.Cmpp3ConnRspPkt)

	// validate the user and password
	// set the status in the connect response.
	resp.Version = 0x30
	addr := req.SrcAddr
	if addr != userS {
		l.Println("handleLogin error:", cmpppacket.ConnRspStatusErrMap[cmpppacket.ErrnoConnInvalidSrcAddr])
		resp.Status = uint32(cmpppacket.ErrnoConnInvalidSrcAddr)
		return false, cmpppacket.ConnRspStatusErrMap[cmpppacket.ErrnoConnInvalidSrcAddr]
	}

	tm := req.Timestamp
	authSrc := md5.Sum(bytes.Join([][]byte{[]byte(userS),
		make([]byte, 9),
		[]byte(passwordS),
		[]byte(fmt.Sprintf("%d", tm))},
		nil))

	if req.AuthSrc != string(authSrc[:]) {
		l.Println("handleLogin error: ", cmpppacket.ConnRspStatusErrMap[cmpppacket.ErrnoConnAuthFailed])
		resp.Status = uint32(cmpppacket.ErrnoConnAuthFailed)
		return false, cmpppacket.ConnRspStatusErrMap[cmpppacket.ErrnoConnAuthFailed]
	}

	authIsmg := md5.Sum(bytes.Join([][]byte{[]byte{byte(resp.Status)},
		authSrc[:],
		[]byte(passwordS)},
		nil))
	resp.AuthIsmg = string(authIsmg[:])
	l.Printf("handleLogin: %s login ok\n", addr)

	return false, nil
}

func handleSubmit(r *cmppserver.Response, p *cmppserver.Packet, l *log.Logger) (bool, error) {
	req, ok := p.Packer.(*cmpppacket.Cmpp3SubmitReqPkt)
	if !ok {
		return true, nil // go on to next handler
	}

	resp := r.Packer.(*cmpppacket.Cmpp3SubmitRspPkt)
	resp.MsgId = 12878564852733378560 //0xb2, 0xb9, 0xda, 0x80, 0x00, 0x01, 0x00, 0x00
	for i, d := range req.DestTerminalId {
		l.Printf("handleSubmit: handle submit from %s ok! msgid[%d], srcId[%s], destTerminalId[%s]\n",
			req.MsgSrc, resp.MsgId+uint64(i), req.SrcId, d)
	}
	return true, nil
}

func main() {
	var handlers = []cmppserver.Handler{
		cmppserver.HandlerFunc(handleLogin),
		cmppserver.HandlerFunc(handleSubmit),
	}

	err := cmppserver.ListenAndServe(":8888",
		cmpppacket.V30,
		5*time.Second,
		3,
		nil,
		handlers...,
	)
	if err != nil {
		log.Println("cmpp ListenAndServ error:", err)
	}
	return
}

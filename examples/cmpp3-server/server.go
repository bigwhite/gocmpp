package main

import (
	"bytes"
	"crypto/md5"
	"log"
	"time"

	"github.com/bigwhite/gocmpp"
	"github.com/bigwhite/gocmpp/utils"
)

const (
	userS     string = "900001"
	passwordS string = "888888"
)

func handleLogin(r *cmpp.Response, p *cmpp.Packet, l *log.Logger) (bool, error) {
	req, ok := p.Packer.(*cmpp.CmppConnReqPkt)
	if !ok {
		// not a connect request, ignore it,
		// go on to next handler
		return true, nil
	}

	resp := r.Packer.(*cmpp.Cmpp3ConnRspPkt)

	// validate the user and password
	// set the status in the connect response.
	resp.Version = 0x30
	addr := req.SrcAddr
	if addr != userS {
		l.Println("handleLogin error:", cmpp.ConnRspStatusErrMap[cmpp.ErrnoConnInvalidSrcAddr])
		resp.Status = uint32(cmpp.ErrnoConnInvalidSrcAddr)
		return false, cmpp.ConnRspStatusErrMap[cmpp.ErrnoConnInvalidSrcAddr]
	}

	tm := req.Timestamp
	authSrc := md5.Sum(bytes.Join([][]byte{[]byte(userS),
		make([]byte, 9),
		[]byte(passwordS),
		[]byte(cmpputils.TimeStamp2Str(tm))},
		nil))

	if req.AuthSrc != string(authSrc[:]) {
		l.Println("handleLogin error: ", cmpp.ConnRspStatusErrMap[cmpp.ErrnoConnAuthFailed])
		resp.Status = uint32(cmpp.ErrnoConnAuthFailed)
		return false, cmpp.ConnRspStatusErrMap[cmpp.ErrnoConnAuthFailed]
	}

	authIsmg := md5.Sum(bytes.Join([][]byte{[]byte{byte(resp.Status)},
		authSrc[:],
		[]byte(passwordS)},
		nil))
	resp.AuthIsmg = string(authIsmg[:])
	l.Printf("handleLogin: %s login ok\n", addr)

	return false, nil
}

func handleSubmit(r *cmpp.Response, p *cmpp.Packet, l *log.Logger) (bool, error) {
	req, ok := p.Packer.(*cmpp.Cmpp3SubmitReqPkt)
	if !ok {
		return true, nil // go on to next handler
	}

	resp := r.Packer.(*cmpp.Cmpp3SubmitRspPkt)
	resp.MsgId = 12878564852733378560 //0xb2, 0xb9, 0xda, 0x80, 0x00, 0x01, 0x00, 0x00
	for i, d := range req.DestTerminalId {
		l.Printf("handleSubmit: handle submit from %s ok! msgid[%d], srcId[%s], destTerminalId[%s]\n",
			req.MsgSrc, resp.MsgId+uint64(i), req.SrcId, d)
	}
	return true, nil
}

func main() {
	var handlers = []cmpp.Handler{
		cmpp.HandlerFunc(handleLogin),
		cmpp.HandlerFunc(handleSubmit),
	}

	err := cmpp.ListenAndServe(":8888",
		cmpp.V30,
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

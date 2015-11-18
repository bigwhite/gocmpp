package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"log"
	"time"

	"github.com/bigwhite/gocmpp/conn"
	"github.com/bigwhite/gocmpp/packet"
	cmppserver "github.com/bigwhite/gocmpp/server"
)

const (
	user     string = "900001"
	password string = "888888"
)

func handleLogin(r *cmppserver.Response, p *cmppserver.Packet) (bool, error) {
	req, ok := p.Packer.(*cmpppacket.CmppConnReqPkt)
	if !ok {
		return true, nil // go on to next handler
	}

	resp := r.Packer.(*cmpppacket.Cmpp3ConnRspPkt)

	// validate the user and password
	// set the status in the connect response.
	resp.Version = 0x30
	addr := req.SrcAddr
	if addr != user {
		log.Println("handleLogin: cmpp connect error:", cmpppacket.ErrConnInvalidSrcAddr)
		resp.Status = 2 //invalid source addr
		return false, cmpppacket.ErrConnInvalidSrcAddr
	}

	tm := req.Timestamp
	authSrc := md5.Sum(bytes.Join([][]byte{[]byte(user),
		make([]byte, 9),
		[]byte(password),
		[]byte(fmt.Sprintf("%d", tm))},
		nil))

	if req.AuthSrc != string(authSrc[:]) {
		log.Println("handleLogin: cmpp connect error:", cmpppacket.ErrConnAuthFailed)
		resp.Status = 3 // auth error
		return false, cmpppacket.ErrConnAuthFailed
	}

	authIsmg := md5.Sum(bytes.Join([][]byte{[]byte{byte(resp.Status)},
		authSrc[:],
		[]byte(password)},
		nil))
	resp.AuthIsmg = string(authIsmg[:])

	log.Println("recv a connection from", addr)

	return false, nil
}

func handleSubmit(r *cmppserver.Response, p *cmppserver.Packet) (bool, error) {
	req, ok := p.Packer.(*cmpppacket.Cmpp3SubmitReqPkt)
	if !ok {
		return true, nil // go on to next handler
	}

	_ = req
	return true, nil
}

func main() {
	err := cmppserver.ListenAndServe(":8888",
		cmppconn.V30,
		5*time.Second,
		3,
		cmppserver.HandlerFunc(handleLogin),
		cmppserver.HandlerFunc(handleSubmit),
	)
	if err != nil {
		log.Println("cmpp ListenAndServ error:", err)
	}
	return
}

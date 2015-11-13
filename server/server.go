// Copyright 2015 Tony Bai.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package cmppserver

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	cmppconn "github.com/bigwhite/gocmpp/conn"
	cmpppacket "github.com/bigwhite/gocmpp/packet"
)

//
var (
	ErrEmptyServerAddr = errors.New("cmpp server listen: empty server addr")
	ErrUnsupportedPkt  = errors.New("cmpp server read packet: receive a unsupported pkt")
)

type Packet struct {
	cmpppacket.Packer
	*cmppconn.Conn
}

type Response struct {
	*Packet
	cmpppacket.Packer
	SeqId uint32
}

type Handler interface {
	ServeCmpp(*Response, *Packet) (bool, error)
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as Cmpp handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(*Response, *Packet) (bool, error)

// ServeHTTP calls f(r, p).
func (f HandlerFunc) ServeCmpp(r *Response, p *Packet) (bool, error) {
	return f(r, p)
}

type Server struct {
	Addr    string
	Handler Handler // handler to invoke, protocolValidator if nil
	Typ     cmppconn.Type

	// ErrorLog specifies an optional logger for errors accepting
	// connections and unexpected behavior from handlers.
	// If nil, logging goes to os.Stderr via the log package's
	// standard logger.
	ErrorLog *log.Logger
}

// A conn represents the server side of a Cmpp connection.
type conn struct {
	*cmppconn.Conn
	remoteAddr string  // network address of remote side
	server     *Server // the Server on which the connection arrived
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each.  The service goroutines read requests and
// then call srv.Handler to reply to them.
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		rw, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				srv.ErrorLog.Printf("cmpp: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		tempDelay = 0
		c, err := srv.newConn(rw)
		if err != nil {
			continue
		}
		go c.serve()
	}
}

func (c *conn) readPacket() (*Response, error) {
	i, err := c.Conn.RecvAndUnpackPkt()
	if err != nil {
		return nil, err
	}
	typ := c.server.Typ

	var pkt *Packet
	var rsp *Response
	switch p := i.(type) {
	case *cmpppacket.CmppConnReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		if typ == cmppconn.V30 {
			rsp = &Response{
				Packet: pkt,
				Packer: &cmpppacket.Cmpp3ConnRspPkt{
					SeqId: p.SeqId,
				},
				SeqId: p.SeqId,
			}
		} else {
			rsp = &Response{
				Packet: pkt,
				Packer: &cmpppacket.Cmpp2ConnRspPkt{
					SeqId: p.SeqId,
				},
				SeqId: p.SeqId,
			}
		}

	case *cmpppacket.Cmpp2SubmitReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &cmpppacket.Cmpp2SubmitRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}

	case *cmpppacket.Cmpp3SubmitReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &cmpppacket.Cmpp3SubmitRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}

	case *cmpppacket.Cmpp2FwdReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &cmpppacket.Cmpp2FwdRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}

	case *cmpppacket.Cmpp3FwdReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &cmpppacket.Cmpp3FwdRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}

	case *cmpppacket.Cmpp2DeliverRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
		}

	case *cmpppacket.Cmpp3DeliverRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
		}

	case *cmpppacket.CmppActiveTestReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &cmpppacket.CmppActiveTestRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}

	case *cmpppacket.CmppActiveTestRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
		}

	case *cmpppacket.CmppTerminateReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &cmpppacket.CmppTerminateRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}

	case *cmpppacket.CmppTerminateRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
		}
	default:
		return nil, cmpppacket.NewOpError(ErrUnsupportedPkt,
			fmt.Sprintf("readPacket: receive unsupported packet type: %#v", p))
	}
	return rsp, nil
}

// Close the connection.
func (c *conn) close() {
	p := &cmpppacket.CmppTerminateReqPkt{}

	err := c.Conn.SendPkt(p, <-c.Conn.SeqId)
	if err != nil {
		c.server.ErrorLog.Printf("cmpp: close connection error: %v\n", err)
	}
	c.Conn.Close()
}

func (c *conn) finishPacket(r *Response) error {
	if r.Packer == nil {
		return nil
	}

	return c.Conn.SendPkt(r.Packer, r.SeqId)
}

// Serve a new connection.
func (c *conn) serve() {
	defer func() {
		if err := recover(); err != nil {
			c.server.ErrorLog.Printf("cmpp: panic serving %v: %v\n", c.remoteAddr, err)
		}
	}()

	for {
		r, err := c.readPacket()
		if err != nil {
			c.Close()
			break
		}

		_, err = c.server.Handler.ServeCmpp(r, r.Packet)
		if err != nil {
			c.Close()
			break
		}

		if err := c.finishPacket(r); err != nil {
			c.Close()
			break
		}
	}
}

// Create new connection from rwc.
func (srv *Server) newConn(rwc net.Conn) (c *conn, err error) {
	c = new(conn)
	c.remoteAddr = rwc.RemoteAddr().String()
	c.server = srv
	c.Conn = cmppconn.New(rwc, srv.Typ)
	c.Conn.SetState(cmppconn.CONN_CONNECTED)
	return c, nil
}

func (srv *Server) listenAndServe() error {
	if srv.Addr == "" {
		return ErrEmptyServerAddr
	}
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

func protocolValidate(r *Response, p *Packet) (bool, error) {
	return false, nil
}

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests.
func ListenAndServe(addr string, typ cmppconn.Type, handlers ...Handler) error {
	if addr == "" {
		return ErrEmptyServerAddr
	}

	var handler Handler
	if handlers != nil {
		handler = HandlerFunc(func(r *Response, p *Packet) (bool, error) {
			for _, h := range handlers {
				next, err := h.ServeCmpp(r, p)
				if err != nil || !next {
					return next, err
				}
			}
			return false, nil
		})
	} else {
		handler = HandlerFunc(protocolValidate)
	}
	server := &Server{Addr: addr, Handler: handler, Typ: typ}
	return server.listenAndServe()
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away. the tcpKeepAliveListener's implementation is copied from
// http package.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(1 * time.Minute) // 1min
	return tc, nil
}

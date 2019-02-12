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

package cmpp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync/atomic"
	"time"
)

// errors for cmpp server
var (
	ErrEmptyServerAddr = errors.New("cmpp server listen: empty server addr")
	ErrNoHandlers      = errors.New("cmpp server: no connection handler")
	ErrUnsupportedPkt  = errors.New("cmpp server read packet: receive a unsupported pkt")
)

type Packet struct {
	Packer
	*Conn
}

type Response struct {
	*Packet
	Packer
	SeqId uint32
}

type Handler interface {
	ServeCmpp(*Response, *Packet, *log.Logger) (bool, error)
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as Cmpp handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
//
// The first return value indicates whether to invoke next handler in
// the chain of handlers.
//
// The second return value shows the error returned from the handler. And
// if it is non-nil, server will close the client connection
// after sending back the corresponding response.
type HandlerFunc func(*Response, *Packet, *log.Logger) (bool, error)

// ServeHTTP calls f(r, p).
func (f HandlerFunc) ServeCmpp(r *Response, p *Packet, l *log.Logger) (bool, error) {
	return f(r, p, l)
}

type Server struct {
	Addr    string
	Handler Handler

	// protocol info
	Typ Type
	T   time.Duration // interval betwwen two active tests
	N   int32         // continuous send times when no response back

	// ErrorLog specifies an optional logger for errors accepting
	// connections and unexpected behavior from handlers.
	// If nil, logging goes to os.Stderr via the log package's
	// standard logger.
	ErrorLog *log.Logger
}

// A conn represents the server side of a Cmpp connection.
type conn struct {
	*Conn
	server *Server // the Server on which the connection arrived

	// for active test
	t       time.Duration // interval betwwen two active tests
	n       int32         // continuous send times when no response back
	done    chan struct{}
	exceed  chan struct{}
	counter int32
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
				srv.ErrorLog.Printf("accept error: %v; retrying in %v", e, tempDelay)
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

		srv.ErrorLog.Printf("accept a connection from %v\n", c.Conn.RemoteAddr())
		go c.serve()
	}
}

func (c *conn) readPacket() (*Response, error) {
	readTimeout := time.Second * 2
	i, err := c.Conn.RecvAndUnpackPkt(readTimeout)
	if err != nil {
		return nil, err
	}
	typ := c.server.Typ

	var pkt *Packet
	var rsp *Response
	switch p := i.(type) {
	case *CmppConnReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		if typ == V30 {
			rsp = &Response{
				Packet: pkt,
				Packer: &Cmpp3ConnRspPkt{
					SeqId: p.SeqId,
				},
				SeqId: p.SeqId,
			}
			c.server.ErrorLog.Printf("receive a cmpp30 connect request from %v[%d]\n",
				c.Conn.RemoteAddr(), p.SeqId)
		} else {
			rsp = &Response{
				Packet: pkt,
				Packer: &Cmpp2ConnRspPkt{
					SeqId: p.SeqId,
				},
				SeqId: p.SeqId,
			}
			c.server.ErrorLog.Printf("receive a cmpp20 connect request from %v[%d]\n",
				c.Conn.RemoteAddr(), p.SeqId)
		}

	case *Cmpp2SubmitReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &Cmpp2SubmitRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}
		c.server.ErrorLog.Printf("receive a cmpp20 submit request from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *Cmpp3SubmitReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &Cmpp3SubmitRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}
		c.server.ErrorLog.Printf("receive a cmpp30 submit request from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *Cmpp2FwdReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &Cmpp2FwdRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}
		c.server.ErrorLog.Printf("receive a cmpp20 forward request from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *Cmpp3FwdReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &Cmpp3FwdRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}
		c.server.ErrorLog.Printf("receive a cmpp30 forward request from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *Cmpp2DeliverRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
		}
		c.server.ErrorLog.Printf("receive a cmpp20 deliver response from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *Cmpp3DeliverRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
		}
		c.server.ErrorLog.Printf("receive a cmpp30 deliver response from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *CmppActiveTestReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &CmppActiveTestRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}
		c.server.ErrorLog.Printf("receive a cmpp active request from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *CmppActiveTestRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
		}
		c.server.ErrorLog.Printf("receive a cmpp active response from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *CmppTerminateReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &CmppTerminateRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}
		c.server.ErrorLog.Printf("receive a cmpp terminate request from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *CmppTerminateRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
		}
		c.server.ErrorLog.Printf("receive a cmpp terminate response from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)
	default:
		return nil, NewOpError(ErrUnsupportedPkt,
			fmt.Sprintf("readPacket: receive unsupported packet type: %#v", p))
	}
	return rsp, nil
}

// Close the connection.
func (c *conn) close() {
	p := &CmppTerminateReqPkt{}

	err := c.Conn.SendPkt(p, <-c.Conn.SeqId)
	if err != nil {
		c.server.ErrorLog.Printf("send cmpp terminate request packet to %v error: %v\n", c.Conn.RemoteAddr(), err)
	}

	close(c.done)
	c.server.ErrorLog.Printf("close connection with %v!\n", c.Conn.RemoteAddr())
	c.Conn.Close()
}

func (c *conn) finishPacket(r *Response) error {
	if _, ok := r.Packer.(*CmppActiveTestRspPkt); ok {
		atomic.AddInt32(&c.counter, -1)
		return nil
	}

	if r.Packer == nil {
		// For response packet received, it need not
		// to send anything back.
		return nil
	}

	return c.Conn.SendPkt(r.Packer, r.SeqId)
}

func startActiveTest(c *conn) {
	exceed, done := make(chan struct{}), make(chan struct{})
	c.done = done
	c.exceed = exceed

	go func() {
		t := time.NewTicker(c.t)
		defer t.Stop()
		for {
			select {
			case <-done:
				// once conn close, the goroutine should exit
				return
			case <-t.C:
				// check whether c.counter exceeds
				if atomic.LoadInt32(&c.counter) >= c.n {
					c.server.ErrorLog.Printf("no cmpp active test response returned from %v for %d times!",
						c.Conn.RemoteAddr(), c.n)
					exceed <- struct{}{}
					break
				}
				// send a active test packet to peer, increase the active test counter
				p := &CmppActiveTestReqPkt{}
				err := c.Conn.SendPkt(p, <-c.Conn.SeqId)
				if err != nil {
					c.server.ErrorLog.Printf("send cmpp active test request to %v error: %v", c.Conn.RemoteAddr(), err)
				} else {
					atomic.AddInt32(&c.counter, 1)
				}
			}
		}
	}()
}

// Serve a new connection.
func (c *conn) serve() {
	defer func() {
		if err := recover(); err != nil {
			c.server.ErrorLog.Printf("panic serving %v: %v\n", c.Conn.RemoteAddr(), err)
		}
	}()

	defer c.close()

	// start a goroutine for sending active test.
	startActiveTest(c)

	for {
		select {
		case <-c.exceed:
			return // close the connection.
		default:
		}

		r, err := c.readPacket()
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				continue
			}
			break
		}

		_, err = c.server.Handler.ServeCmpp(r, r.Packet, c.server.ErrorLog)
		if err1 := c.finishPacket(r); err1 != nil {
			break
		}

		if err != nil {
			break
		}
	}
}

// Create new connection from rwc.
func (srv *Server) newConn(rwc net.Conn) (c *conn, err error) {
	c = new(conn)
	c.server = srv
	c.Conn = NewConn(rwc, srv.Typ)
	c.Conn.SetState(CONN_CONNECTED)
	c.n = c.server.N
	c.t = c.server.T
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

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests.
func ListenAndServe(addr string, typ Type, t time.Duration, n int32, logWriter io.Writer, handlers ...Handler) error {
	if addr == "" {
		return ErrEmptyServerAddr
	}

	if handlers == nil {
		return ErrNoHandlers
	}

	var handler Handler
	handler = HandlerFunc(func(r *Response, p *Packet, l *log.Logger) (bool, error) {
		for _, h := range handlers {
			next, err := h.ServeCmpp(r, p, l)
			if err != nil || !next {
				return next, err
			}
		}
		return false, nil
	})

	if logWriter == nil {
		logWriter = os.Stderr
	}
	server := &Server{Addr: addr, Handler: handler, Typ: typ,
		T: t, N: n,
		ErrorLog: log.New(logWriter, "cmppserver: ", log.LstdFlags)}
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

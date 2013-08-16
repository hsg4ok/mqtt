package mqtt

import (
    "errors"
    "fmt"
    "io"
    "log"
    "net"
 //   "net/url"
 //   "strings"
    "time"
    "runtime"
)

const serverDebug = true


type ServerSocket struct {
    ReadTimeout time.Duration // -1 == no timeout
    WriteTimeout time.Duration // -1 == no timeout
    listeners map[string]net.Listener
    shutdown bool
    forceShutdown bool
    connections map[ClientID]*ServerConnection
}

const DefaultServerKeepAlive = KeepAlive(10)
const DefaultServerScheme = MQTTScheme
var DefaultServerReadTimeout = 10 * time.Second
var DefaultServerWriteTimeout = 5 * time.Second

func (ss *ServerSocket) Listen(addr string) (err error) {
    host, port, err := net.SplitHostPort(addr)
    if err != nil {
        host = addr
    }
    if port == "" {
        addr = fmt.Sprintf("%s:%d", host, MQTTPort)
    } else {
        addr = fmt.Sprintf("%s:%s", host, port)
    }
    if ss.listeners[addr] != nil {
        err = errors.New("Already listening on addr")
        return
    }
    var l net.Listener
    l, err = net.Listen("tcp", addr)
    if err != nil {
        return
    }
    if ss.listeners == nil {
        ss.listeners = make(map[string]net.Listener)
    }
    ss.listeners[addr] = l
    return
}

func (ss *ServerSocket) RunServer() {
    for addr, l := range ss.listeners {
        go func() { (*ss).eventLoop(addr, &l) }()
    }
}

func (ss *ServerSocket) eventLoop(addr string, l *net.Listener) {
    for ! ss.shutdown && ! ss.forceShutdown {
        conn, err := (*l).Accept()
        if err != nil {
            continue
        }
        go func() { (*ss).handleConnection(&conn) }()
    }
}

func (ss *ServerSocket) hasConnection(client_id ClientID) bool {
    return ss.connections != nil && ss.connections[client_id] != nil
}

func (ss *ServerSocket) removeConnection(sc *ServerConnection) {
    delete(ss.connections, sc.connectPacket.ClientID)
}

func (ss *ServerSocket) addConnection(sc *ServerConnection) {
    if ss.connections == nil {
        ss.connections = make(map[ClientID]*ServerConnection)
    }
    ss.connections[sc.connectPacket.ClientID] = sc
}

func (ss *ServerSocket) handleConnection(conn *net.Conn) {
    defer (*conn).Close()
    cc := new(ServerConnection)
    cc.ss = ss
    var r io.Reader
    var w io.Writer
    if serverDebug {
        r = &DebugReader{Reader:*conn}
        w = &DebugWriter{Writer:*conn}
    } else {
        r = *conn
        w = *conn
    }
    // CONNECT
    ack, err := expect(r, Connect, cc)
    if err != nil {
        log.Println(err)
        log.Println("Error: expected a CONNECT packet")
        return
    }
    log.Println("@@@@ conected")
    defer ss.removeConnection(cc)
    ss.addConnection(cc)
    // CONNACK
    log.Println("@@@@ sending CONACK")
    err = ConnAckPacket{ConnAckReturnCode: ack}.write(w)
    if err != nil {
        log.Println("Error: unable to write a CONNACK packet")
        return
    }
    // just sit in a loop until an error or DISCONNECT
    // may need another goroutine for sending subscription updates
    cc.ping = make(chan int, 1)
    log.Println("@@@ waiting around")
    for ! ss.shutdown && ! ss.forceShutdown {
        log.Println("reading a packet")
        _, err = read(r, cc)
        log.Println("read a packet")
        if err != nil {
            log.Println(err)
            log.Println("Error: read packet failed")
            return
        }
        select {
        case <- cc.ping:
            log.Println("sending pong")
            err = pingRespPacket(w)
            log.Println("sent pong")
            if err != nil {
                log.Println(err)
                log.Println("Error: write ping packet failed")
                return
            }
        }
    }
}

func (ss *ServerSocket) ForceShutdown() {
    ss.forceShutdown = true
    runtime.Gosched()
}

func (ss *ServerSocket) Shutdown() {
    ss.shutdown = true
}


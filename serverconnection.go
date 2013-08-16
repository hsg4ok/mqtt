package mqtt

import (
    "errors"
    "log"
 //   "net/url"
 //   "strings"
    "sync"
)



type ServerConnection struct {
    ss *ServerSocket
    connectPacket ConnectPacket
    mutex sync.Mutex
    mutex_locked bool
    subscriptions map[Subscription]QoSLevel
    connected bool
    ping chan int
}


func (sc *ServerConnection) connect     (cp ConnectPacket, ack ConnAckReturnCode) (rack ConnAckReturnCode, err error) {
    if ! sc.ss.shutdown && ! sc.ss.forceShutdown {
        rack = ConnAckRefusedUnavailable
        err = errors.New("Server is not running")
        return
    }
    if ack != ConnAckAccepted {
        rack = ack
        err = errors.New("Connection not ascepted")
        return
    }
    if sc.connected {
        rack = ConnAckRefusedNotAuthorized
        err = errors.New("already connected")
        return
    }
    if sc.ss.hasConnection(cp.ClientID) {
       rack = ConnAckRefusedNotAuthorized
       err = errors.New("Duplicate clientid")
       return
    }
    log.Println("CONNECT")
    sc.connectPacket = cp
    sc.connected = true
    return
}


func (sc *ServerConnection) connAck     (ConnAckPacket) error {
    log.Println("CONNACK")
    return errors.New("TODO")
}

func (sc *ServerConnection) publish     (fh FixedHeader, pub PublishPacket) error {
    log.Println("PUBLISH")
    if fh.Dup {
        log.Println("dropping DUP packet")
        return nil
    }
    if fh.Retain {
        sc.ss.retained.retain(fh, pub)
    }
    for client_id, _ /*connection*/ := range sc.ss.connections {
        log.Println("publish processing ", client_id)
    }
    return nil
}
func (sc *ServerConnection) pubAck      (PubAckPacket) error {
    log.Println("PUBACK")
    return errors.New("TODO")
}
func (sc *ServerConnection) pubRec      (PubRecPacket) error {
    log.Println("PUBREC")
    return errors.New("TODO")
}
func (sc *ServerConnection) pubRel      (PubRelPacket) error {
    log.Println("PUBREL")
    return errors.New("TODO")
}
func (sc *ServerConnection) pubComp     (PubCompPacket) error {
    log.Println("PUBCOMP")
    return errors.New("TODO")
}
func (sc *ServerConnection) subscribe   (SubscribePacket) error {
    log.Println("SUBSCRIBE")
    return errors.New("TODO")
}
func (sc *ServerConnection) subAck      (SubAckPacket) error {
    log.Println("SUBACK")
    return errors.New("TODO")
}
func (sc *ServerConnection) unsubscribe (UnsubscribePacket) error {
    log.Println("UNSUBSCRIBE")
    return errors.New("TODO")
}
func (sc *ServerConnection) unsubAck    (UnsubAckPacket) error {
    log.Println("UNSUBACK")
    return errors.New("TODO")
}
func (sc *ServerConnection) pingReq     () error {
    log.Println("PINGREQ")
    sc.ping <- 1 // send a PINGRESP
    return nil
}

func (sc *ServerConnection) pingResp    () error {
    log.Println("PINGRESP")
    return errors.New("Server cannot process PINGRESP packets")
}

func (sc *ServerConnection) disconnect  () error {
    if ! sc.connected {
        return errors.New("already disconnected")
    }
    log.Println("DISCONNECT")
    sc.connected = false
    return nil
}

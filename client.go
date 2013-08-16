package mqtt

import (
    "errors"
    "fmt"
    "io"
    "log"
    "net"
    "net/url"
    "sync"
    "time"
)

const debug = true

const DefaultClientKeepAlive = KeepAlive(10)
const DefaultClientScheme = MQTTScheme
var DefaultClientReadTimeout = 10 * time.Second
var DefaultClientWriteTimeout = 5 * time.Second

type ClientSocket struct {
    ReadTimeout time.Duration // -1 == no timeout
    WriteTimeout time.Duration // -1 == no timeout
    connectPacket ConnectPacket
    connected bool
    ssl bool
    conn net.Conn
    _reader io.Reader
    _writer io.Writer
    url *url.URL
    mutex sync.Mutex
    mutex_locked bool
    _nextMessageID MessageID
}


// Connect
// Publish
// Subscribe
// Unsubscribe
// Ping
// Disconnect


// mqtt     1883/tcp
// mqtt-ssl 8883/tcp 

//func NewSocket(url string) (cs ClientSocket, err error) {
//}

// -- required --
// Host
// KeepAlive
// -- optional --
// ClientID 1 to 23 chars
// Scheme  mqtt (default) or mqtt-ssl
// Port    default 1883 (mqtt) or 8883 (mqtt-ssl)
// Username
// Password
// Retain
// QoS
// WillFlag
// Clean

func (cs *ClientSocket) SetUserName(u string) {
    cs.connectPacket.UserName = String(u)
    cs.connectPacket.ConnectFlags.UserName = true
}

func (cs *ClientSocket) SetPassword(p string) {
    cs.connectPacket.Password = String(p)
    cs.connectPacket.ConnectFlags.Password = true
}

func (cs *ClientSocket) SetWillRetainFlag(r bool) {
    cs.connectPacket.ConnectFlags.WillRetain = r
    cs.connectPacket.ConnectFlags.WillFlag = true
}

func (cs *ClientSocket) SetQoSLevel(q QoSLevel) {
    cs.connectPacket.ConnectFlags.WillQoS = q
    cs.connectPacket.ConnectFlags.WillFlag = true
}

func (cs *ClientSocket) SetCleanSessionFlag(c bool) {
    cs.connectPacket.ConnectFlags.CleanSession = c
}

func (cs *ClientSocket) SetWill(topic, msg string) error {
    cs.connectPacket.WillTopic = Topic(topic)
    will := String(msg)
    if ! validWillMessage(will) {
        return errors.New("invalid will message")
    }
    cs.connectPacket.WillMessage = will
    cs.connectPacket.ConnectFlags.WillFlag = true
    return nil
}

func (cs *ClientSocket) SetClientID(client_id ClientID) (err error) {
    cs.connectPacket.ClientID = client_id
    if ! cs.connectPacket.ClientID.IsValid() {
        err = errors.New("invalid clientid")
        return
    }
    return
}

func (cs *ClientSocket) SetKeepAlive(k KeepAlive) {
    cs.connectPacket.KeepAlive = k
}

func (cs *ClientSocket) setDefaults() {
    if cs.connectPacket.ClientID == "" {
        cs.SetClientID(generateClientID())
    }
    if cs.connectPacket.KeepAlive == 0 {
        cs.SetKeepAlive(DefaultClientKeepAlive)
    }
    if cs.ReadTimeout == 0 {
        cs.ReadTimeout = DefaultClientReadTimeout
    }
    if cs.WriteTimeout == 0 {
        cs.WriteTimeout = DefaultClientWriteTimeout
    }
    return
}

func (cs *ClientSocket) nextMessageID() MessageID {
    cs._nextMessageID++
    return cs._nextMessageID
}

func (cs *ClientSocket) parseURL(connect_url string) (err error) {
    cs.url, err = url.Parse(connect_url)
    if err != nil {
        return
    }
    // check scheme
    if cs.url.Scheme == "" {
        cs.url.Scheme = DefaultClientScheme
    }
    switch cs.url.Scheme {
    case MQTTScheme:
        cs.ssl = false
    case MQTTSSLScheme:
        cs.ssl = true
    default:
        err = errors.New("unknown mqtt url scheme")
        return
    }
    // default port
    host, port, err2 := net.SplitHostPort(cs.url.Host)
    if err2 != nil {
        host = cs.url.Host
    }
    if port == "" {
        var defaultPort int
        if cs.ssl {
            defaultPort = MQTTSSLPort
        } else {
            defaultPort = MQTTPort
        }
        cs.url.Host = fmt.Sprintf("%s:%d", host, defaultPort)
    }
    return
}

func (cs ClientSocket) setReadTimeout() {
    cs.conn.SetReadDeadline(time.Now().Add(cs.ReadTimeout))
}

func (cs ClientSocket) setWriteTimeout() {
    cs.conn.SetWriteDeadline(time.Now().Add(cs.WriteTimeout))
}

func (cs *ClientSocket) Connect(connect_url string) (err error) {
    cs.setDefaults()
    err = cs.parseURL(connect_url)
    if err != nil {
        return
    }

    cs.conn, err = net.Dial("tcp", cs.url.Host)
    if err != nil {
        return
    }
    defer func() {
        if err != nil {
            cs.Close()
            log.Println("Closed")
            return
        }
    }()

    // write a CONNECT packet
    log.Println("writing a CONNECT packet")
    err = cs.connectPacket.write(cs.writer())
    if err != nil {
        log.Println("ERR wrote a CONNECT packet")
        return
    }
    log.Println("wrote a CONNECT packet")

    // read a CONNACK packet
    log.Println("reading a CONNACK packet")
    _, err = expect(cs.reader(), ConnAck, cs)
    if err != nil {
        return
    }
    cs.connected = true
    // start sending keepalive packets
    go cs.keepConnectionAlive()
    return
}

const MaxClientPingErrors = 5

func (cs ClientSocket) reader() io.Reader {
    if cs._reader == nil {
        if debug {
            cs._reader = &DebugReader{Reader: io.Reader(cs.conn)}
        } else {
            cs._reader = io.Reader(cs.conn)
        }
    }
    cs.setReadTimeout()
    return cs._reader
}

func (cs ClientSocket) writer() io.Writer {
    if cs._writer == nil {
        if debug {
            cs._writer = &DebugWriter{Writer: io.Writer(cs.conn)}
        } else {
            cs._writer = io.Writer(cs.conn)
        }
    }
    cs.setWriteTimeout()
    return cs._writer
}

func (cs ClientSocket) keepConnectionAlive() {
    log.Println("keepConnectionAlive() started")
    errors := 0
    for cs.connected && errors < MaxClientPingErrors {
        cs.pingPong(&errors)
        time.Sleep(1 * time.Second)
    }
    log.Println("keepConnectionAlive() finished")
}

func (cs ClientSocket) pingPong(errors *int) {
    cs.lock()
    defer cs.unlock()
    log.Println("Ping?")
    err := pingReqPacket(cs.writer())
    if err != nil {
        log.Println("Ping failed")
        *errors++
        return
    } else {
        log.Println("Ping")
        *errors = 0
    }
    log.Println("Pong?")
    _, err = expect(cs.reader(), PingResp, cs)
    if err != nil {
        log.Println("Pong failed")
        *errors++
    } else {
        *errors = 0
        log.Println("Pong")
    }
}

// PUBLISH an unretained message
func (cs ClientSocket) Publish(topic Topic, qos QoSLevel, data []byte) error {
    switch qos {
    case FireAndForget:
        return cs.PublishFireAndForget(topic, false, data)
    case AcknowledgedDelivery:
        return cs.PublishAcknowledgedDelivery(topic, false, data)
    case AssuredDelivery:
        return cs.PublishAssuredDelivery(topic, false, data)
    default:
        return errors.New("Unknown QoS")
    }
}

// PUBLISH a retained message
func (cs ClientSocket) PublishRetain(topic Topic, qos QoSLevel, retain bool, data []byte) error {
    switch qos {
    case FireAndForget:
        return cs.PublishFireAndForget(topic, true, data)
    case AcknowledgedDelivery:
        return cs.PublishAcknowledgedDelivery(topic, true, data)
    case AssuredDelivery:
        return cs.PublishAssuredDelivery(topic, true, data)
    default:
        return errors.New("Unknown QoS")
    }
}

///    send: PUBLISH (QoS == 0)
// response: none
func (cs ClientSocket) PublishFireAndForget(topic Topic, retain bool, data []byte) error {
    cs.lock()
    defer cs.unlock()
    p := newPacket(Publish)
    p.FixedHeader.Retain = retain
    return PublishPacket{Topic: topic, Data: data}.writePacket(p, cs.writer())
}

///    send: PUBLISH (QoS == 1)
// response: PUBACK
func (cs ClientSocket) PublishAcknowledgedDelivery(topic Topic, retain bool, data []byte) (err error) {
    cs.lock()
    defer cs.unlock()
    // C -> S PUBLISH
    p := newPacket(Publish)
    p.FixedHeader.QoSLevel = AcknowledgedDelivery
    p.FixedHeader.Retain = retain
    err = PublishPacket{Topic: topic, MessageID: cs.nextMessageID(), Data: data}.writePacket(p, cs.writer())
    if err != nil {
        return
    }
    // C <- S PUBACK
    _, err = expect(cs.reader(), PubAck, cs)
    if err != nil {
        return
    }
    return
}

///    send: PUBLISH (QoS == 2)
// response: PUBREC
//     send: PUBREL
// response: PUBCOMP
func (cs ClientSocket) PublishAssuredDelivery(topic Topic, retain bool, data []byte) (err error) {
    cs.lock()
    defer cs.unlock()
    // C -> S PUBLISH
    p := newPacket(Publish)
    p.FixedHeader.QoSLevel = AssuredDelivery
    p.FixedHeader.Retain = retain
    err = PublishPacket{Topic: topic, MessageID: cs.nextMessageID(), Data: data}.writePacket(p, cs.writer())
    if err != nil {
        return
    }
    // C <- S PUBREC
    _, err = expect(cs.reader(), PubRec, cs)
    if err != nil {
        return
    }
    // C -> S PUBREL
    p = newPacket(PubRel)
    p.FixedHeader.QoSLevel = AcknowledgedDelivery
    err = PublishPacket{MessageID: cs._nextMessageID}.writePacket(p, cs.writer())
    if err != nil {
        return
    }
    // C <- S PUBCOMP
    _, err = expect(cs.reader(), PubComp, cs)
    if err != nil {
        return
    }
    return
}

func (cs ClientSocket) lock() {
    log.Println("lock")
    cs.mutex.Lock()
    cs.mutex_locked = true
}

func (cs ClientSocket) unlock() {
    log.Println("unlock")
    if cs.mutex_locked {
        cs.mutex.Unlock()
        cs.mutex_locked = false
    }
}

func (cs ClientSocket) Close() {
    if cs.conn != nil {
        if cs.connected {
            cs.connected = false
            log.Println("Closing")
            cs.lock()
            disconnectPacket(cs.writer())
            cs.unlock()
        }
        cs.conn.Close()
        cs.conn = nil
    }
}

//func (cs ClientSocket) Publish

// implement PacketProcessor

func (cs ClientSocket) connect(cp ConnectPacket, ack ConnAckReturnCode) (rack ConnAckReturnCode, err error) {
    rack = ack
    err = errors.New("Client cannot process CONNECT packets")
    return
}

func (cs ClientSocket) connAck(ack ConnAckPacket) error {
    if ! ack.IsValid() {
        return errors.New(ack.String())
    }
    return nil // ok
}

// server sent a publish because we SUBSCRIBEd to one/more topics
func (cs ClientSocket) publish(fh FixedHeader, p PublishPacket) error {
    return errors.New("TODO")
}

func (cs ClientSocket) pubAck(p PubAckPacket) error {
    if p.MessageID != cs._nextMessageID {
        return errors.New("Wrong PUBACK packet")
    }
    return nil
}

func (cs ClientSocket) pubRec(p PubRecPacket) error {
    return errors.New("TODO")
}

func (cs ClientSocket) pubRel(p PubRelPacket) error {
    return errors.New("TODO")
}

func (cs ClientSocket) pubComp(p PubCompPacket) error {
    return nil
}

func (cs ClientSocket) subscribe(p SubscribePacket) error {
    return errors.New("TODO")
}

func (cs ClientSocket) subAck(p SubAckPacket) error {
    return errors.New("TODO")
}

func (cs ClientSocket) unsubscribe(p UnsubscribePacket) error {
    return errors.New("Client cannot process UNSUBSCRIBE packets")
}

// make sure the message id belongs to an outstanding UNSUBSCRIBE we sent
func (cs ClientSocket) unsubAck(p UnsubAckPacket) error {
    return errors.New("TODO")
}

func (cs ClientSocket) pingReq() error {
    return errors.New("Client cannot process PINGREQ packets")
}

func (cs ClientSocket) pingResp() error {
    return nil
}

func (cs ClientSocket) disconnect() error {
    return errors.New("Client cannot process DISCONNECT packets")
}

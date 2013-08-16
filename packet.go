package mqtt

import (
    "bytes"
    "errors"
    "io"
    "log"
)

type Packet struct {
    FixedHeader
    payloadReader *bytes.Reader
    payloadWriter *bytes.Buffer
}

func (p *Packet) writer() *bytes.Buffer {
    if p.payloadWriter == nil {
        p.payloadWriter = bytes.NewBuffer(nil)
    }
    return p.payloadWriter
}

func newPacket(m MessageType) Packet {
    return Packet{ FixedHeader: FixedHeader{MessageType: m} }
}

func writePacket(w io.Writer, m MessageType) error {
   return newPacket(m).write(w)
}

func disconnectPacket(w io.Writer) error {
    return writePacket(w, Disconnect)
}

func pingReqPacket(w io.Writer) error {
    return writePacket(w, PingReq)
}

func pingRespPacket(w io.Writer) error {
    return writePacket(w, PingResp)
}

func readPacket(r io.Reader) (p Packet, err error) {
    log.Println("readpacket start")
    var fixed_header FixedHeader
    fixed_header, err = readFixedHeader(r)
    if err != nil {
        log.Println("readpacket err")
        return
    }
    var length Length
    length, err = readLength(r)
    if err != nil {
        log.Println("readpacket err")
        return
    }
    var payload []byte
    payload, err = readBytes(int(length), r)
    if err != nil {
        log.Println("readpacket err")
        return
    }
    p = Packet{
        FixedHeader: fixed_header,
        payloadReader: bytes.NewReader(payload),
    }
    log.Println("readpacket done")
    return
}

func dispatchDecodePacket(p Packet, proc PacketProcessor) (ack ConnAckReturnCode, err error) {
    switch p.FixedHeader.MessageType {
    case Connect:
        var c ConnectPacket
        c, ack, err = readConnectPayload(p.payloadReader)
        if err != nil {
            return
        }
        ack, err = proc.connect(c, ack)
    case ConnAck:
        var ca ConnAckPacket
        ca, err = readConnAckPayload(p.payloadReader)
        if err != nil {
            return
        }
        err = proc.connAck(ca)
    case Publish:
        var pub PublishPacket
        pub, err = readPublishPayload(p.payloadReader, p.FixedHeader.QoSLevel)
        if err != nil {
            return
        }
        err = proc.publish(p.FixedHeader, pub)
    case PubAck:
        var pa PubAckPacket
        pa, err = readPubAckPayload(p.payloadReader, p.FixedHeader.QoSLevel)
        if err != nil {
            return
        }
        err = proc.pubAck(pa)
    case PubRec:
        var prec PubRecPacket
        prec, err = readPubRecPayload(p.payloadReader)
        if err != nil {
            return
        }
        err = proc.pubRec(prec)
    case PubRel:
        var prel PubRelPacket
        prel, err = readPubRelPayload(p.payloadReader)
        if err != nil {
            return
        }
        err = proc.pubRel(prel)
    case PubComp:
        var pcomp PubCompPacket
        pcomp, err = readPubCompPayload(p.payloadReader)
        if err != nil {
            return
        }
        err = proc.pubComp(pcomp)
    case Subscribe:
        var sub SubscribePacket
        sub, err = readSubscribePayload(p.payloadReader)
        if err != nil {
            return
        }
        err = proc.subscribe(sub)
    case SubAck:
        var sa SubAckPacket
        sa, err = readSubAckPayload(p.payloadReader)
        if err != nil {
            return
        }
        err = proc.subAck(sa)
    case Unsubscribe:
        var unsub UnsubscribePacket
        unsub, err = readUnsubscribePayload(p.payloadReader)
        if err != nil {
            return
        }
        err = proc.unsubscribe(unsub)
    case UnsubAck:
        var ua UnsubAckPacket
        ua, err = readUnsubAckPayload(p.payloadReader)
        if err != nil {
            return
        }
        err = proc.unsubAck(ua)
    case PingReq:
        err = proc.pingReq()
    case PingResp:
        err = proc.pingResp()
    case Disconnect:
        err = proc.disconnect()
    default:
        err = errors.New("Unknown packet type")
    }
    return
}

func expect(r io.Reader, msg MessageType, proc PacketProcessor) (ack ConnAckReturnCode, err error) {
    p, err := readPacket(r)
    if err != nil {
        return
    }
    if msg != p.FixedHeader.MessageType {
        err = errors.New("Unexpected packet type")
        return
    }
    ack, err = dispatchDecodePacket(p, proc)
    return
}

func read(r io.Reader, proc PacketProcessor) (ack ConnAckReturnCode, err error) {
    log.Println("read()ing a packet")
    p, err := readPacket(r)
    log.Println("read() a packet")
    if err != nil {
        log.Println("packet read error packet")
        return
    }
    log.Println("trying to print packet type")
    log.Println("read a packet of type", p.FixedHeader.String())
    ack, err = dispatchDecodePacket(p, proc)
    return
}

func (p Packet) write(w io.Writer) (err error) {
    log.Println("writing packet", p.FixedHeader.MessageType.String())
    err = p.FixedHeader.write(w)
    if err != nil {
        return
    }
    if p.payloadWriter != nil {
        payload := p.payloadWriter.Bytes()
        l := len(payload)
        err = Length(l).write(w)
        if err != nil {
            return
        }
        var n int
        n, err = w.Write(payload)
        if err != nil {
            return
        }
        if n != int(l) {
            err = errors.New("Unable to write packet")
            return
        }
    } else {
        err = writeByte(0, w)
        if err != nil {
            return
        }
    }
    return
}

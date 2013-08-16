package mqtt

import (
    "log"
    "io"
    "bytes"
    "errors"
)

// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#connect

// Client -> Server
type ConnectPacket struct {
    ConnectFlags
    KeepAlive
    ClientID
    WillTopic   Topic  // optional
    WillMessage String // optional
    UserName    String // optional
    Password    String // optional
}


func (cp ConnectPacket) write(w io.Writer) error {
    return cp.writePacket(newPacket(Connect), w)
}

func (cp ConnectPacket) writePacket(p Packet, w io.Writer) (err error) {
    p.FixedHeader.MessageType = Connect
    err = cp.writePayload(p.writer())
    if err != nil {
        return
    }
    log.Println("writing to network")
    err = p.write(w)
    log.Println("wrote to network")
    return
}

func (cp ConnectPacket) writePayload(w io.Writer) (err error) {
    // bytes 1..8 protocol name
    err = writeBytes(ProtocolName, w)
    if err != nil {
        return
    }
    // byte 9  protocol version
    err = writeByte(ProtocolVersion, w)
    if err != nil {
        return
    }
    // byte 10 connect flags
    err = cp.ConnectFlags.write(w)
    if err != nil {
        return
    }
    // byte 11..12 keep alive timer
    err = cp.KeepAlive.write(w)
    if err != nil {
        return
    }
    log.Println("wrote keepalive", cp.KeepAlive)
    // byte 13+ client id
    err = cp.ClientID.write(w)
    if err != nil {
        return
    }
    log.Println("wrote clientid", cp.ClientID)
    if cp.ConnectFlags.WillFlag {
        err = cp.WillTopic.write(w)
        if err != nil {
            return
        }
        log.Println("writed willtopic")
        err = cp.WillMessage.write(w)
        if err != nil {
            return
        }
        log.Println("writed willmessage")
    }
    if cp.ConnectFlags.UserName {
        err = cp.UserName.write(w)
        if err != nil {
            return
        }
        log.Println("writed username")
    }
    if cp.ConnectFlags.Password {
        err = cp.Password.write(w)
        if err != nil {
            return
        }
        log.Println("writed password")
    }
    log.Println("wrote connect packet payload")
    return
}

func readConnectPayload(r io.Reader) (result ConnectPacket, ack ConnAckReturnCode, err error) {
    // bytes 1..8  protocol name
    protocol_name, err := readBytes(len(ProtocolName), r)
    if err != nil {
        return
    }
    if bytes.Compare(protocol_name, ProtocolName) != 0 {
        ack = ConnAckRefusedProtocolVersion
        err = errors.New("Bad protocol name")
        return
    }
    // byte 9  protocol version
    protocol_version, err := readByte(r)
    if err != nil {
        return
    }
    if protocol_version != ProtocolVersion {
        ack = ConnAckRefusedProtocolVersion
        err = errors.New("Connection Refused: unacceptable protocol version") // 0x04
        return
    }
    // byte 10 connect flags
    result.ConnectFlags, err = readConnectFlags(r)
    if err != nil {
        return
    }
    // byte 11..12 keep alive timer
    result.KeepAlive, err = readKeepAlive(r)
    if err != nil {
        return
    }
    // bytes 13+ client id
    result.ClientID, err = readClientID(r)
    if err != nil {
        ack = ConnAckRefusedIdentifier
        err = errors.New("Connection Refused: identifier rejected")
        return
    }
    // will topic & message
    if result.ConnectFlags.WillFlag {
        var will_topic Topic
        will_topic, err = readTopic(r)
        if err != nil {
            ack = ConnAckRefusedIdentifier
            err = errors.New("Missing or bad will topic")
            return
        }
        var will_message String
        will_message, err = readString(r)
        if err != nil {
            ack = ConnAckRefusedIdentifier
            err = errors.New("Missing will message")
            return
        }
        result.WillTopic = will_topic
        result.WillMessage = will_message
    }
    // username
    if result.ConnectFlags.UserName {
        var username String
        username, err = readString(r)
        if err != nil {
            ack = ConnAckRefusedBadUserOrPass
            return
        }
        result.UserName = username
    }
    // password
    if result.ConnectFlags.Password {
        var password String
        password, err = readString(r)
        if err != nil {
            // backward compatibility for empty username
            if result.ConnectFlags.UserName {
                password = result.UserName
                result.UserName = ""
            } else {
                ack = ConnAckRefusedBadUserOrPass
                err = errors.New("Missing password")
                return
            }
        }
        result.Password = password
    }
    ack = ConnAckAccepted
    return
}

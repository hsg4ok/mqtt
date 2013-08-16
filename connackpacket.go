package mqtt

import (
    "io"
    "errors"
)

type ConnAckReturnCode byte

type ConnAckPacket struct {
    ConnAckReturnCode
}

// Server -> Client

const (
    ConnAckAccepted ConnAckReturnCode = iota // 0
    ConnAckRefusedProtocolVersion // 1
    ConnAckRefusedIdentifier // 2
    ConnAckRefusedUnavailable // 3
    ConnAckRefusedBadUserOrPass // 4
    ConnAckRefusedNotAuthorized // 5
    ConnAckReturnCodeMax = ConnAckRefusedNotAuthorized
)

func (c ConnAckPacket) IsValid() bool {
    return c.ConnAckReturnCode <= ConnAckReturnCodeMax
}

func (c ConnAckPacket) String() string {
    switch c.ConnAckReturnCode {
    case ConnAckAccepted:               return "Connection Accepted" // 0
    case ConnAckRefusedProtocolVersion: return "Connection Refused: unacceptable protocol version" // 1
    case ConnAckRefusedIdentifier:      return "Connection Refused: identifier rejected" // 2
    case ConnAckRefusedUnavailable:     return "Connection Refused: server unavailable" // 3
    case ConnAckRefusedBadUserOrPass:   return "Connection Refused: bad user name or password" // 4
    case ConnAckRefusedNotAuthorized:   return "Connection Refused: not authorized" // 5
    default:                            return "Unknown or reserved value"
    }
}

func (c ConnAckPacket) write(w io.Writer) error {
    return c.writePacket(newPacket(ConnAck), w)
}

func (c ConnAckPacket) writePacket(p Packet, w io.Writer) (err error) {
    p.FixedHeader.MessageType = ConnAck
    err = c.writePayload(p.writer())
    if err != nil {
        return
    }
    err = p.write(w)
    return
}

func (c ConnAckPacket) writePayload(w io.Writer) error {
    if ! c.IsValid() {
        return errors.New("Invalid ConnAck return code")
    }
    return writeBytes([]byte {0x00, byte(c.ConnAckReturnCode)}, w)
}

func readConnAckPayload(r io.Reader) (result ConnAckPacket, err error) {
    // byte 1: reserved (ignored)
    // byte 2: return code
    b, err := readBytes(2, r)
    if err != nil {
        return
    }
    result = ConnAckPacket{ ConnAckReturnCode(b[1]) }
    if ! result.IsValid() {
        err = errors.New("Invalid ConnAck return code")
        return
    }
    return
}

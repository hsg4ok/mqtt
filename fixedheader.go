package mqtt

// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html

import (
    "io"
    "errors"
    "fmt"
    "log"
)

type FixedHeader struct {
    // byte 1
    MessageType // 7..4
    Dup bool // 3
    QoSLevel // 2..1
    Retain bool // 0
    // byte 2..4 length (variable length)
    // byte 3..5 payload (variable length)
}

func (h FixedHeader) write(w io.Writer) error {
    if ! h.QoSLevel.IsValid() {
        return errors.New("invalid qoslevel")
    }
    if ! h.MessageType.IsValid() {
        return errors.New("invalid message type")
    }
    encoded := byte(h.MessageType) << 4
    encoded |= bool2byte(h.Dup)    << 3
    encoded |= byte(h.QoSLevel)    << 1
    encoded |= bool2byte(h.Retain) << 0

    return writeByte(encoded, w)
}

func readFixedHeader(r io.Reader) (result FixedHeader, err error) {
    log.Println("readFixedHeader start")
    log.Println("RFH byte start")
    x, err := readByte(r)
    log.Println("RFH byte done")
    if err != nil {
        return
    }
    result = FixedHeader {
        MessageType: MessageType(x >> 4),
        Dup:         byte2bool(  x >> 3),
        QoSLevel:    QoSLevel((  x >> 1) & 0x03),
        Retain:      byte2bool(  x >> 0),
    }
    if ! result.QoSLevel.IsValid() {
        err = errors.New("invalid qoslevel")
        return
    }
    if ! result.MessageType.IsValid() {
        err = errors.New("invalid message type")
        return
    }
    return
}

func (h FixedHeader) String() string {
    return fmt.Sprintf("%#v", h)
}

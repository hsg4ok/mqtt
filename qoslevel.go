package mqtt

import (
    "io"
    "errors"
)

type QoSLevel byte // 2 bits

const (
    FireAndForget QoSLevel = iota  // 0 : 00b
    AcknowledgedDelivery           // 1 : 01b
    AssuredDelivery                // 2 : 10b
    QoSLevelMax = AssuredDelivery
)

func readQoS(r io.Reader) (result QoSLevel, err error) {
    var tmp byte
    tmp, err = readByte(r)
    if err != nil {
        return
    }
    result = QoSLevel(tmp)
    if ! result.IsValid() {
        err = errors.New("Invalid QoS")
        return
    }
    return
}

func (q QoSLevel) write(w io.Writer) error {
    if ! q.IsValid() {
        return errors.New("Invalid QoS")
    }
    return writeByte(byte(q), w)
}

func (q QoSLevel) IsValid() bool {
    return q <= QoSLevelMax
}

package mqtt

import (
    "bytes"
    "errors"
    "io"
)

type SubAckPacket struct {
    // variable header
    MessageID
    // payload
    GrantedQoS []QoSLevel
}

func (p SubAckPacket) write(w io.Writer) (err error) {
    if len(p.GrantedQoS) == 0 {
        err = errors.New("Invalid grantedqos length")
        return
    }
    err = p.MessageID.write(w)
    if err != nil {
        return
    }
    for _, qos := range p.GrantedQoS {
        err = qos.write(w)
        if err != nil {
            return
        }
    }
    return
}

func readSubAckPayload(r *bytes.Reader) (p SubAckPacket, err error) {
    p.MessageID, err = readMessageID(r)
    if err != nil {
        return
    }
    qos_size := r.Len()
    if qos_size == 0 {
        err = errors.New("Missing grantedqos")
        return
    }
    p.GrantedQoS = make([]QoSLevel, qos_size)
    for i := range p.GrantedQoS {
        p.GrantedQoS[i], err = readQoS(r)
        if err != nil {
            return
        }
    }
    return
}

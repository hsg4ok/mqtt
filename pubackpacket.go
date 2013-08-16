package mqtt

import (
    "io"
    "errors"
)

type PubAckPacket struct {
    MessageID
}

func readPubAckPayload(r io.Reader, qos QoSLevel) (p PubAckPacket, err error) {
    if qos != AcknowledgedDelivery {
        err = errors.New("Should not have received a PUBACK packet")
        return
    }
    p.MessageID, err = readMessageID(r)
    return
}

func (p PubAckPacket) writePayload(w io.Writer) error {
    return p.MessageID.write(w)
}

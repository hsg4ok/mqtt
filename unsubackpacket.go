package mqtt

import (
    "io"
)

type UnsubAckPacket struct {
    MessageID
}

func (p UnsubAckPacket) write(w io.Writer) error {
    return p.MessageID.write(w)
}

func readUnsubAckPayload(r io.Reader) (p UnsubAckPacket, err error) {
    p.MessageID, err = readMessageID(r)
    return
}

package mqtt

import "io"

type PubRecPacket struct {
    MessageID
}

func readPubRecPayload(r io.Reader) (p PubRecPacket, err error) {
    p.MessageID, err = readMessageID(r)
    return
}

func (p PubRecPacket) write(w io.Writer) error {
    return p.MessageID.write(w)
}

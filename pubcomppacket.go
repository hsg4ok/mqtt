package mqtt

import "io"

// This message is either the response from the server to a PUBREL message from a publisher, or the response from a subscriber to a PUBREL message from the server. It is the fourth and last message in the QoS 2 protocol flow.

type PubCompPacket struct {
    MessageID // required
    // no payload
}

func readPubCompPayload(r io.Reader) (p PubCompPacket, err error) {
    p.MessageID, err = readMessageID(r)
    return
}

func (p PubCompPacket) write(w io.Writer) (err error) {
    err = p.MessageID.write(w)
    return
}

package mqtt

import "io"


// C -> S When the server receives a PUBREL message from a publisher, the server makes the original message available to interested subscribers, and sends a PUBCOMP message with the same Message ID to the publisher. 
// C <- S When a subscriber receives a PUBREL message from the server, the subscriber makes the message available to the subscribing application and sends a PUBCOMP message to the server.

type PubRelPacket struct {
    MessageID // required
    // There is no payload.
}

func readPubRelPayload(r io.Reader) (p PubRelPacket, err error) {
    p.MessageID, err = readMessageID(r)
    return
}

func (pr PubRelPacket) write(w io.Writer) error {
    return pr.writePacket(newPacket(PubRel), w)
}

func (pr PubRelPacket) writePacket(p Packet, w io.Writer) error {
    p.FixedHeader.MessageType = PubRel
    p.FixedHeader.Dup = false
    p.FixedHeader.QoSLevel = AcknowledgedDelivery
    return pr.writePayload(p.writer())
}

func (pr PubRelPacket) writePayload(w io.Writer) error {
    return pr.MessageID.write(w)
}

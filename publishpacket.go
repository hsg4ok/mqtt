package mqtt

import (
    "io"
    "bytes"
)

type PublishPacket struct {
    Topic
    MessageID
    Data []byte
}

func (pp PublishPacket) write(w io.Writer) error {
    return pp.writePacket(newPacket(Publish), w)
}

func (pp PublishPacket) writePacket(p Packet, w io.Writer) error {
    p.FixedHeader.MessageType = Publish
    return pp.writePayload(p.writer(), p.FixedHeader.QoSLevel)
}

func (p PublishPacket) writePayload(w io.Writer, qos QoSLevel) (err error) {
    // topic
    err = p.Topic.write(w)
    if err != nil {
        return
    }
    // message id (qos 1 or 2)
    if qos == AcknowledgedDelivery || qos == AssuredDelivery {
        err = p.MessageID.write(w)
        if err != nil {
            return
        }
    }
    // data
    err = writeBytes(p.Data, w)
    return
}

func readPublishPayload(r *bytes.Reader, qos QoSLevel) (result PublishPacket, err error) {
    result.Topic, err = readTopic(r)
    if err != nil {
        return
    }
    if qos == AcknowledgedDelivery || qos == AssuredDelivery {
        result.MessageID, err = readMessageID(r)
        if err != nil {
            return
        }
    }
    result.Data, err = readRemainingBytes(r)
    return
}


package mqtt

import (
    "bytes"
    "errors"
    "io"
)

type Unsubscription struct {
    Topic
}

func (s Unsubscription) write(w io.Writer) (err error) {
    err = s.Topic.write(w)
    return
}

func readUnsubscription(r io.Reader) (s Unsubscription, err error) {
    s.Topic, err = readTopic(r)
    return
}

type UnsubscribePacket struct {
    Unsubscriptions []Unsubscription
}

func (p UnsubscribePacket) write(w io.Writer) (err error) {
    for _, unsub := range p.Unsubscriptions {
        err = unsub.write(w)
        if err != nil {
            return
        }
    }
    return
}

func readUnsubscribePayload(r *bytes.Reader) (s UnsubscribePacket, err error) {
    if r.Len() == 0 {
        err = errors.New("cannot unsubscribe from zero topics")
        return
    }
    var unsub Unsubscription
    for r.Len() > 0 {
        unsub, err = readUnsubscription(r)
        if err != nil {
            return
        }
        s.Unsubscriptions = append(s.Unsubscriptions, unsub)
    }
    return
}

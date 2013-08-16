package mqtt

import (
    "bytes"
    "errors"
    "io"
)


type Subscription struct {
    WildCardTopic
    QoSLevel
}

func (s Subscription) write(w io.Writer) (err error) {
    err = s.WildCardTopic.write(w)
    if err != nil {
        return
    }
    err = s.QoSLevel.write(w)
    return
}

func readSubscription(r io.Reader) (s Subscription, err error) {
    s.WildCardTopic, err = readWildCardTopic(r)
    if err != nil {
        return
    }
    s.QoSLevel, err = readQoS(r)
    return
}

type SubscribePacket struct {
    // variable header
    MessageID
    // payload
    Subscriptions []Subscription
}

func (p SubscribePacket) write(w io.Writer) (err error) {
    err = p.MessageID.write(w)
    if err != nil {
        return
    }
    for _, sub := range p.Subscriptions {
        err = sub.write(w)
        if err != nil {
            return
        }
    }
    return
}

func readSubscribePayload(r *bytes.Reader) (s SubscribePacket, err error) {
    s.MessageID, err = readMessageID(r)
    if err != nil {
        return
    }
    if r.Len() == 0 {
        err = errors.New("zero subscriptions")
        return
    }
    var sub Subscription
    for r.Len() > 0 {
        sub, err = readSubscription(r)
        if err != nil {
            return
        }
        s.Subscriptions = append(s.Subscriptions, sub)
    }
    return
}

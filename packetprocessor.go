package mqtt

type PacketProcessor interface {
    connect     (ConnectPacket, ConnAckReturnCode) (ConnAckReturnCode, error)
    connAck     (ConnAckPacket) error
    publish     (FixedHeader, PublishPacket) error
    pubAck      (PubAckPacket) error
    pubRec      (PubRecPacket) error
    pubRel      (PubRelPacket) error
    pubComp     (PubCompPacket) error
    subscribe   (SubscribePacket) error
    subAck      (SubAckPacket) error
    unsubscribe (UnsubscribePacket) error
    unsubAck    (UnsubAckPacket) error
    pingReq     () error
    pingResp    () error
    disconnect  () error
}

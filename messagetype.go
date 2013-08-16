package mqtt

type MessageType uint8

const (
    Reserved0 MessageType = iota // 0
    Connect                      // 1
    ConnAck                      // 2
    Publish                      // 3
    PubAck                       // 4
    PubRec                       // 5
    PubRel                       // 6
    PubComp                      // 7
    Subscribe                    // 8
    SubAck                       // 9
    Unsubscribe                  // 10
    UnsubAck                     // 11
    PingReq                      // 12
    PingResp                     // 13
    Disconnect                   // 14
    MessageTypeMin = Connect     // 1
    MessageTypeMax = Disconnect  // 14
)

func (m MessageType) IsValid() bool {
    return m >= MessageTypeMin && m <= MessageTypeMax
}

func (m MessageType) String() string {
    switch m {
        case Connect:     return "CONNECT"
        case ConnAck:     return "CONNACK"
        case Publish:     return "PUBLISH"
        case PubAck:      return "PUBACK"
        case PubRec:      return "PUBREC"
        case PubRel:      return "PUBREL"
        case PubComp:     return "PUBCOMP"
        case Subscribe:   return "SUBSCRIBE"
        case SubAck:      return "SUBACK"
        case Unsubscribe: return "UNSUBSCRIBE"
        case UnsubAck:    return "UNSUBACK"
        case PingReq:     return "PINGREQ"
        case PingResp:    return "PINGRESP"
        case Disconnect:  return "DISCONNECT"
        default:          return "(unknown)"
    }
}

package mqtt

import (
    "io"
    "errors"
    "fmt"
)

type ConnectFlags struct {
    UserName     bool // 7
    Password     bool // 6
    WillRetain   bool // 5
    WillQoS      QoSLevel // 4..3
    WillFlag     bool // 2
    CleanSession bool // 1
    // Reserved // 0
}

func (cf ConnectFlags) write(w io.Writer) error {
    if ! cf.WillQoS.IsValid() {
        return errors.New("Invalid WillQoS")
    }
    if ! cf.WillFlag {
        if cf.WillRetain || cf.WillQoS != 0 {
            return errors.New("Will options requires WillFlag")
        }
    }

    encoded := bool2byte(cf.UserName)     << 7
    encoded |= bool2byte(cf.Password)     << 6
    encoded |= bool2byte(cf.WillRetain)   << 5
    encoded |= byte(cf.WillQoS)           << 3
    encoded |= bool2byte(cf.WillFlag)     << 2
    encoded |= bool2byte(cf.CleanSession) << 1

    return writeByte(encoded, w)
}

func readConnectFlags(r io.Reader) (result ConnectFlags, err error) {
    x, err := readByte(r)
    if err != nil {
        return
    }
    if (x & 1) == 1 {
        err = errors.New("Reserved bit should not be used")
        return
    }
    result = ConnectFlags{
        UserName:     byte2bool(x >> 7),
        Password:     byte2bool(x >> 6),
        WillRetain:   byte2bool(x >> 5),
        WillQoS:      QoSLevel(x >> 3) & 0x03,
        WillFlag:     byte2bool(x >> 2),
        CleanSession: byte2bool(x >> 1),
    }
    if ! result.WillQoS.IsValid() {
        err = fmt.Errorf("WillQoS invalid %d", result.WillQoS)
        return
    }
    if ! result.WillFlag {
        if result.WillRetain || result.WillQoS != 0 {
            err = errors.New("Will options requires WillFlag")
            return
        }
    }
    return
}

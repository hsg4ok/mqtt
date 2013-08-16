package mqtt

import (
    "bytes"
    "testing"
)

var fixed_headers = []FixedHeader {
    FixedHeader{ MessageType: Connect, Dup: false, QoSLevel: FireAndForget, Retain: false },
    FixedHeader{ MessageType: Disconnect, Dup: true, QoSLevel: AssuredDelivery, Retain: true },
}

var fixed_headers_writed = []byte {0x10, 0xED}

var fixed_headers_invalid = []FixedHeader {
    FixedHeader{ MessageType: MessageTypeMax+1, Dup: false, QoSLevel: FireAndForget, Retain: false },
    FixedHeader{ MessageType: Disconnect, Dup: true, QoSLevel: QoSLevelMax+1, Retain: true },
}

var fixed_headers_writed_invalid = []byte {0x00, 0xFF}

func TestFixedHeaderWrite(t *testing.T) {
    for i, h := range fixed_headers {
        w := bytes.NewBuffer(nil)
        err := h.write(w)
        if err != nil {
            t.Error("encoding error")
        }
        if bytes.Compare(w.Bytes(), []byte{fixed_headers_writed[i]}) != 0 {
            t.Error("encoding wrong", w.Bytes(), h)
        }
    }
}

func TestFixedHeaderWriteInvalid(t *testing.T) {
    for _, h := range fixed_headers_invalid {
        w := bytes.NewBuffer(nil)
        err := h.write(w)
        if err == nil {
            t.Error("should not have writed")
        }
    }
}

func BenchmarkFixedHeaderWrite(b *testing.B) {
    h := fixed_headers[1]
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        w := bytes.NewBuffer(nil)
        b.StartTimer()
        h.write(w)
    }
}

func TestReadFixedHeader(t *testing.T) {
    for i, e := range fixed_headers_writed {
        r := bytes.NewReader([]byte{e})
        h, err := readFixedHeader(r)
        if err != nil {
            t.Error("decoding error")
        }
        if h != fixed_headers[i] {
            t.Error("decoding wrong")
        }
    }
}

func TestReadFixedHeaderInvalid(t *testing.T) {
    for _, e := range fixed_headers_writed_invalid {
        r := bytes.NewReader([]byte{e})
        h, err := readFixedHeader(r)
        if err == nil {
            t.Error("should not have readd", e, h)
        }
    }
}

func BenchmarkReadFixedHeader(b *testing.B) {
    e := fixed_headers_writed[1]
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        r := bytes.NewBuffer([]byte{e})
        b.StartTimer()
        readFixedHeader(r)
    }
}

package mqtt

import (
    "testing"
    "bytes"
)

var valid_clientids = []ClientID{ "1", "jim", "bob", "12345678901234567890123" }
var valid_clientids_writed = [][]byte{
    []byte{0, 1, 49},
    []byte{0, 3, 106, 105, 109},
    []byte{0, 3, 98, 111, 98},
    []byte{0, 23, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51},
}


var invalid_clientids = []ClientID{ "", "12345678901234567890123toolong" }


func TestClientIDWrite(t *testing.T) {
    for i, x := range valid_clientids {
        w := bytes.NewBuffer(nil)
        err := x.write(w)
        if err != nil {
            t.Error("encoding error")
        }
        if bytes.Compare(w.Bytes(), valid_clientids_writed[i]) != 0 {
            t.Error("encoding wrong")
        }
    }
}

func TestClientIDWriteInvalid(t *testing.T) {
    for _, x := range invalid_clientids {
        w := bytes.NewBuffer(nil)
        err := x.write(w)
        if err == nil {
            t.Error("encoding invalid", w.Bytes())
        }
    }
}

func BenchmarkWriteClientID(b *testing.B) {
    c := ClientID("1")
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        w := bytes.NewBuffer(nil)
        b.StartTimer()
        c.write(w)
    }
}

func TestReadClientID(t *testing.T) {
    for i, x := range valid_clientids_writed {
        r := bytes.NewReader(x)
        c, err := readClientID(r)
        if err != nil {
            t.Error("encoding error")
        }
        if string(c) != string(valid_clientids[i]) {
            t.Error("encoding wrong")
        }
    }
}

func BenchmarkReadClientID(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        r := bytes.NewReader(valid_clientids_writed[1])
        b.StartTimer()
        readClientID(r)
    }
}

func TestGenerateClientID(t *testing.T) {
    for i := 0; i < 10; i ++ {
        c := generateClientID()
        if ! c.IsValid() {
            t.Error("generated an invalid clientid", c)
        }
    }
}

func BenchmarkGenerateClientID(b *testing.B) {
    for i := 0; i < b.N; i++ {
        generateClientID()
    }
}

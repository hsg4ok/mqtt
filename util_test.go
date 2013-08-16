package mqtt

import (
//    "io"
    "bytes"
    "testing"
)

func TestBool2Byte(t *testing.T) {
    if bool2byte(true) != 1 {
        t.Error("true fails")
    }
    if bool2byte(false) != 0 {
        t.Error("false fails")
    }
}

func BenchmarkBool2Byte(b *testing.B) {
    for i := 0; i < b.N; i++ {
        bool2byte(true)
    }
}

func TestByte2Bool(t *testing.T) {
    if byte2bool(1) != true {
        t.Error("true fails")
    }
    if byte2bool(0) != false {
        t.Error("false fails")
    }
}

func BenchmarkByte2Bool(b *testing.B) {
    for i := 0; i < b.N; i++ {
        byte2bool(1)
    }
}

var test_byte_arrays = [][]byte{ []byte{}, []byte{1}, []byte{2, 3} }


func TestReadRemainingBytes(t *testing.T) {
    for _, test_case := range test_byte_arrays {
        r := bytes.NewReader(test_case)
        actual, err := readRemainingBytes(r)
        if err != nil {
            t.Error("unexpected error")
        }
        if bytes.Compare(actual, test_case) != 0 {
            t.Error("read failed")
        }
    }
}

func BenchmarkReadRemainingBytes(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        r := bytes.NewReader(test_byte_arrays[1])
        b.StartTimer()
        readRemainingBytes(r)
    }
}

func TestWriteBytes(t *testing.T) {
    for _, test_case := range test_byte_arrays {
        w := bytes.NewBuffer(nil)
        err := writeBytes(test_case, w)
        if err != nil {
            t.Error("unexpected error")
        }
        if bytes.Compare(w.Bytes(), test_case) != 0 {
            t.Error("write failed")
        }
    }
}

func BenchmarkWriteBytes(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        w := bytes.NewBuffer(nil)
        b.StartTimer()
        writeBytes(test_byte_arrays[1], w)
    }
}

var test_bytes = []byte{0, 1, 53, 255}

func TestWriteByte(t *testing.T) {
    for _, test_case := range test_bytes {
        w := bytes.NewBuffer(nil)
        err := writeByte(test_case, w)
        if err != nil {
            t.Error("unexpected error")
        }
        if bytes.Compare(w.Bytes(), []byte{test_case}) != 0 {
            t.Error("write failed")
        }
    }
}

func BenchmarkWriteByte(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        w := bytes.NewBuffer(nil)
        b.StartTimer()
        writeByte(test_bytes[1], w)
    }
}

var test_uint16 = []uint16 {0, 1, 5, 99, 255, 513, 65535}
var test_uint16_writed = [][]byte {
    []byte{0, 0},
    []byte{0, 1},
    []byte{0, 5},
    []byte{0, 99},
    []byte{0, 255},
    []byte{2, 1},
    []byte{255, 255},
}

func TestWriteUint16(t *testing.T) {
    for i, test_case := range test_uint16 {
        w := bytes.NewBuffer(nil)
        err := writeUint16(test_case, w)
        if err != nil {
            t.Error("unexpected error")
        }

        if bytes.Compare(w.Bytes(), test_uint16_writed[i]) != 0 {
            t.Error("write failed")
        }
   }
}

func BenchmarkWriteUint16(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        w := bytes.NewBuffer(nil)
        b.StartTimer()
        writeUint16(test_uint16[1], w)
    }
}

func TestReadBytes(t *testing.T) {
    for _, test_case := range test_byte_arrays {
        r := bytes.NewReader(test_case)
        actual, err := readBytes(len(test_case), r)
        if err != nil {
            t.Error("unexpected error")
        }
        if bytes.Compare(test_case, actual) != 0 {
            t.Error("read failed")
        }
    }
}

func BenchmarkReadBytes(b *testing.B) {
    bench := test_byte_arrays[1]
    l := len(bench)

    for i := 0; i < b.N; i++ {
        b.StopTimer()
        r := bytes.NewReader(bench)
        b.StartTimer()
        readBytes(l, r)
    }
}

func TestReadByte(t *testing.T) {
    for _, b := range test_bytes {
        r := bytes.NewReader([]byte{b})
        actual, err := readByte(r)
        if err != nil {
            t.Error("unexpected error")
        }
        if b != actual {
            t.Error("read failed")
        }
    }
}

func BenchmarkReadByte(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        r := bytes.NewReader([]byte{3})
        b.StartTimer()
        readByte(r)
    }
}

func TestReadUint16(t *testing.T) {
    for i, test_case := range test_uint16_writed {
        r := bytes.NewReader(test_case)
        actual, err := readUint16(r)
        if err != nil {
            t.Error("unexpected error")
        }
        if actual != test_uint16[i] {
            t.Error("write failed")
        }
    }
}

func BenchmarkReadUint16(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        r := bytes.NewReader(test_uint16_writed[1])
        b.StartTimer()
        readUint16(r)
    }
}

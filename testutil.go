package mqtt

import (
    "io"
    "fmt"
    "log"
    "bytes"
)

type DebugWriter struct {
    io.Writer
}

func (w *DebugWriter) Write(p []byte) (nn int, err error) {
    log.Println("Writing ", decstring(p))
    log.Println("Writing ", hexstring(p))
    log.Println("Writing ", strstring(p))
    nn, err = w.Writer.Write(p)
    if err != nil {
        log.Println("Write error!", err)
    } else if len(p) != nn {
        log.Println("Did not write enough bytes")
    }
    return
}

type DebugReader struct {
    io.Reader
}

func (r *DebugReader) Read(p []byte) (n int, err error) {
    log.Println("Reading")
    n, err = r.Reader.Read(p)
    if err != nil {
        if err == io.EOF {
            log.Println("EOF")
        } else {
            log.Println("Read error!", err)
        }
    } else {
        buf := p[:n]
        log.Println("Read", decstring(buf))
        log.Println("Read", hexstring(buf))
        log.Println("Read", strstring(buf))
    }
    return
}

func hexstring(b []byte) string {
    w := bytes.NewBuffer(nil)
    last := len(b)-1
    w.WriteString(" ")
    for i, x := range b {
        fmt.Fprintf(w, "%02X", x)
        if i != last {
            w.WriteString("  ")
        }
    }
    return w.String()
}

func decstring(b []byte) string {
    w := bytes.NewBuffer(nil)
    last := len(b)-1
    for i, x := range b {
        fmt.Fprintf(w, "%3d", x)
        if i != last {
            w.WriteString(" ")
        }
    }
    return w.String()
}

func strstring(b []byte) string {
    w := bytes.NewBuffer(nil)
    last := len(b)-1
    for i, x := range b {
        if x < 32 || x > 127 {
            x = ' '
        }
        s := fmt.Sprintf("%3c", x)
        if len(s) != 3 {
            s = "   "
        }
        w.WriteString(s)
        if i != last {
            w.WriteString(" ")
        }
    }
    return w.String()
}


func bytes2go(b []byte) string {
    return fmt.Sprintf("%#v", b)
}


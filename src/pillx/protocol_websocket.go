package pillx

import (
    "bytes"
    "encoding/binary"
    "io"
    "log"
)

type WebSocketHeader struct {
    OpcodeByte  byte
    ProloadByte byte
}

func (gateway *GateWayProtocol) Analyze(client *Response) (err error) {
    var err error
    header := new(WebSocketHeader)
    gateway.Header = header
    buf := client.conn.buf

    header.OpcodeByte, err = buf.ReadByte()
    fin := header.fin[0] >> 7
    if fin == 0 {

    }

    opcode, err = header.OpcodeByte & 0x0f
    if opcode == 8 {
        log.Print("Connection closed")
        //self.Close()
        break
    }

    header.ProloadByte, err = buf.ReadByte()
    mask := header.ProloadByte >> 7
    proload := header.ProloadByte & 0x7f

    var (
        lengthBytes  []byte
        length      uint64
        l            uint16
        maskKeyBytes []byte
    )

    switch {
    case payload < 126:
        length = uint64(payload)

    case payload == 126:
        lengthBytes = make([]byte, 2)
        bug.Read(lengthBytes)
        binary.Read(bytes.NewReader(lengthBytes), binary.BigEndian, &l)
        length = uint64(l)

    case payload == 127:
        lengthBytes = make([]byte, 8)
        buf.Read(lengthBytes)
        binary.Read(bytes.NewReader(lengthBytes), binary.BigEndian, &length)
    }

    if mask == 1 {
        maskKeyBytes = make([]byte, 4)
        buf.Read(maskKeyBytes)
    }

    buf = make([]byte, length)
    io.ReadFull(self.Conn, buf)
}
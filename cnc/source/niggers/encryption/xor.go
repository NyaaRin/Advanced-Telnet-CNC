package encryption

import (
	"fmt"
)

var tableKey = `rEFdypQFiezEYXLAejy3fFXDVirHL0yWdnk6hynYRweP6RgdDJig4zhrGcmUBedtHTXGvnE9jm22up2mnZx2FwFiZYzmqW8Nk0kfmykkjDWZL7cYRuMBvPT5iRZe5HmPfp2xi0ZGRppjQGvz9TF3mQKUFiExPPWZ6RC4MQZNPDa7rTP4ZWzcqGivDF6MPjYx4hh6GwWDH4RJt0UjVhyq7QxNSKejeghapi9KcSHzHjw7HjEGJf1XiC3n0yphEYUF`

func xorCrypt(input, key string) (output string) {
	for i := 0; i < len(input); i++ {
		output += string(input[i] ^ key[i%len(key)])
	}

	return output
}

func byteToInt(b byte) int {
	return int(b)
}

func Crypt(input string) {
	encryptedString := xorCrypt(input, tableKey)
	fmt.Printf("XOR'ing %d bytes: ", len(encryptedString))
	for _, b := range encryptedString {
		fmt.Printf("\\x%02X", b)
	}
	fmt.Println()
}

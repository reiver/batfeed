package libfeeds

import (
	"math/rand"
)

var symbols = []byte{'0','1','2','3','4','5','6','7','8','9','A','B','C','D','E','F','G','H','I','J','K','L','M','N','O','P','Q','R','S','T','U','V','W','X','Y','Z'}

// randstr returns a 40 byte random string made up of number-symbols and uppercase-letters.
func randstr() string {

	var buffer [40]byte

	var lenBuffer int = len(buffer)
	var lenSymbols int = len(symbols)

	for index:=0; index<lenBuffer; index++ {
		var symbol byte = symbols[rand.Intn(lenSymbols)]
		buffer[index] = symbol
	}

	return string(buffer[:])
}

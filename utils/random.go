package utils

import (
	"math/rand"
	"time"
)

func RandomBoolean() bool {
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(2)
	return randomInt == 1
}

func RandomUrl() string {
	return "https://" + RandomCharacter(RandomInt(3, 13), false) + ".net"
}

func RandomCharacter(length int, capital bool) string {
	start := int('a')
	stop := int('z')

	bytes := make([]byte, length)

	if capital {
		start = int('A')
		stop = int('Z')
	}

	for i := 0; i < length; i++ {
		bytes[i] = byte(RandomInt(start, stop))
	}

	return string(bytes)
}

func RandomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())

	return min + rand.Intn(max-min+1)
}

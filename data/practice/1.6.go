package practice

import (
	"fmt"
	"strings"
)

/*
 Implement a method to perform basic string compression using the counts of repeated chartacters.

 aabcccccaaa == a2b1c5a3
 Only return the compressed string if its smaller.
 Assume its only upper and lower case chars.

*/

func Compress(s string) (string, error) {
	b := strings.Builder{}
	var lastChar rune
	counter := 0

	for _, v := range s {
		if v != lastChar {
			// Countner should only be 0 on the first iteration, but in general
			// we should only encode non-zero numbers.
			if counter > 0 {
				_, err := b.WriteString(fmt.Sprintf("%c%d", lastChar, counter))
				if err != nil {
					return "", err
				}
			}
			lastChar = v
			counter = 1
			continue
		}
		counter++
	}

	if counter > 0 {
		_, err := b.WriteString(fmt.Sprintf("%c%d", lastChar, counter))
		if err != nil {
			return "", err
		}
	}

	if b.Len() < len(s) {
		return b.String(), nil
	}
	return s, nil
}

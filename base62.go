package pksuid

var base62Map = [256]byte{
	'0': '0', '5': '5', 'A': 'A', 'F': 'F', 'K': 'K', 'P': 'P',
	'1': '1', '6': '6', 'B': 'B', 'G': 'G', 'L': 'L', 'Q': 'Q',
	'2': '2', '7': '7', 'C': 'C', 'H': 'H', 'M': 'M', 'R': 'R',
	'3': '3', '8': '8', 'D': 'D', 'I': 'I', 'N': 'N', 'S': 'S',
	'4': '4', '9': '9', 'E': 'E', 'J': 'J', 'O': 'O', 'T': 'T',
	'U': 'U', 'a': 'a', 'g': 'g', 'm': 'm', 's': 's', 'y': 'y',
	'V': 'V', 'b': 'b', 'h': 'h', 'n': 'n', 't': 't', 'z': 'z',
	'W': 'W', 'c': 'c', 'i': 'i', 'o': 'o', 'u': 'u',
	'X': 'X', 'd': 'd', 'j': 'j', 'p': 'p', 'v': 'v',
	'Y': 'Y', 'e': 'e', 'k': 'k', 'q': 'q', 'w': 'w',
	'Z': 'Z', 'f': 'f', 'l': 'l', 'r': 'r', 'x': 'x',
}

func isBase62Bytes(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	for _, c := range b {
		if base62Map[c] == 0 {
			return false
		}
	}
	return true
}

package utilities

var hexTable = []byte{
	'0', '1', '2', '3',
	'4', '5', '6', '7',
	'8', '9', 'A', 'B',
	'C', 'D', 'E', 'F',
}

func Int64Hex(i int64) string {
	c := 0xF

	tmp := i
	result := make([]byte, 0, 16)
	for shift := 60; shift >= 0; shift -= 4 {
		result = append(result, hexTable[c&int(tmp>>shift)])
	}

	return string(result)
}

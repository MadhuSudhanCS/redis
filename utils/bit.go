package utils

func appendByte(slice []byte, size int) []byte {
	if size >= len(slice) { // if necessary, reallocate
		newSlice := make([]byte, size+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	return slice
}

func PutBit(slice []byte, offset, value int) ([]byte, uint8) {
	index := offset / 8
	newOffset := uint(7 - offset%8)
	slice = appendByte(slice, index)

	oldBitValue := GetBit(slice, offset)

	if value == 1 {
		slice[index] = setBit(slice[index], newOffset)
	}

	if value == 0 {
		slice[index] = clearBit(slice[index], newOffset)
	}

	return slice, oldBitValue
}

func setBit(n byte, pos uint) byte {
	intValue := int(n)
	intValue |= (1 << pos)
	return byte(intValue)
}

func clearBit(n byte, pos uint) byte {
	intValue := int(n)
	mask := ^(1 << pos)
	intValue &= mask
	return byte(intValue)
}

func GetBit(slice []byte, offset int) uint8 {
	index := offset / 8
	newOffset := uint(7 - offset%8)

	if index >= len(slice) {
		return 0
	}

	return slice[index] & (1 << newOffset) >> newOffset
}

package api

// ASCIIDecode to decode the STX adn ETX in packets
func ASCIIDecode(ascii []uint8) interface{} {
	str := ""
	//fmt.Println(ascii)
	for _, v := range ascii {
		if v == 0x02 {
			str += "[STX]"
		} else if v == 0x03 {
			str += "[ETX]"
		} else {
			str += string(v)
		}
	}
	return str
}

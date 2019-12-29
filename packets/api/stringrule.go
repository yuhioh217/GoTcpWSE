package api

import "strings"

// StringRule define the string type in packets
func StringRule(str string) (string, string) {
	if strings.Contains(str, "[ETX]0000[ETX]") {
		return "id", strings.Split(str, "[ETX]0000[ETX]")[1]
	}

	substr := strings.Split(str, "[ETX]")
	switch substr[0] {
	case "0124":
		return "best_five_left", substr[1]
	case "0125":
		return "best_five_right", substr[1]
	case "0102", "0602":
		return "timestamp", substr[1]
	case "0119", "0619":
		return "type", substr[1]
	case "0104", "0604":
		return "deal_price", substr[1]
	case "0114", "0614":
		return "order_count", substr[1]
	case "0113", "0613":
		return "total_count", substr[1]
	case "0115", "0615":
		return "total_amount", substr[1]
	default:
		return "", ""
	}
}

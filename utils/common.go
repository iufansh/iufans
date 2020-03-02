package utils

import "strconv"

// ID转化为推广码
func GenInviteCode(id int64) string {
	if id == 0 {
		return "999"
	}
	return strconv.FormatInt(id*3 + 3051, 10)
}

// 推广码转化为ID
func ReverseInviteCode(ic string) int64 {
	icInt, err := strconv.ParseInt(ic, 10, 64)
	if err != nil {
		return -1
	}
	if icInt = icInt - 3051; icInt < 3 {
		return 0
	}
	if icInt % 3 != 0 {
		return 0
	}
	return icInt / 3
}

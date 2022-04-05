package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// StringToInt 字符串转int
func StringToInt(string string) int {
	parseInt, err := strconv.Atoi(string)
	if err != nil {
		return 0
	}
	return parseInt
}

// StringToFloat64 字符串转float64
func StringToFloat64(string string) float64 {
	parseInt, err := strconv.ParseFloat(string, 64)
	if err != nil {
		return 0
	}
	return parseInt
}

// StringToInt64 字符串转int64
func StringToInt64(string2 string) int64 {
	parseInt, err := strconv.ParseInt(string2, 10, 64)
	if err != nil {
		return 0
	}
	return parseInt
}

// StringToBool 字符串转bool
func StringToBool(string2 string) bool {
	boo, err := strconv.ParseBool(string2)
	if err != nil {
		return false
	}
	return boo
}

// Int64ToString int64转字符串
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

// RemoveReplicaSliceString slice(string类型)元素去重
func RemoveReplicaSliceString(slc []string) []string {
	var result []string
	tempMap := make(map[string]bool)
	for _, e := range slc {
		if _, ok := tempMap[e]; !ok {
			tempMap[e] = true
			result = append(result, e)
		}
	}
	return result
}

// Md5 将字符串生成MD5
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	cipher := h.Sum(nil)
	return hex.EncodeToString(cipher)
}

func String2TimeDuration(str string) (time.Duration, error) {
	supportsSuffix := []string{"s", "m", "h", "w", "d"}
	for _, suffix := range supportsSuffix {
		if strings.HasSuffix(str, suffix) {
			num := str[0 : len(str)-len(suffix)]
			numInt := StringToInt64(num)
			switch suffix {
			case "s":
				return time.Duration(numInt) * time.Second, nil
			case "m":
				return time.Duration(numInt) * time.Minute, nil
			case "h":
				return time.Duration(numInt) * time.Hour, nil
			case "w":
				return time.Duration(numInt) * 7 * 24 * time.Hour, nil
			case "d":
				return time.Duration(numInt) * 24 * time.Hour, nil
			default:
			}
		}
	}
	return 0, errors.New(fmt.Sprintf("'%s' conver to time.Duration fail,unsupported type,support type%v", str, supportsSuffix))
}

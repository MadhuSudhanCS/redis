package cmd

import (
	"bytes"
	"fmt"
)

//RedisProtocol Error Msg
func GetErrMsg(msg interface{}) []byte {
	errString := "-%s\r\n"
	return []byte(fmt.Sprintf(errString, msg))
}

//RedisProtocol Simple String Msg
func GetSimpleString(msg string) []byte {
	str := "+%s\r\n"
	return []byte(fmt.Sprintf(str, msg))
}

//RedisProtocol NULL Msg
func GetNullReply() []byte {
	return []byte("$-1\r\n")
}

//RedisProtocol Bulk String
func GetBulkString(msg string) []byte {
	str := "$%d\r\n%s\r\n"
	return []byte(fmt.Sprintf(str, len(msg), msg))
}

//RedisProtocol Integer Msg
func GetIntegerMsg(value interface{}) []byte {
	str := ":%d\r\n"
	return []byte(fmt.Sprintf(str, value))
}

//RedisProtocol ArrayMsg
func GetArrayMsg(values []interface{}) []byte {
	var buffer bytes.Buffer
	length := "*%d\r\n"

	buffer.WriteString(fmt.Sprintf(length, len(values)))
	for _, value := range values {
		switch t := value.(type) {
		case string:
			buffer.WriteString(string(GetBulkString(t)))
		case int:
			buffer.WriteString(string(GetIntegerMsg(t)))
		}
	}

	return []byte(buffer.String())
}

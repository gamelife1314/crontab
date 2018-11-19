
package common

import (
	"encoding/json"
	"strings"
)

// ExtractWorkerIp extract ip from worker key
func ExtractWorkerIp(rawIp string) string {
	return strings.TrimPrefix(rawIp, CronWorkerDir)
}

// BuildResponse build http response
func BuildResponse(errorNum int, msg string, data interface{}) (resp []byte, err error) {
	var (
		response Response
	)

	response.ErrorCode = errorNum
	response.ErrorMsg = msg
	response.Data = data

	resp, err = json.Marshal(response)

	return
}

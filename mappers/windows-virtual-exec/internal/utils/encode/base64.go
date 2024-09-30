package encode

import (
	"encoding/base64"

	"k8s.io/klog/v2"
)

func DecodeBase64(encodedString string) string {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		klog.Error("Decode error: ", err)
		return ""
	}
	return string(decodedBytes)
}

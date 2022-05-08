package ghwebhook

import (
	"fmt"
	"regexp"
)

func ReplaceVersion(content []byte, regex, newVersion string) []byte {

	r, err := regexp.Compile(regex)
	if err != nil {
		fmt.Printf("Error compiling regex %v: %v\n", regex, err)
		return content
	}

	ret := r.ReplaceAll(content, []byte(fmt.Sprintf("${key}%s${other}", newVersion)))
	// fmt.Println(string(ret))
	return ret
}

// convenience method to ease testing
func ReplaceVersionString(content, regex, newVersion string) string {
	return string(ReplaceVersion([]byte(content), regex, newVersion))
}

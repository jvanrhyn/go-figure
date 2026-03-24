package figure

import (
	"errors"
	"strconv"
	"strings"
)

const signature = "flf2"
const reverseFlag = "1"

var charDelimiters = [3]string{"@", "#", "$"}
var hardblanksBlacklist = [2]byte{'a', '2'}

var errMalformedHeader = errors.New("malformed FIGlet header")

func getHeight(fields []string) (int, error) {
	if len(fields) < 2 {
		return 0, errMalformedHeader
	}
	h, err := strconv.Atoi(fields[1])
	if err != nil || h <= 0 {
		return 0, errMalformedHeader
	}
	return h, nil
}

func getBaseline(fields []string) (int, error) {
	if len(fields) < 3 {
		return 0, errMalformedHeader
	}
	b, err := strconv.Atoi(fields[2])
	if err != nil || b <= 0 {
		return 0, errMalformedHeader
	}
	return b, nil
}

func getHardblank(fields []string) byte {
	if len(fields) < 1 || len(fields[0]) == 0 {
		return ' '
	}
	hardblank := fields[0][len(fields[0])-1]
	if hardblank == hardblanksBlacklist[0] || hardblank == hardblanksBlacklist[1] {
		return ' '
	}
	return hardblank
}

func getReverse(fields []string) bool {
	return len(fields) > 6 && fields[6] == reverseFlag
}

func lastCharLine(text string, height int) bool {
	suffixLen := 2
	if height == 1 && len(text) > 0 {
		suffixLen = 1
	}
	for _, delim := range charDelimiters {
		if strings.HasSuffix(text, strings.Repeat(delim, suffixLen)) {
			return true
		}
	}
	return false
}

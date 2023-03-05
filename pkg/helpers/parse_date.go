package helpers

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	regexps = []string{
		"^(?P<day>\\d{2})-(?P<month>\\d{2})-(?P<year>\\d{4}) (?P<hour>\\d{2}):(?P<min>\\d{2})",
		"^(?P<day>\\d{2})(\\.|:|-|/)(?P<month>\\d{2}) (?P<hour>\\d{2}):(?P<min>\\d{2})",
		"^(?P<hour>\\d{2}):(?P<min>\\d{2})",
		"^(?P<year>\\d{4})-(?P<month>\\d{2})-(?P<day>\\d{2}) (?P<hour>\\d{2}):(?P<min>\\d{2})",
		"^(?P<rel>tomorrow) (?P<hour>\\d{2}):(?P<min>\\d{2})",
	}
	compiledRegexps             []*regexp.Regexp
	compiled                    = false
	ParseDateNoDateError        = errors.New("prefix not found")
	ParseDateInvalidFormatError = errors.New("invalid format")
)

func getRegexps() []*regexp.Regexp {
	if compiled {
		return compiledRegexps
	}
	compiledRegexps = make([]*regexp.Regexp, 0)
	for _, rgxp := range regexps {
		if compile, err := regexp.Compile(rgxp); err == nil {
			compiledRegexps = append(compiledRegexps, compile)
		}
	}
	compiled = true
	return compiledRegexps
}

func ParseDate(prefix string, str string, timeNow time.Time) (time.Time, int, error) {
	var err error
	tIdx := -1
	if prefix != "" {
		tIdx = strings.Index(str, prefix+" ")
		if tIdx == -1 {
			return timeNow, 0, ParseDateNoDateError
		}
	}
	tIdx = tIdx + len(prefix) + 1
	strToSearch := str[tIdx:]

	var res int
	var submatches []string

	testRegexps := getRegexps()
	for _, testRegexp := range testRegexps {
		groupNames := testRegexp.SubexpNames()
		if submatches = testRegexp.FindStringSubmatch(strToSearch); submatches == nil {
			continue
		}

		day := timeNow.Day()
		month := int(timeNow.Month())
		year := timeNow.Year()
		hour := timeNow.Hour()
		min := timeNow.Minute()
		patterntGood := true

		for gIdx, submatch := range submatches {
			if gIdx == 0 {
				continue
			}
			name := groupNames[gIdx]

			if name == "day" {
				if res, err = strconv.Atoi(submatch); err == nil && res > 0 && res < 32 {
					day = res
				} else {
					patterntGood = false
					break
				}
			}
			if name == "month" {
				if res, err = strconv.Atoi(submatch); err == nil && res > 0 && res < 13 {
					month = res
				} else {
					patterntGood = false
					break
				}
			}
			if name == "year" {
				if res, err = strconv.Atoi(submatch); err == nil && res > 2019 {
					year = res
				} else {
					patterntGood = false
					break
				}
			}
			if name == "hour" {
				if res, err = strconv.Atoi(submatch); err == nil && res > -1 && res < 24 {
					hour = res
				} else {
					patterntGood = false
					break
				}
			}
			if name == "min" {
				if res, err = strconv.Atoi(submatch); err == nil && res > -1 && res < 60 {
					min = res
				} else {
					patterntGood = false
					break
				}
			}
			if name == "rel" {
				if submatch == "tomorrow" {
					tomorrow := timeNow.Add(24 * time.Hour)
					day, month, year = tomorrow.Day(), int(tomorrow.Month()), tomorrow.Year()
				}
			}
		}

		if patterntGood {
			newTime := time.Date(year, time.Month(month), day, hour, min, 0, 0, timeNow.Location())
			return newTime, tIdx + len(submatches[0]), nil
		}
	}

	return timeNow, 0, ParseDateInvalidFormatError
}

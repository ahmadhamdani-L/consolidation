package helper

import (
	"context"
	"errors"
	"fmt"
	"mcash-finance-console-core/pkg/redis"
	"mcash-finance-console-core/pkg/util/response"
	"net/url"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func ReplaceWholeWord(text string, oldWord string, newWord string) string {
	var patternLength = len(oldWord)
	var textLength = len(text)

	var copyIndex = 0
	var textIndex = 0
	var patternIndex = 0
	var newString strings.Builder
	var lps = computeLPSArray(oldWord)

	for textIndex < textLength {
		if oldWord[patternIndex] == text[textIndex] {
			patternIndex++
			textIndex++
		}
		if patternIndex == patternLength {
			startIndex := textIndex - patternIndex
			endIndex := textIndex - patternIndex + patternLength - 1

			if checkIfWholeWord(text, startIndex, endIndex) {
				if copyIndex != startIndex {
					newString.WriteString(text[copyIndex:startIndex])
				}
				newString.WriteString(newWord)
				copyIndex = endIndex + 1
			}

			patternIndex = 0
			textIndex = endIndex + 1

		} else if textIndex < textLength && oldWord[patternIndex] != text[textIndex] {

			if patternIndex != 0 {
				patternIndex = lps[patternIndex-1]

			} else {
				textIndex = textIndex + 1
			}

		}
	}
	newString.WriteString(text[copyIndex:])

	return newString.String()
}

func computeLPSArray(pattern string) []int {
	var length = 0
	var i = 1
	var patternLength = len(pattern)

	var lps = make([]int, patternLength)

	lps[0] = 0

	for i = 1; i < patternLength; {
		if pattern[i] == pattern[length] {
			length++
			lps[i] = length
			i++

		} else {

			if length != 0 {
				length = lps[length-1]

			} else {
				lps[i] = length
				i++
			}
		}
	}
	return lps
}

func checkIfWholeWord(text string, startIndex int, endIndex int) bool {
	startIndex = startIndex - 1
	endIndex = endIndex + 1

	if (startIndex < 0 && endIndex >= len(text)) ||
		(startIndex < 0 && endIndex < len(text) && isNonWord(text[endIndex])) ||
		(startIndex >= 0 && endIndex >= len(text) && isNonWord(text[startIndex])) ||
		(startIndex >= 0 && endIndex < len(text) && isNonWord(text[startIndex]) && isNonWord(text[endIndex])) {
		return true
	}

	return false
}

func isNonWord(c byte) bool {
	return !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_'))
}

func MultiVersionFilter(queryParam url.Values) ([]int, error) {
	payload := []int{}
	versionFilter := queryParam["version[]"]
	if len(versionFilter) > 0 {
		var version []int
		for _, v := range versionFilter {
			versionInt, err := strconv.Atoi(v)
			if err != nil {
				return []int{}, errors.New("Not Integer Value!")
			}
			version = append(version, versionInt)
		}
		payload = version
	}
	return payload, nil
}

func MultiRoleFilter(queryParam url.Values) ([]int, error) {
	payload := []int{}
	roleFilter := queryParam["role_id[]"]
	if len(roleFilter) > 0 {
		var role []int
		for _, v := range roleFilter {
			roleInt, err := strconv.Atoi(v)
			if err != nil {
				return []int{}, errors.New("Not Integer Value!")
			}
			role = append(role, roleInt)
		}
		payload = role
	}
	return payload, nil
}

func MultiCoaGroupIDFilter(queryParam url.Values) ([]int, error) {
	payload := []int{}
	coaGroupFilter := queryParam["coa_group_id[]"]
	if len(coaGroupFilter) > 0 {
		var coaGroupID []int
		for _, v := range coaGroupFilter {
			coaGroupIDInt, err := strconv.Atoi(v)
			if err != nil {
				return []int{}, errors.New("Not Integer Value!")
			}
			coaGroupID = append(coaGroupID, coaGroupIDInt)
		}
		payload = coaGroupID
	}
	return payload, nil
}

func MultiStatusFilter(queryParam url.Values) ([]int, error) {
	payload := []int{}
	statusFilter := queryParam["status[]"]
	if len(statusFilter) > 0 {
		var status []int
		for _, v := range statusFilter {
			statusInt, err := strconv.Atoi(v)
			if err != nil {
				return []int{}, errors.New("Not Integer Value!")
			}
			status = append(status, statusInt)
		}
		payload = status
	}
	return payload, nil
}

func CompanyValidation(userID int, inputCompany int) bool {
	pass := false

	allCompany, err := redis.RedisClient.Get(context.Background(), fmt.Sprintf("access_all_company:user:%d", userID)).Result()
	// mau is nil nya true atau false tapi kalau dia errornya is nil dia ignore
	if redis.IsNil(err) || !redis.IsNil(err) {

	} else if err != nil {
		return false
	}

	if allCompany == "1" {
		return true
	}

	listCompany, err := redis.RedisClient.LRange(context.Background(), fmt.Sprintf("access_company:user:%d", userID), 0, -1).Result()
	if err != nil && redis.IsNil(err) == true {
		return false
	}

	for _, v := range listCompany {
		intVal, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println(err)
			return false
		}
		if inputCompany == intVal {
			return true
		}
	}

	return pass
}

func ErrorHandler(err error) *response.Error {
	switch err {
	case gorm.ErrRecordNotFound:
		return response.ErrorBuilder(&response.ErrorConstant.NotFound, err)
	case gorm.ErrInvalidValue, gorm.ErrInvalidValueOfLength, gorm.ErrInvalidField, gorm.ErrInvalidData, gorm.ErrModelValueRequired, gorm.ErrRegistered:
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err)
	default:
		return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
}

func AssignAmount(amount *float64) *float64 {
	if amount == nil {
		zero := 0.0
		return &zero
	}
	return amount
}

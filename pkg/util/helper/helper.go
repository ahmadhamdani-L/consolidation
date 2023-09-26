package helper

import "strings"

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

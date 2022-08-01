package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/nyaruka/phonenumbers"
	"github.com/vonmutinda/organono/app/entities"
	"gopkg.in/guregu/null.v3"
)

var nonDigitRe = regexp.MustCompile("[^0-9eE+]")

func ParsePhoneNumber(phoneNumberStr string) (entities.PhoneNumber, error) {

	phoneNumberStr = nonDigitRe.ReplaceAllString(phoneNumberStr, "")

	phoneNumberStr = sanitizeScientificNotation(phoneNumberStr)

	if len(phoneNumberStr) < 7 {
		return entities.PhoneNumber{}, NewErrorWithCode(
			fmt.Errorf("number (%v) is less than 7 digits", phoneNumberStr),
			ErrorCodeInvalidPhone,
			"check phone number length",
		)
	}

	if len(phoneNumberStr) > 15 {
		return entities.PhoneNumber{}, NewErrorWithCode(
			fmt.Errorf("number (%v) is greater than 15 digits", phoneNumberStr),
			ErrorCodeInvalidPhone,
			"check phone number length",
		)
	}

	countryCode, err := getCountryCode(phoneNumberStr)
	if err != nil {
		return entities.PhoneNumber{}, err

	}

	num, err := phonenumbers.Parse(phoneNumberStr, countryCodes[countryCode])
	if err != nil {
		return entities.PhoneNumber{}, NewErrorWithCode(
			fmt.Errorf("invalid phone number"),
			ErrorCodeInvalidPhone,
			"parse phone number and country code",
		)
	}

	ok := phonenumbers.IsValidNumber(num)
	if !ok {
		return entities.PhoneNumber{}, NewErrorWithCode(
			fmt.Errorf("invalid phone number"),
			ErrorCodeInvalidPhone,
			"check if valid number",
		)
	}

	countryCode = strconv.FormatInt(int64(num.GetCountryCode()), 10)

	number := strconv.FormatInt(int64(num.GetNationalNumber()), 10)

	if num.GetItalianLeadingZero() {
		number = "0" + number
	}

	countryCode = "+" + countryCode
	number = strings.TrimPrefix(number, countryCode)

	phoneNumber := entities.PhoneNumber{
		CountryCode: null.StringFrom(countryCode),
		Number:      null.StringFrom(number),
	}

	return phoneNumber, nil
}

func getCountryCode(phoneNumber string) (string, error) {

	valid := false
	countryCode := ""

	if phoneNumber[0:1] == "+" {
		phoneNumber = phoneNumber[1:]
	}

	for i := 0; i < 3; i++ {
		countryCode = phoneNumber[0 : i+1]
		if _, ok := countryCodes[countryCode]; !ok {
			continue
		}

		valid = true
		break
	}

	if !valid {
		return "", NewError(
			errors.New("missing country code"),
			"get phone country code",
		)
	}

	return countryCode, nil
}

func sanitizeScientificNotation(phoneNumber string) string {
	phoneNumber = strings.ToLower(phoneNumber)
	index := strings.Index(phoneNumber, "e")
	if index > -1 {
		phoneNumber = phoneNumber[:strings.Index(phoneNumber, "e")]
	}

	return phoneNumber
}

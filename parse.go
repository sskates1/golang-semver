package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	//MAXLENGTH max size of semvar
	MAXLENGTH = 256
	//MaxUint64 used to calculate max int
	MaxUint64 = 1<<64 - 1
)

func parse(version string, loose string) (*SemVer, error) {
	if len(version) > MAXLENGTH {
		errorMessage := fmt.Sprintf("Version string longer than %d max length", MAXLENGTH)
		err := errors.New(errorMessage)
		return nil, err
	}
	var tester *regexp.Regexp
	regexes := getRegexes()
	if loose == "loose" {
		tester = regexes["LOOSE"]
	} else {
		tester = regexes["FULL"]
	}
	if !testValid(version, tester) {
		err := errors.New("Invalid semver format")
		return nil, err
	}

	fmt.Println(version, loose)
	semver, err := newSemVer(version, loose)
	if err != nil {
		panicIfError(err, "Creating new semver failed")
	}
	return semver, nil
}

func valid(version string, loose string) string {
	semver, err := parse(version, loose)
	if err != nil {
		panicIfError(err, "Parse failed")
	}
	return semver.version()
}

func clean(version string, loose string) string {
	trimmedVersion := strings.TrimSpace(version)
	if strings.Index(trimmedVersion, "v") == 0 {
		trimmedVersion = strings.Replace(trimmedVersion, "v", "", 1)
	}
	semver, err := parse(trimmedVersion, loose)
	if err != nil {
		panicIfError(err, "Parse failed")
	}
	return semver.version()
}

func newSemVer(version string, loose string) (*SemVer, error) {
	semver := new(SemVer)
	semver.raw = strings.TrimSpace(version)
	semver.loose = loose

	regexs := getRegexes()
	var regex *regexp.Regexp
	if len(loose) > 0 {
		regex = regexs["LOOSE"]
	} else {
		regex = regexs["FULL"]
	}
	matched := regex.FindAllString(semver.raw, -1)
	errors := make([]error, 3)
	semver.major, errors[0] = strconv.Atoi(matched[0])
	semver.minor, errors[1] = strconv.Atoi(matched[1])
	semver.patch, errors[2] = strconv.Atoi(matched[2])
	for k, err := range errors {
		if err != nil {
			fmt.Println(k)
			panicIfError(err, "failed string to int conversion")
		}
	}

	return semver, nil
}

func testValid(version string, regex *regexp.Regexp) bool {

	return true
}

package main

import (
	"regexp"
)

func getRegexes() map[string]*regexp.Regexp {
	regValues := make(map[string]string)

	// The following Regular Expressions can be used for tokenizing,
	// validating, and parsing SemVer version strings.

	// ## Numeric Identifier
	// A single `0`, or a non-zero digit followed by zero or more digits.
	regValues["NUMERICIDENTIFIER"] = "0|[1-9]\\d*"
	regValues["NUMERICIDENTIFIERLOOSE"] = "[0-9]+"

	// ## Non-numeric Identifier
	// Zero or more digits, followed by a letter or hyphen, and then zero or
	// more letters, digits, or hyphens.
	regValues["NONNUMERICIDENTIFIER"] = "\\d*[a-zA-Z-][a-zA-Z0-9-]*"
	regValues["NONNUMERICIDENTIFIER"] = "\\d*[a-zA-Z-][a-zA-Z0-9-]*"

	// ## Main Version
	// Three dot-separated numeric identifiers.
	regValues["MAINVERSION"] = "(" + regValues["NUMERICIDENTIFIER"] + ")\\." +
		"(" + regValues["NUMERICIDENTIFIER"] + ")\\." +
		"(" + regValues["NUMERICIDENTIFIER"] + ")"
	regValues["MAINVERSIONLOOSE"] = "(" + regValues["NUMERICIDENTIFIERLOOSE"] + ")\\." +
		"(" + regValues["NUMERICIDENTIFIERLOOSE"] + ")\\." +
		"(" + regValues["NUMERICIDENTIFIERLOOSE"] + ")"

	// ## Pre-release Version Identifier
	// A numeric identifier, or a non-numeric identifier.
	regValues["PRERELEASEIDENTIFIER"] = "(?:" + regValues["NUMERICIDENTIFIER"] +
		"|" + regValues["NONNUMERICIDENTIFIER"] + ")"
	regValues["PRERELEASEIDENTIFIERLOOSE"] = "(?:" + regValues["NUMERICIDENTIFIERLOOSE"] +
		"|" + regValues["NONNUMERICIDENTIFIER"] + ")"

	// ## Pre-release Version
	// Hyphen, followed by one or more dot-separated pre-release version
	// identifiers.
	regValues["PRERELEASE"] = "(?:-(" + regValues["PRERELEASEIDENTIFIER"] +
		"(?:\\." + regValues["PRERELEASEIDENTIFIER"] + ")*))"
	regValues["PRERELEASELOOSE"] = "(?:-?(" + regValues["PRERELEASEIDENTIFIERLOOSE"] +
		"(?:\\." + regValues["PRERELEASEIDENTIFIERLOOSE"] + ")*))"

	// ## Build Metadata Identifier
	// Any combination of digits, letters, or hyphens.
	regValues["BUILDIDENTIFIER"] = "[0-9A-Za-z-]+"

	// ## Build Metadata
	// Plus sign, followed by one or more period-separated build metadata
	// identifiers.
	regValues["BUILD"] = "(?:\\+(" + regValues["BUILDIDENTIFIER"] +
		"(?:\\." + regValues["BUILDIDENTIFIER"] + ")*))"

	// ## Full Version String
	// A main version, followed optionally by a pre-release version and
	// build metadata.

	// Note that the only major, minor, patch, and pre-release sections of
	// the version string are capturing groups.  The build metadata is not a
	// capturing group, because it should not ever be used in version
	// comparison.

	FULLPLAIN := "v?" + regValues["MAINVERSION"] +
		regValues["PRERELEASE"] + "?" +
		regValues["BUILD"] + "?"

	regValues["FULL"] = "^" + FULLPLAIN + "$"

	// like full, but allows v1.2.3 and =1.2.3, which people do sometimes.
	// also, 1.0.0alpha1 (prerelease without the hyphen) which is pretty
	// common in the npm registry.
	LOOSEPLAIN := "[v=\\s]*" + regValues["MAINVERSIONLOOSE"] +
		regValues["PRERELEASELOOSE"] + "?" +
		regValues["BUILD"] + "?"

	regValues["LOOSE"] = "^" + LOOSEPLAIN + "$"

	regValues["GTLT"] = "((?:<|>)?=?)"

	// Something like "2.*" or "1.2.x".
	// Note that "x.x" is a valid xRange identifer, meaning "any version"
	// Only the first item is strictly required.
	regValues["XRANGEIDENTIFIERLOOSE"] = regValues["NUMERICIDENTIFIERLOOSE"] + "|x|X|\\*"
	regValues["XRANGEIDENTIFIER"] = regValues["NUMERICIDENTIFIER"] + "|x|X|\\*"

	regValues["XRANGEPLAIN"] = "[v=\\s]*(" + regValues["XRANGEIDENTIFIER"] + ")" +
		"(?:\\.(" + regValues["XRANGEIDENTIFIER"] + ")" +
		"(?:\\.(" + regValues["XRANGEIDENTIFIER"] + ")" +
		"(?:" + regValues["PRERELEASE"] + ")?" +
		regValues["BUILD"] + "?" +
		")?)?"

	regValues["XRANGEPLAINLOOSE"] = "[v=\\s]*(" + regValues["XRANGEIDENTIFIERLOOSE"] + ")" +
		"(?:\\.(" + regValues["XRANGEIDENTIFIERLOOSE"] + ")" +
		"(?:\\.(" + regValues["XRANGEIDENTIFIERLOOSE"] + ")" +
		"(?:" + regValues["PRERELEASELOOSE"] + ")?" +
		regValues["BUILD"] + "?" +
		")?)?"

	regValues["XRANGE"] = "^" + regValues["GTLT"] + "\\s*" + regValues["XRANGEPLAIN"] + "$"
	regValues["XRANGELOOSE"] = "^" + regValues["GTLT"] + "\\s*" + regValues["XRANGEPLAINLOOSE"] + "$"

	// Tilde ranges.
	// Meaning is "reasonably at or greater than"
	regValues["LONETILDE"] = "(?:~>?)"

	regValues["TILDETRIM"] = "(\\s*)" + regValues["LONETILDE"] + "\\s+"
	// re["TILDETRIM"] = new RegExp(regValues["TILDETRIM"], "g")
	// tildeTrimReplace := "$1~"

	regValues["TILDE"] = "^" + regValues["LONETILDE"] + regValues["XRANGEPLAIN"] + "$"
	regValues["TILDELOOSE"] = "^" + regValues["LONETILDE"] + regValues["XRANGEPLAINLOOSE"] + "$"

	// Caret ranges.
	// Meaning is "at least and backwards compatible with"
	regValues["LONECARET"] = "(?:\\^)"

	regValues["CARETTRIM"] = "(\\s*)" + regValues["LONECARET"] + "\\s+"
	// re[CARETTRIM] = new RegExp(regValues["CARETTRIM"], "g");
	// var caretTrimReplace = "$1^"

	regValues["CARET"] = "^" + regValues["LONECARET"] + regValues["XRANGEPLAIN"] + "$"
	regValues["CARETLOOSE"] = "^" + regValues["LONECARET"] + regValues["XRANGEPLAINLOOSE"] + "$"

	// A simple gt/lt/eq thing, or just "" to indicate "any version"
	regValues["COMPARATORLOOSE"] = "^" + regValues["GTLT"] + "\\s*(" + LOOSEPLAIN + ")$|^$"
	regValues["COMPARATOR"] = "^" + regValues["GTLT"] + "\\s*(" + FULLPLAIN + ")$|^$"

	// An expression to strip any whitespace between the gtlt and the thing
	// it modifies, so that `> 1.2.3` ==> `>1.2.3`
	regValues["COMPARATORTRIM"] = "(\\s*)" + regValues["GTLT"] +
		"\\s*(" + LOOSEPLAIN + "|" + regValues["XRANGEPLAIN"] + ")"

	// this one has to use the /g flag
	// re["COMPARATORTRIM"] = new RegExp(regValues["COMPARATORTRIM"], "g")
	// comparatorTrimReplace = "$1$2$3"

	// Something like `1.2.3 - 1.2.4`
	// Note that these all use the loose form, because they"ll be
	// checked against either the strict or loose comparator form
	// later.
	regValues["HYPHENRANGE"] = "^\\s*(" + regValues["XRANGEPLAIN"] + ")" +
		"\\s+-\\s+" +
		"(" + regValues["XRANGEPLAIN"] + ")" +
		"\\s*$"

	regValues["HYPHENRANGELOOSE"] = "^\\s*(" + regValues["XRANGEPLAINLOOSE"] + ")" +
		"\\s+-\\s+" +
		"(" + regValues["XRANGEPLAINLOOSE"] + ")" +
		"\\s*$"

	// Star ranges basically just allow anything at all.
	regValues["STAR"] = "(<|>)?=?\\s*\\*"

	regexes := make(map[string]*regexp.Regexp)
	var err error

	for key, value := range regValues {
		regexes[key], err = regexp.Compile(value)
		panicIfError(err, key)
	}

	return regexes
}

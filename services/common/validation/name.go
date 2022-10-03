package validation

import "regexp"

const (
	segment            string = "[a-z0-9]([-_a-z0-9]*[a-z0-9])?"
	dnsLabelNameRegexp        = "(" + segment + "\\.)*" + segment

	maxNameLength int = 253 // max length of DNS label
)

var sysctlRegexp = regexp.MustCompile("^" + dnsLabelNameRegexp + "$")

// IsValidDnsLabelName checks that the given string is a valid name to us in url (DNS label)
func IsValidDnsLabelName(name string) bool {
	if len(name) > maxNameLength {
		return false
	}
	return sysctlRegexp.MatchString(name)
}

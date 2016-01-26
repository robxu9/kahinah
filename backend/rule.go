package kahinah

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
)

// Rule represents a rule for grouping updates to advisories.
// if a regexp is defined, the update must match the regexp.
// the update must match _all_ defined regexp to return true.
// nil means no match needed - so an empty value Rule will return true for match
type Rule struct {
	RuleName     string         `json:"rule_name"`
	Platform     string         `json:"platform,omitempty"`
	PlatformCut  *regexp.Regexp `json:"platform_cut,omitempty"`
	Name         string         `json:"name,omitempty"`
	NameCut      *regexp.Regexp `json:"name_cut,omitempty"`
	EVR          string         `json:"evr,omitempty"`
	EVRCut       *regexp.Regexp `json:"evr_cut,omitempty"`
	Submitter    string         `json:"submitter,omitempty"`
	SubmitterCut *regexp.Regexp `json:"submitter_cut,omitempty"`
	Date         string         `json:"date,omitempty"`
	DateCut      *regexp.Regexp `json:"date_cut,omitempty"`
	Type         string         `json:"type,omitempty"`
	TypeCut      *regexp.Regexp `json:"type_cut,omitempty"`
	Packages     string         `json:"packages,omitempty"`
	PackagesCut  *regexp.Regexp `json:"packages_cut,omitempty"`
}

// Check if there is any advisory that this updated matches with (via rules)
// If not, returns nil
func (k *K) MatchAdvisory(u *Update) *Advisory {

}

// check if a Rule matches an update with an advisory.
func (r *Rule) Matches(u *Update, a *Advisory) bool {
	if r.Platform != nil && !r.Platform.MatchString(u.Platform) {
		return false
	}
	if r.Name != nil && !r.Name.MatchString(u.Name) {
		return false
	}
	if r.EVR != nil && !r.EVR.MatchString(u.EVR) {
		return false
	}
	if r.Submitter != nil && !r.Submitter.MatchString(u.Submitter) {
		return false
	}
	if r.Date != nil && !r.Date.MatchString(strconv.FormatInt(u.Date.Unix(), 10)) {
		return false
	}
	if r.Type != nil && !r.Type.MatchString(u.Type) {
		return false
	}
	if r.Packages != nil {
		found := false
		for _, v := range u.Packages {
			for _, p := range v {
				if r.Packages.MatchString(p) {
					found = true
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// readRules will parse rules into the connector, erroring as necessary
func (k *K) ReadRules(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	k.Rules = []*Rule{}
	names := map[string]bool{}
	var lastError error

	for _, v := range files {
		if v.IsDir() {
			log.Printf("[rules] skipping %v - we don't iterate into directories", v.Name())
			continue
		}

		// open the path
		file, err := ioutil.ReadFile(filepath.Join(path, v.Name()))
		if err != nil {
			log.Printf("[rules] skipping %v - can't open: %v", v.Name(), err)
			lastError = err
			continue
		}

		parsedRules := []*Rule{}
		if err = json.Unmarshal(file, &parsedRules); err != nil {
			log.Printf("[rules] skipping %v - can't unmarshal: %v", v.Name(), err)
			lastError = err
			continue
		}

		// validate JSON
		for _, parsed := range parsedRules {
			if names[parsed.RuleName] {
				log.Printf("[rules] skipping %v in %v, duplicate name", parsed.RuleName, v.Name())
				continue
			}
			names[parsed.RuleName] = true
		}

		k.Rules = append(k.Rules, parsedRules...)

		log.Printf("[rules] added %v to the ruleset", v.Name())
	}

	return lastError
}

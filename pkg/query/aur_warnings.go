package query

import (
	"strings"

	"github.com/leonelquinteros/gotext"

	"github.com/Jguer/aur"
	"github.com/Jguer/go-alpm/v2"

	"github.com/Jguer/yay/v12/pkg/stringset"
	"github.com/Jguer/yay/v12/pkg/text"
)

type AURWarnings struct {
	Orphans   []string
	OutOfDate []string
	Missing   []string
	Ignore    stringset.StringSet

	log *text.Logger
}

func NewWarnings(logger *text.Logger) *AURWarnings {
	if logger == nil {
		logger = text.GlobalLogger
	}
	return &AURWarnings{Ignore: make(stringset.StringSet), log: logger}
}

func (warnings *AURWarnings) AddToWarnings(remote map[string]alpm.IPackage, aurPkg *aur.Pkg) {
	name := aurPkg.Name
	pkg, ok := remote[name]
	if !ok {
		return
	}

	if aurPkg.Maintainer == "" && !pkg.ShouldIgnore() {
		warnings.Orphans = append(warnings.Orphans, name)
	}

	if aurPkg.OutOfDate != 0 && !pkg.ShouldIgnore() {
		warnings.OutOfDate = append(warnings.OutOfDate, name)
	}
}

func (warnings *AURWarnings) CalculateMissing(remoteNames []string, remote map[string]alpm.IPackage, aurData map[string]*aur.Pkg) {
	for _, name := range remoteNames {
		if _, ok := aurData[name]; !ok && !remote[name].ShouldIgnore() {
			warnings.Missing = append(warnings.Missing, name)
		}
	}
}

func (warnings *AURWarnings) Print() {
	normalMissing, debugMissing := filterDebugPkgs(warnings.Missing)

	if len(normalMissing) > 0 {
		warnings.log.Warnln(gotext.Get("Packages not in AUR:"), formatNames(normalMissing))
	}

	if len(debugMissing) > 0 {
		warnings.log.Warnln(gotext.Get("Missing AUR Debug Packages:"), formatNames(debugMissing))
	}

	if len(warnings.Orphans) > 0 {
		warnings.log.Warnln(gotext.Get("Orphan (unmaintained) AUR Packages:"), formatNames(warnings.Orphans))
	}

	if len(warnings.OutOfDate) > 0 {
		warnings.log.Warnln(gotext.Get("Flagged Out Of Date AUR Packages:"), formatNames(warnings.OutOfDate))
	}
}

func filterDebugPkgs(names []string) (normal, debug []string) {
	normal = make([]string, 0, len(names))
	debug = make([]string, 0, len(names))

	for _, name := range names {
		if strings.HasSuffix(name, "-debug") {
			debug = append(debug, name)
		} else {
			normal = append(normal, name)
		}
	}

	return
}

func formatNames(names []string) string {
	return " " + text.Cyan(strings.Join(names, "  "))
}

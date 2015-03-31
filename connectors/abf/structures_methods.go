package abf

import (
	"fmt"
	"time"

	"github.com/robxu9/kahinah/kahinah"
)

const (
	PackageTypeSource = "source"
	PackageTypeBinary = "binary"
)

func (b *BuildList) Name() string {
	for _, v := range b.BuildList.Packages {
		if v.Type == PackageTypeSource {
			if v.Epoch == 0 {
				return fmt.Sprintf("%s %s-%s.%s", v.Name, v.Version, v.Release, b.BuildList.Arch.Name)
			}
			return fmt.Sprintf("%s %d:%s-%s.%s", v.Name, v.Epoch, v.Version, v.Release, b.BuildList.Arch.Name)
		}
	}

	if b.BuildList.PackageVersion == "" {
		// wut
		return fmt.Sprintf("%s unknown.%s", b.BuildList.Project.Name, b.BuildList.Arch.Name)
	}
	return fmt.Sprintf("%s %s.%s", b.BuildList.Project.Name, b.BuildList.PackageVersion, b.BuildList.Arch.Name)
}

func (b *BuildList) Type() kahinah.UpdateType {
	switch b.BuildList.UpdateType {
	case "security":
		return kahinah.SECURITY
	case "bugfix":
		return kahinah.BUGFIX
	case "enhancement":
		return kahinah.ENHANCEMENT
	case "recommended":
		return kahinah.BUGFIX
	case "newpackage":
		return kahinah.NEW
	default:
		return kahinah.NONE
	}
}

func (b *BuildList) Packages() []*kahinah.UpdatePackage {
	pkgs := make([]*kahinah.UpdatePackage, len(b.BuildList.Packages))

	for k, v := range b.BuildList.Packages {
		pkgs[k] = &kahinah.UpdatePackage{
			Name:    v.Name,
			Epoch:   uint64(v.Epoch),
			Version: v.Version,
			Release: v.Release,
			Arch:    b.BuildList.Arch.Name,
			Url:     v.URL,
		}

		switch v.Type {
		case PackageTypeSource:
			pkgs[k].Type = "src"
		default:
			pkgs[k].Type = "rpm"
		}
	}

	return pkgs
}

func (b *BuildList) Changes() []*kahinah.UpdateChange {
	// FIXME not implemented yet, requires fixing changelog on server side
	return []*kahinah.UpdateChange{
		&kahinah.UpdateChange{
			ChangeAt: time.Unix(b.BuildList.CreatedAt, 0),
			For:      b.Name(),
			By:       "maintainer: " + b.BuildList.Project.Maintainer.Uname,
			Details:  fmt.Sprintf("see diff at %s/diff/%s...%s", b.BuildList.Project.URL, b.BuildList.LastPublishedCommitHash, b.BuildList.CommitHash),
		},
	}
}

package abf

import (
	"fmt"
	"strings"
)

func (a *APIBuildList) Name() string {
	return a.BuildList.Project.Name
}

func (w *WebBuildList) Name() string {
	return w.BuildList.ItemGroups.Group[0].Name
}

func (a *APIBuildList) EVR() string {
	for _, v := range a.BuildList.Packages {
		if v.Type == "source" {
			if v.Epoch == 0 {
				return fmt.Sprintf("%s-%s", v.Version, v.Release)
			}
			return fmt.Sprintf("%d:%s-%s", v.Epoch, v.Version, v.Release)
		}
	}

	return a.BuildList.PackageVersion
}

func (w *WebBuildList) EVR() string {
	for _, v := range w.BuildList.Packages {
		if strings.Contains(v.Fullname, ".src.rpm") {
			if v.Epoch == 0 {
				return fmt.Sprintf("%s-%s", v.Version, v.Release)
			}
			return fmt.Sprintf("%d:%s-%s", v.Epoch, v.Version, v.Release)
		}
	}

	return w.BuildList.ItemGroups.Group[0].Path.Text
}

func (w *WebBuildList) Packages() []string {
	pkgs := make([]string, len(w.BuildList.Packages))

	for k, v := range w.BuildList.Packages {
		// XXX: epoches are not in package names
		pkgs[k] = v.Fullname
	}

	return pkgs
}

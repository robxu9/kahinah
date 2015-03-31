package abf

// gojson helped with BuildLists and BuildList, a combination of query + examples

// https://abf.io/api/v1/build_lists.json?per_page=100&filter[status]=0&filter[ownership]=index&filter[platform_id]=$PLATFORM_ID
type BuildLists struct {
	BuildLists []struct {
		ID        int    `json:"id"`
		ProjectID int    `json:"project_id"`
		Status    int    `json:"status"`
		URL       string `json:"url"`
	} `json:"build_lists"`
	URL string `json:"url"`
}

// https://abf.io/api/v1/build_lists/$BUILDLIST_ID.json
type BuildList struct {
	BuildList struct {
		Advisory struct {
			Description string `json:"description"`
			ID          int    `json:"id"`
			Name        string `json:"name"`
			URL         string `json:"url"`
		} `json:"advisory"`
		Arch struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"arch"`
		AutoCreateContainer bool   `json:"auto_create_container"`
		AutoPublishStatus   string `json:"auto_publish_status"`
		BuildForPlatform    struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"build_for_platform"`
		BuildLogURL     string `json:"build_log_url"`
		CommitHash      string `json:"commit_hash"`
		ContainerPath   string `json:"container_path"`
		ContainerStatus int    `json:"container_status"`
		CreatedAt       int64  `json:"created_at"`
		Duration        int    `json:"duration"`
		ExtraBuildLists []struct {
			ContainerPath string `json:"container_path"`
			ID            int    `json:"id"`
			Status        string `json:"status"`
			URL           string `json:"url"`
		} `json:"extra_build_lists"`
		ExtraParams       map[string]interface{} `json:"extra_params"` // FIXME
		ExtraRepositories []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Platform struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"platform"`
			URL string `json:"url"`
		} `json:"extra_repositories"`
		ID           int `json:"id"`
		IncludeRepos []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Platform struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"platform"`
			URL string `json:"url"`
		} `json:"include_repos"`
		LastPublishedCommitHash string `json:"last_published_commit_hash"`
		Logs                    []struct {
			FileName string `json:"file_name"`
			Size     string `json:"size"`
			URL      string `json:"url"`
		} `json:"logs"`
		MassBuild struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"mass_build"`
		NewCore        bool   `json:"new_core"`
		PackageVersion string `json:"package_version"`
		Packages       []struct {
			DependentProjects []struct {
				DependentPackages []string `json:"dependent_packages"`
				Fullname          string   `json:"fullname"`
				GitURL            string   `json:"git_url"`
				ID                int      `json:"id"`
				Maintainer        struct {
					Email string `json:"email"`
					ID    int    `json:"id"`
					Name  string `json:"name"`
					Uname string `json:"uname"`
					URL   string `json:"url"`
				} `json:"maintainer"`
				Name   string `json:"name"`
				SshURL string `json:"ssh_url"`
				URL    string `json:"url"`
			} `json:"dependent_projects"`
			Epoch     int    `json:"epoch"`
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Release   string `json:"release"`
			Type      string `json:"type"`
			UpdatedAt int    `json:"updated_at"`
			URL       string `json:"url"`
			Version   string `json:"version"`
		} `json:"packages"`
		Priority int `json:"priority"`
		Project  struct {
			Fullname   string `json:"fullname"`
			GitURL     string `json:"git_url"`
			ID         int    `json:"id"`
			Maintainer struct {
				Email string `json:"email"`
				ID    int    `json:"id"`
				Name  string `json:"name"`
				Uname string `json:"uname"`
				URL   string `json:"url"`
			} `json:"maintainer"`
			Name   string `json:"name"`
			SshURL string `json:"ssh_url"`
			URL    string `json:"url"`
		} `json:"project"`
		Publisher struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"publisher"`
		SaveBuildroot    bool `json:"save_buildroot"`
		SaveToRepository struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Platform struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"platform"`
			URL string `json:"url"`
		} `json:"save_to_repository"`
		Status          int    `json:"status"`
		UpdateType      string `json:"update_type"`
		UpdatedAt       int64  `json:"updated_at"`
		URL             string `json:"url"`
		UseCachedChroot bool   `json:"use_cached_chroot"`
		UseExtraTests   bool   `json:"use_extra_tests"`
		User            struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Type  string `json:"type"`
			Uname string `json:"uname"`
			URL   string `json:"url"`
		} `json:"user"`
	} `json:"build_list"`
}

// https://abf.io/api/v1/users/$USER.json
type User struct {
	User struct {
		AvatarURL              string `json:"avatar_url"`
		BuildPriority          int    `json:"build_priority"`
		Company                string `json:"company"`
		CreatedAt              int    `json:"created_at"`
		Email                  string `json:"email"`
		HtmlURL                string `json:"html_url"`
		ID                     int    `json:"id"`
		Language               string `json:"language"`
		Location               string `json:"location"`
		Name                   string `json:"name"`
		OwnProjectsCount       int    `json:"own_projects_count"`
		ProfessionalExperience string `json:"professional_experience"`
		Site                   string `json:"site"`
		Uname                  string `json:"uname"`
		UpdatedAt              int    `json:"updated_at"`
		URL                    string `json:"url"`
	} `json:"user"`
}

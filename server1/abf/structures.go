package abf

// gojson helped with BuildLists and BuildList, a combination of query + examples

// https://abf.io/build_lists?per_page=100&filter[status]=0&filter[ownership]=everything&filter[save_to_platform_id]=$PLATFORM_ID
type BuildLists struct {
	BuildLists []struct {
		ArchID                  int    `json:"arch_id"`
		BuildForPlatformID      int    `json:"build_for_platform_id"`
		CommitHash              string `json:"commit_hash"`
		GroupID                 int    `json:"group_id"`
		ID                      int    `json:"id"`
		LastPublishedCommitHash string `json:"last_published_commit_hash"`
		ProjectID               int    `json:"project_id"`
		ProjectVersion          string `json:"project_version"`
		SaveToPlatformID        int    `json:"save_to_platform_id"`
		SaveToRepositoryID      int    `json:"save_to_repository_id"`
		Status                  int    `json:"status"`
		UpdatedAt               string `json:"updated_at"`
		UpdatedAtUtc            string `json:"updated_at_utc"`
		UserID                  int    `json:"user_id"`
		VersionRelease          string `json:"version_release"`
	} `json:"build_lists"`
	Dictionary struct {
		Arches []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"arches"`
		Platforms []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Personal bool   `json:"personal"`
		} `json:"platforms"`
		Projects []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Owner string `json:"owner"`
		} `json:"projects"`
		Repositories []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"repositories"`
		Users []struct {
			Fullname string `json:"fullname"`
			ID       int    `json:"id"`
			Uname    string `json:"uname"`
		} `json:"users"`
	} `json:"dictionary"`
	Filter struct {
		ArchID             interface{} `json:"arch_id"`
		BuildForPlatformID interface{} `json:"build_for_platform_id"`
		CreatedAtEnd       interface{} `json:"created_at_end"`
		CreatedAtStart     interface{} `json:"created_at_start"`
		ID                 interface{} `json:"id"`
		IsCircle           interface{} `json:"is_circle"`
		MassBuildID        interface{} `json:"mass_build_id"`
		NewCore            interface{} `json:"new_core"`
		Ownership          string      `json:"ownership"`
		ProjectName        interface{} `json:"project_name"`
		ProjectVersion     interface{} `json:"project_version"`
		SaveToPlatformID   interface{} `json:"save_to_platform_id"`
		SaveToRepositoryID interface{} `json:"save_to_repository_id"`
		Status             interface{} `json:"status"`
		UpdatedAtEnd       interface{} `json:"updated_at_end"`
		UpdatedAtStart     interface{} `json:"updated_at_start"`
	} `json:"filter"`
	Page         interface{} `json:"page"`
	ServerStatus struct {
		Publish struct {
			BuildTasks   int `json:"build_tasks"`
			DefaultTasks int `json:"default_tasks"`
			LowTasks     int `json:"low_tasks"`
			Tasks        int `json:"tasks"`
			Workers      int `json:"workers"`
		} `json:"publish"`
		Rpm struct {
			BuildTasks   int `json:"build_tasks"`
			DefaultTasks int `json:"default_tasks"`
			LowTasks     int `json:"low_tasks"`
			OtherWorkers int `json:"other_workers"`
			Tasks        int `json:"tasks"`
			Workers      int `json:"workers"`
		} `json:"rpm"`
	} `json:"server_status"`
	TotalItems int `json:"total_items"`
}

// https://abf.io/api/v1/build_lists/$BUILDLIST_ID.json
type APIBuildList struct {
	BuildList struct {
		Advisory interface{} `json:"advisory"`
		Arch     struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"arch"`
		AutoCreateContainer bool   `json:"auto_create_container"`
		AutoPublishStatus   string `json:"auto_publish_status"`
		BuildForPlatform    struct {
			ID           int    `json:"id"`
			Name         string `json:"name"`
			PlatformType string `json:"platform_type"`
			URL          string `json:"url"`
			Visibility   string `json:"visibility"`
		} `json:"build_for_platform"`
		BuildLogURL       string        `json:"build_log_url"`
		CommitHash        string        `json:"commit_hash"`
		ContainerPath     string        `json:"container_path"`
		ContainerStatus   int           `json:"container_status"`
		CreatedAt         int           `json:"created_at"`
		Duration          int           `json:"duration"`
		ExtraBuildLists   []interface{} `json:"extra_build_lists"`
		ExtraParams       struct{}      `json:"extra_params"`
		ExtraRepositories []interface{} `json:"extra_repositories"`
		ID                int           `json:"id"`
		IncludeRepos      []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Platform struct {
				ID           int    `json:"id"`
				Name         string `json:"name"`
				PlatformType string `json:"platform_type"`
				URL          string `json:"url"`
				Visibility   string `json:"visibility"`
			} `json:"platform"`
			URL string `json:"url"`
		} `json:"include_repos"`
		LastPublishedCommitHash string `json:"last_published_commit_hash"`
		Logs                    []struct {
			FileName string  `json:"file_name"`
			Size     float64 `json:"size"`
			URL      string  `json:"url"`
		} `json:"logs"`
		MassBuild      interface{} `json:"mass_build"`
		NewCore        bool        `json:"new_core"`
		PackageVersion string      `json:"package_version"`
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
				Name       string `json:"name"`
				URL        string `json:"url"`
				Visibility string `json:"visibility"`
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
			Name       string `json:"name"`
			URL        string `json:"url"`
			Visibility string `json:"visibility"`
		} `json:"project"`
		Publisher struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Type  string `json:"type"`
			Uname string `json:"uname"`
			URL   string `json:"url"`
		} `json:"publisher"`
		SaveToRepository struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Platform struct {
				ID           int    `json:"id"`
				Name         string `json:"name"`
				PlatformType string `json:"platform_type"`
				URL          string `json:"url"`
				Visibility   string `json:"visibility"`
			} `json:"platform"`
			URL string `json:"url"`
		} `json:"save_to_repository"`
		Status          int    `json:"status"`
		UpdateType      string `json:"update_type"`
		UpdatedAt       int    `json:"updated_at"`
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

// https://abf.io/build_lists/$BUILDLIST_ID
type WebBuildList struct {
	BuildList struct {
		CanCancel                bool   `json:"can_cancel"`
		CanCreateContainer       bool   `json:"can_create_container"`
		CanPublish               bool   `json:"can_publish"`
		CanPublishInFuture       bool   `json:"can_publish_in_future"`
		CanPublishIntoRepository bool   `json:"can_publish_into_repository"`
		CanPublishIntoTesting    bool   `json:"can_publish_into_testing"`
		CanRejectPublish         bool   `json:"can_reject_publish"`
		ContainerPath            string `json:"container_path"`
		ContainerStatus          int    `json:"container_status"`
		DependentProjectsExists  bool   `json:"dependent_projects_exists"`
		ExtraBuildListsPublished bool   `json:"extra_build_lists_published"`
		HumanDuration            string `json:"human_duration"`
		ID                       int    `json:"id"`
		ItemGroups               struct {
			Group []struct {
				Level int    `json:"level"`
				Name  string `json:"name"`
				Path  struct {
					Href string `json:"href"`
					Text string `json:"text"`
				} `json:"path"`
				Status int `json:"status"`
			} `json:"group"`
		} `json:"item_groups"`
		Packages []struct {
			DependentProjects []struct {
				DependentPackages []string `json:"dependent_packages"`
				Name              string   `json:"name"`
				NewURL            string   `json:"new_url"`
				URL               string   `json:"url"`
			} `json:"dependent_projects"`
			Epoch    interface{} `json:"epoch"`
			Fullname string      `json:"fullname"`
			ID       int         `json:"id"`
			Name     string      `json:"name"`
			Release  string      `json:"release"`
			Sha1     string      `json:"sha1"`
			URL      string      `json:"url"`
			Version  string      `json:"version"`
		} `json:"packages"`
		Publisher struct {
			Fullname string `json:"fullname"`
			Path     string `json:"path"`
		} `json:"publisher"`
		Results []struct {
			FileName string  `json:"file_name"`
			Sha1     string  `json:"sha1"`
			Size     float64 `json:"size"`
			URL      string  `json:"url"`
		} `json:"results"`
		Status       int    `json:"status"`
		UpdateType   string `json:"update_type"`
		UpdatedAt    string `json:"updated_at"`
		UpdatedAtUtc string `json:"updated_at_utc"`
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

package pypi

import "time"

type Package struct {
	Info       Info                 `json:"info"`
	LastSerial int                  `json:"last_serial"`
	Releases   map[string][]Release `json:"releases"`
	Urls       []Urls               `json:"urls"`
}

type Downloads struct {
	LastDay   int `json:"last_day"`
	LastMonth int `json:"last_month"`
	LastWeek  int `json:"last_week"`
}

type ProjectUrls struct {
	Download string `json:"Download"`
	Homepage string `json:"Homepage"`
}

type Info struct {
	Author                 string      `json:"author"`
	AuthorEmail            string      `json:"author_email"`
	Classifiers            []string    `json:"classifiers"`
	Description            string      `json:"description"`
	DescriptionContentType string      `json:"description_content_type"`
	DownloadURL            string      `json:"download_url"`
	Downloads              Downloads   `json:"downloads"`
	HomePage               string      `json:"home_page"`
	Keywords               string      `json:"keywords"`
	License                string      `json:"license"`
	Maintainer             string      `json:"maintainer"`
	MaintainerEmail        string      `json:"maintainer_email"`
	Name                   string      `json:"name"`
	PackageURL             string      `json:"package_url"`
	Platform               string      `json:"platform"`
	ProjectURL             string      `json:"project_url"`
	ProjectUrls            ProjectUrls `json:"project_urls"`
	ReleaseURL             string      `json:"release_url"`
	RequiresDist           []string    `json:"requires_dist"`
	RequiresPython         string      `json:"requires_python"`
	Summary                string      `json:"summary"`
	Version                string      `json:"version"`
}

type Digests struct {
	Md5    string `json:"md5"`
	Sha256 string `json:"sha256"`
}

type Release struct {
	CommentText       string    `json:"comment_text"`
	Digests           Digests   `json:"digests"`
	Downloads         int       `json:"downloads"`
	Filename          string    `json:"filename"`
	HasSig            bool      `json:"has_sig"`
	Md5Digest         string    `json:"md5_digest"`
	Packagetype       string    `json:"packagetype"`
	PythonVersion     string    `json:"python_version"`
	Size              int       `json:"size"`
	UploadTime        string    `json:"upload_time"`
	UploadTimeIso8601 time.Time `json:"upload_time_iso_8601"`
	URL               string    `json:"url"`
}

type Urls struct {
	CommentText       string    `json:"comment_text"`
	Digests           Digests   `json:"digests"`
	Downloads         int       `json:"downloads"`
	Filename          string    `json:"filename"`
	HasSig            bool      `json:"has_sig"`
	Md5Digest         string    `json:"md5_digest"`
	Packagetype       string    `json:"packagetype"`
	PythonVersion     string    `json:"python_version"`
	Size              int       `json:"size"`
	UploadTime        string    `json:"upload_time"`
	UploadTimeIso8601 time.Time `json:"upload_time_iso_8601"`
	URL               string    `json:"url"`
}

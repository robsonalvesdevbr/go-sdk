package entity

type GoVersion struct {
	Version string   `json:"version"`
	Stable  bool     `json:"stable"`
	Files   []GoFile `json:"files,omitempty"`
}

type GoFile struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Sha256   string `json:"sha256"`
	Kind     string `json:"kind"`
}

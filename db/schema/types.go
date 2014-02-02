package schema

import (
	"time"
)

type Project struct {
	Id   int
	Name string
	Path string
	MCId string `db:"mcid"`
}

type ProjectEvent struct {
	Id        int
	Path      string
	Event     string
	EventTime time.Time `db:"event_time"`
	ProjectId int       `db:"project_id"`
}

type DataDir struct {
	Id         int
	ProjectID  int    `db:"project_id"`
	MCId       string `db:"mcid"`
	Name       string
	Path       string
	ParentMCId string `db:"parent_mcid"`
	Parent     int
}

type DataFile struct {
	Id         int
	MCId       string    `db:"mcid"`
	DataDirID  int       `db:"datadir_id"`
	ProjectID  int       `db:"project_id"`
	LastUpload time.Time `db:"last_upload"`
	MTime      time.Time `db:"mtime"`
	ParentMCId string    `db:"parent_mcid"`
	Parent     int
	Name       string
	Path       string
	Size       int
	Checksum   string
	Version    int
}

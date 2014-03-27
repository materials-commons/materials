package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/materials-commons/materials/db/model"
	"github.com/materials-commons/materials/db/schema"
)

var (
	// ProjectsModel is the model for projects
	ProjectsModel *model.Model

	// ProjectEventsModel is the model for project events
	ProjectEventsModel *model.Model

	// ProjectFilesModel is the model for project files
	ProjectFilesModel *model.Model

	// Projects is the query model for projects
	Projects *model.Query

	// ProjectEvents is the query model for project events
	ProjectEvents *model.Query

	// ProjectFiles is the query model for project files
	ProjectFiles *model.Query
)

// Use sets the database connection for all the models.
func Use(db *sqlx.DB) {
	Projects = ProjectsModel.Q(db)
	ProjectEvents = ProjectEventsModel.Q(db)
	ProjectFiles = ProjectFilesModel.Q(db)
}

func init() {
	pQueries := model.ModelQueries{
		Insert: "insert into projects (name, path, mcid) values (:name, :path, :mcid)",
	}
	ProjectsModel = model.New(schema.Project{}, "projects", pQueries)

	peQueries := model.ModelQueries{
		Insert: `insert into project_events (path, event, event_time, project_id)
                 values (:path, :event, :event_time, :project_id)`,
	}
	ProjectEventsModel = model.New(schema.ProjectEvent{}, "project_events", peQueries)

	pfQueries := model.ModelQueries{
		Insert: `insert into project_files (path, size, checksum, mtime, isdir, project_id)
                 values (:path, :size, :checksum, :mtime, :isdir, :project_id)`,
	}
	ProjectFilesModel = model.New(schema.ProjectFile{}, "project_files", pfQueries)
}

package models

import (
	"fmt"
	"log"

	"github.com/bashbunni/project-management/utils"
	"gorm.io/gorm"
)

const notFound uint = 0

// Entity
type Project struct {
	gorm.Model
	Name      string
}

// TODO: make all functions return errors for testability

// Create a new project instance.
// DeletedAt defaults to the zero value for time.Time.
func NewProject(id uint, name string) *Project {
	return &Project{Name: name}
}

// Implement list.Item for Bubbletea TUI
// TODO: change this
func (p Project) Title() string       { return p.Name }
func (p Project) Description() string { return fmt.Sprintf("%d", p.ID) }
func (p Project) FilterValue() string { return p.Name }

// Interface
type ProjectRepository interface {
	PrintProjects()
	HasProjects() bool
	GetProjectByID(projectID uint) (Project, error)
	GetAllProjects() ([]Project, error)
	CreateProject(name string) (Project, error)
	DeleteProject(projectID uint, er EntryRepository) error
	RenameProject(projectID uint) error
}

// Gorm implementation
type GormProjectRepository struct {
	DB *gorm.DB
}

func (g *GormProjectRepository) GetProjectByID(projectID uint) (Project, error) {
	var project Project
	if err := g.DB.Where("id = ?", projectID).Find(&project).Error; err != nil {
		return project, utils.ErrProjectNotFound
	}
	return project, nil
}

func (g *GormProjectRepository) PrintProjects() {
	projects := g.GetAllProjects()
	for _, project := range projects {
		fmt.Printf(Format, project.ID, project.Name)
	}
}

func (g *GormProjectRepository) GetAllProjects() ([]Project, error) {
	var projects []Project
	if err := g.DB.Find(&projects).Error; err != nil {
		return projects, utils.ErrEmptyTable
	}
	return projects, nil
}

func (g *GormProjectRepository) HasProjects() bool {
	if len(g.GetAllProjects()) == 0 {
		return false
	}
	return true
}


func (g *GormProjectRepository) CreateProject(name string) (Project, error) {
	proj := Project{Name: name}
	if err := g.DB.Create(&proj).Error; err != nil {
		return proj, utils.ErrCannotCreateProject
	}
	return proj, nil
}

// TODO: check for cascade delete functionality for GORM
func (g *GormProjectRepository) DeleteProject(projectID uint, er EntryRepository) error {
	er.DeleteEntries(projectID)
	if err := g.DB.Delete(&Project{}, projectID).Error; err != nil {
		return utils.ErrCannotDeleteProject
	}
	return nil
}

// TODO: make pe's Project a *Project instead to simplify?
func (g *GormProjectRepository) RenameProject(pe *ProjectWithEntries) {
	name := NewProjectPrompt()
	var project Project
	if err := g.DB.Where("id = ?", pe.Project.ID).First(&project).Error; err != nil {
		log.Fatalf("Unable to rename project: %q", err)
	}
	project.Name = name
	pe.Project.Name = name
	if err := g.DB.Save(&project).Error; err != nil {
		log.Fatalf("Unable to save project: %q", err)
	}
}

// TODO: move helper?
// helpers

func NewProjectPrompt() string {
	var name string
	fmt.Println("what would you like to name your project?")
	fmt.Scanf("%s", &name)
	return name
}

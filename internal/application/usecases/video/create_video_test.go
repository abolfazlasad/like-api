package videousecase_test

import (
	"testing"

	videousecase "like-api/internal/application/usecases/video"
	"like-api/internal/domain/entities"
	"like-api/internal/infrastructure/database/memory"
)

func TestCreateVideo_Success(t *testing.T) {
	userRepo, _, videoRepo, _ := memory.InitRepositories()
	userRepo.Create(entities.User{ID: "u1"})

	uc := videousecase.NewCreateVideoUseCase(userRepo, videoRepo)
	out := uc.Execute(videousecase.CreateVideoInput{
		UserID:      "u1",
		Title:       "My Video",
		Description: "A description",
		VideoURL:    "https://example.com/v.mp4",
	})

	if out.Error != nil {
		t.Fatalf("unexpected error: %v", out.Error)
	}
	if out.Video.ID == "" {
		t.Error("expected non-empty video ID")
	}
	if out.Video.Title != "My Video" {
		t.Errorf("unexpected title: %s", out.Video.Title)
	}
	if out.Video.UserID != "u1" {
		t.Errorf("unexpected userID: %s", out.Video.UserID)
	}
	if out.Video.LikesCount != 0 || out.Video.ViewsCount != 0 {
		t.Error("new video should have zero counts")
	}
	if out.Video.CreatedAt == "" {
		t.Error("CreatedAt should be set")
	}
}

func TestCreateVideo_IsPersisted(t *testing.T) {
	userRepo, _, videoRepo, _ := memory.InitRepositories()
	userRepo.Create(entities.User{ID: "u1"})

	uc := videousecase.NewCreateVideoUseCase(userRepo, videoRepo)
	out := uc.Execute(videousecase.CreateVideoInput{
		UserID:   "u1",
		Title:    "Persisted",
		VideoURL: "https://example.com/v.mp4",
	})

	if out.Error != nil {
		t.Fatalf("unexpected error: %v", out.Error)
	}

	_, err := videoRepo.FindByID(out.Video.ID)
	if err != nil {
		t.Errorf("video not found in repository after create: %v", err)
	}
}

func TestCreateVideo_UniqueIDs(t *testing.T) {
	userRepo, _, videoRepo, _ := memory.InitRepositories()
	userRepo.Create(entities.User{ID: "u1"})

	uc := videousecase.NewCreateVideoUseCase(userRepo, videoRepo)

	out1 := uc.Execute(videousecase.CreateVideoInput{UserID: "u1", Title: "V1", VideoURL: "https://a.com/1.mp4"})
	out2 := uc.Execute(videousecase.CreateVideoInput{UserID: "u1", Title: "V2", VideoURL: "https://a.com/2.mp4"})

	if out1.Video.ID == out2.Video.ID {
		t.Error("two videos should not share the same ID")
	}
}

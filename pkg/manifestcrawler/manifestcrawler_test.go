package manifestcrawler

import (
	"testing"
)

func TestGetImagesByRepo(t *testing.T) {
	// GetImagesByRepo(folder string, repoUrl string, extensions []string) ([]string, error)

	// got images
	extensions := []string{}
	images_empty_ext, _, _, _ := GetImagesByRepo("../test/repo", "quay.io", extensions)
	got_image_count_empty_ext := len(images_empty_ext)
	want_image_count_empty_ext := 4
	if got_image_count_empty_ext != want_image_count_empty_ext {
		t.Errorf("GetImagesByRepo(): ext[]. Wrong amout of images returned. got %d != want %d", got_image_count_empty_ext, want_image_count_empty_ext)
	}

	// got images
	extensions = []string{".yaml", ".yml"}
	images, _, _, _ := GetImagesByRepo("../test/repo", "quay.io", extensions)
	got_image_count := len(images)
	want_image_count := 4
	if got_image_count != want_image_count {
		t.Errorf("GetImagesByRepo(): ext[yml, yaml]. Wrong amout of images returned. got %d != want %d", got_image_count, want_image_count)
	}

	// got images
	extensions = []string{".yml"}
	images_yml, _, _, _ := GetImagesByRepo("../test/repo", "quay.io", extensions)
	got_image_count_yml := len(images_yml)
	want_image_count_yml := 0
	if got_image_count_yml != want_image_count_yml {
		t.Errorf("GetImagesByRepo(): ext[yml]. Wrong amout of images returned. got %d != want %d", got_image_count_yml, want_image_count_yml)
	}
}

func TestGetImages(t *testing.T) {

}

func TestGetImage(t *testing.T) {

}

func TestGetImageProjects(t *testing.T) {

}

func TestBuildUniModels(t *testing.T) {

}

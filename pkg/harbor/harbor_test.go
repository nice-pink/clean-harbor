package harbor

import (
	"testing"

	"github.com/nice-pink/clean-harbor/pkg/models"
	"github.com/nice-pink/clean-harbor/pkg/test/mock"
	// "github.com/nice-pink/clean-harbor/pkg/test/payload"
)

// all

func TestGetAll(t *testing.T) {
	config := HarborConfig{
		DryRun:    true,
		HarborUrl: "api.url",
		BasicAuth: models.Auth{BasicUser: "user", BasicPassword: "password"},
	}

	body := "[]"

	r := &mock.MockRequester{JsonBody: body}
	h := NewHarbor(r, config)
	projects := h.GetAll()
	if len(projects) == 0 {
		t.Error("nothing returned")
	}
}

// helper

func (h *Harbor) TestGetQuery(t *testing.T) {
	// empty query
	got_query_empty := h.GetQuery(0, 0)
	want_query_empty := ""
	if got_query_empty != want_query_empty {
		t.Errorf("got %q != want %q", got_query_empty, want_query_empty)
	}

	// page query
	got_query_page := h.GetQuery(2, 0)
	want_query_page := "?page=2"
	if got_query_page != want_query_page {
		t.Errorf("got %q != want %q", got_query_page, want_query_page)
	}

	// page_size query
	got_query_page_size := h.GetQuery(0, 2)
	want_query_page_size := "?page_size=2"
	if got_query_page_size != want_query_page_size {
		t.Errorf("got %q != want %q", got_query_page_size, want_query_page_size)
	}

	// page and size query
	got_query_page_and_size := h.GetQuery(1, 2)
	want_query_page_and_size := "?page=1&page_size=2"
	if got_query_page_and_size != want_query_page_and_size {
		t.Errorf("got %q != want %q", got_query_page_and_size, want_query_page_and_size)
	}
}

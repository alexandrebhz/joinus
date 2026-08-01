package utils

import "testing"

func TestExcerpt(t *testing.T) {
	if got := Excerpt("hello", 10); got != "hello" {
		t.Fatalf("short string: got %q", got)
	}
	got := Excerpt("abcdefghijklmnopqrstuvwxyz", 10)
	if got != "abcdefghi…" {
		t.Fatalf("truncate: got %q", got)
	}
}

func TestClampPagination(t *testing.T) {
	page, size := ClampPagination(0, 0, MaxPageSizePublic)
	if page != 1 || size != DefaultPageSize {
		t.Fatalf("defaults: page=%d size=%d", page, size)
	}
	page, size = ClampPagination(9999, 9999, MaxPageSizePublic)
	if page != MaxPage || size != MaxPageSizePublic {
		t.Fatalf("caps: page=%d size=%d", page, size)
	}
	_, size = ClampPagination(1, 100, MaxPageSizeAuth)
	if size != 100 {
		t.Fatalf("auth size: got %d", size)
	}
}

func TestSanitizeOrder(t *testing.T) {
	col, dir := SanitizeOrder("title; drop table", "asc", JobOrderColumns, "created_at")
	if col != "created_at" || dir != "ASC" {
		t.Fatalf("reject inject: col=%s dir=%s", col, dir)
	}
	col, dir = SanitizeOrder("title", "ASC", JobOrderColumns, "created_at")
	if col != "title" || dir != "ASC" {
		t.Fatalf("allow: col=%s dir=%s", col, dir)
	}
	if !IsHTTPURL("https://example.com/apply") {
		t.Fatal("https should pass")
	}
	if IsHTTPURL("javascript:alert(1)") {
		t.Fatal("javascript: must fail")
	}
}

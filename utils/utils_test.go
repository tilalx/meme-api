package utils

import (
	"testing"

	"Meme_Api/models"
)

// helpers

func makeMemes(urls ...string) []models.Meme {
	memes := make([]models.Meme, len(urls))
	for i, u := range urls {
		memes[i] = models.Meme{URL: u, Title: u}
	}
	return memes
}

func makeMemesWithFlags(count int, nsfw, spoiler bool) []models.Meme {
	memes := make([]models.Meme, count)
	for i := range memes {
		memes[i] = models.Meme{URL: "https://i.redd.it/x.jpg", NSFW: nsfw, Spoiler: spoiler}
	}
	return memes
}

// isImageURL

func TestIsImageURL(t *testing.T) {
	cases := []struct {
		url  string
		want bool
	}{
		{"https://i.redd.it/abc123.jpg", true},
		{"https://i.redd.it/abc123.png", true},
		{"https://example.com/image.jpg", true},
		{"https://example.com/image.jpeg", true},
		{"https://example.com/image.png", true},
		{"https://example.com/image.gif", true},
		{"https://example.com/image.webp", true},
		{"https://example.com/image.JPG", true},  // case-insensitive
		{"https://example.com/image.jpg?v=1", true}, // query string stripped
		{"https://www.reddit.com/r/memes", false},
		{"https://example.com/video.mp4", false},
		{"https://example.com/", false},
		{"", false},
	}

	for _, tc := range cases {
		got := isImageURL(tc.url)
		if got != tc.want {
			t.Errorf("isImageURL(%q) = %v, want %v", tc.url, got, tc.want)
		}
	}
}

// RemoveNonImagePosts

func TestRemoveNonImagePosts(t *testing.T) {
	input := makeMemes(
		"https://i.redd.it/abc.jpg",
		"https://www.reddit.com/r/memes/comments/abc",
		"https://example.com/img.png",
		"https://example.com/video.mp4",
		"https://example.com/img.webp",
	)

	got := RemoveNonImagePosts(input)

	if len(got) != 3 {
		t.Fatalf("expected 3 image posts, got %d", len(got))
	}
}

func TestRemoveNonImagePostsAllFiltered(t *testing.T) {
	input := makeMemes("https://reddit.com/link", "https://v.redd.it/video")
	got := RemoveNonImagePosts(input)
	if len(got) != 0 {
		t.Errorf("expected 0 posts, got %d", len(got))
	}
}

// PickMemeAt

func TestPickMemeAt(t *testing.T) {
	memes := makeMemes("a", "b", "c")

	if PickMemeAt(memes, 0).URL != "a" {
		t.Error("index 0 should return first meme")
	}
	if PickMemeAt(memes, 2).URL != "c" {
		t.Error("index 2 should return third meme")
	}
	// wraparound
	if PickMemeAt(memes, 3).URL != "a" {
		t.Error("index 3 should wrap to first meme")
	}
	if PickMemeAt(memes, 4).URL != "b" {
		t.Error("index 4 (4%3=1) should return second meme")
	}
}

// PickNMemesFrom

func TestPickNMemesFrom(t *testing.T) {
	memes := makeMemes("a", "b", "c", "d", "e")

	got := PickNMemesFrom(memes, 0, 3)
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
	if got[0].URL != "a" || got[1].URL != "b" || got[2].URL != "c" {
		t.Errorf("unexpected memes: %v", got)
	}
}

func TestPickNMemesFromWraparound(t *testing.T) {
	memes := makeMemes("a", "b", "c")

	// start=2, n=3 → c, a, b
	got := PickNMemesFrom(memes, 2, 3)
	if got[0].URL != "c" || got[1].URL != "a" || got[2].URL != "b" {
		t.Errorf("unexpected wraparound result: %v", got)
	}
}

// GetRandomN

func TestGetRandomNInRange(t *testing.T) {
	for i := 0; i < 1000; i++ {
		n := GetRandomN(10)
		if n < 0 || n >= 10 {
			t.Errorf("GetRandomN(10) = %d, out of range [0,10)", n)
		}
	}
}

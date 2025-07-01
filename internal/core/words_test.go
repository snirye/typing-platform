package core

import (
	"testing"
)

func TestNewWordManager(t *testing.T) {
	wm := NewWordManager()

	if wm == nil {
		t.Fatal("NewWordManager() returned nil")
	}

	if len(wm.Words) == 0 {
		t.Error("WordManager should have words loaded")
	}

	if wm.Difficulty != 1 {
		t.Errorf("Expected default difficulty 1, got %d", wm.Difficulty)
	}
}

func TestGetRandomWord(t *testing.T) {
	wm := NewWordManager()

	word := wm.GetRandomWord()
	if word == "" {
		t.Error("GetRandomWord() returned empty string")
	}

	// Test that we get different words (run multiple times)
	words := make(map[string]bool)
	for i := 0; i < 10; i++ {
		w := wm.GetRandomWord()
		words[w] = true
	}

	// Should have at least 2 different words in 10 attempts
	if len(words) < 2 {
		t.Error("GetRandomWord() doesn't seem to be returning random words")
	}
}

func TestSetDifficulty(t *testing.T) {
	wm := NewWordManager()

	// Test valid difficulty levels
	for i := 1; i <= 3; i++ {
		wm.SetDifficulty(i)
		if wm.Difficulty != i {
			t.Errorf("Expected difficulty %d, got %d", i, wm.Difficulty)
		}
	}

	// Test invalid difficulty levels
	originalDifficulty := wm.Difficulty
	wm.SetDifficulty(0)
	if wm.Difficulty != originalDifficulty {
		t.Error("SetDifficulty should not accept 0")
	}

	wm.SetDifficulty(4)
	if wm.Difficulty != originalDifficulty {
		t.Error("SetDifficulty should not accept 4")
	}
}

func TestIsWordComplete(t *testing.T) {
	wm := NewWordManager()

	tests := []struct {
		word     string
		typed    string
		expected bool
	}{
		{"hello", "hello", true},
		{"hello", "Hell", false},
		{"hello", "hello!", false},
		{"HELLO", "hello", true}, // Case insensitive
		{"test", "test", true},
		{"test", "tes", false},
		{"", "", true},
	}

	for _, test := range tests {
		result := wm.IsWordComplete(test.word, test.typed)
		if result != test.expected {
			t.Errorf("IsWordComplete(%q, %q) = %v, expected %v",
				test.word, test.typed, result, test.expected)
		}
	}
}

func TestIsValidChar(t *testing.T) {
	wm := NewWordManager()

	tests := []struct {
		word     string
		typed    string
		char     rune
		expected bool
	}{
		{"hello", "", 'h', true},
		{"hello", "h", 'e', true},
		{"hello", "he", 'l', true},
		{"hello", "", 'x', false},
		{"hello", "h", 'x', false},
		{"HELLO", "", 'h', true},       // Case insensitive
		{"HELLO", "", 'H', true},       // Case insensitive
		{"hello", "hello", 'x', false}, // Already complete
	}

	for _, test := range tests {
		result := wm.IsValidChar(test.word, test.typed, test.char)
		if result != test.expected {
			t.Errorf("IsValidChar(%q, %q, %q) = %v, expected %v",
				test.word, test.typed, string(test.char), result, test.expected)
		}
	}
}

func TestDifficultyFiltering(t *testing.T) {
	wm := NewWordManager()

	// Override words with known lengths for testing
	wm.Words = []string{
		"go",          // 2 chars
		"cat",         // 3 chars
		"word",        // 4 chars
		"hello",       // 5 chars
		"longer",      // 6 chars
		"testing",     // 7 chars
		"programming", // 11 chars
	}

	// Test difficulty 1 (3-5 chars)
	wm.SetDifficulty(1)
	shortWords := make(map[string]bool)
	for i := 0; i < 20; i++ {
		word := wm.GetRandomWord()
		shortWords[word] = true
		if len(word) < 3 || len(word) > 5 {
			t.Errorf("Difficulty 1 returned word '%s' with length %d", word, len(word))
		}
	}

	// Test difficulty 3 (6+ chars)
	wm.SetDifficulty(3)
	longWords := make(map[string]bool)
	for i := 0; i < 20; i++ {
		word := wm.GetRandomWord()
		longWords[word] = true
		if len(word) < 6 {
			t.Errorf("Difficulty 3 returned word '%s' with length %d", word, len(word))
		}
	}
}

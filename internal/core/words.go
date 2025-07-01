package core

import (
	"math/rand"
	"strings"
	"time"
)

// WordManager handles word selection and management
type WordManager struct {
	Words      []string
	UsedWords  map[string]bool
	Difficulty int
	rng        *rand.Rand
}

// NewWordManager creates a new word manager
func NewWordManager() *WordManager {
	// Default word list - in a real implementation, this would load from assets/words.txt
	defaultWords := []string{
		"the", "and", "for", "are", "but", "not", "you", "all", "can", "her", "was", "one",
		"our", "had", "day", "get", "use", "man", "new", "now", "way", "may", "say", "each",
		"which", "their", "time", "will", "about", "would", "there", "could", "other", "after",
		"first", "never", "these", "think", "where", "being", "every", "great", "might", "shall",
		"still", "those", "while", "write", "place", "right", "where", "sound", "again", "below",
		"between", "important", "children", "example", "sentence", "following", "without", "another",
		"different", "thought", "through", "before", "picture", "country", "together", "followed",
		"programming", "computer", "keyboard", "function", "variable", "algorithm", "structure",
		"interface", "development", "framework", "library", "package", "compile", "execute",
	}

	return &WordManager{
		Words:      defaultWords,
		UsedWords:  make(map[string]bool),
		Difficulty: 1,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GetRandomWord returns a random word based on difficulty
func (wm *WordManager) GetRandomWord() string {
	var availableWords []string

	// Filter words based on difficulty
	for _, word := range wm.Words {
		wordLen := len(word)
		switch wm.Difficulty {
		case 1: // Easy: 3-5 characters
			if wordLen >= 3 && wordLen <= 5 {
				availableWords = append(availableWords, word)
			}
		case 2: // Medium: 4-8 characters
			if wordLen >= 4 && wordLen <= 8 {
				availableWords = append(availableWords, word)
			}
		case 3: // Hard: 6+ characters
			if wordLen >= 6 {
				availableWords = append(availableWords, word)
			}
		default: // All words
			availableWords = append(availableWords, word)
		}
	}

	if len(availableWords) == 0 {
		availableWords = wm.Words // Fallback to all words
	}

	// Select random word
	word := availableWords[wm.rng.Intn(len(availableWords))]
	return strings.ToLower(word)
}

// SetDifficulty sets the difficulty level (1-3)
func (wm *WordManager) SetDifficulty(level int) {
	if level >= 1 && level <= 3 {
		wm.Difficulty = level
	}
}

// IsWordComplete checks if a word is completely typed
func (wm *WordManager) IsWordComplete(word, typed string) bool {
	return strings.ToLower(word) == strings.ToLower(typed)
}

// IsValidChar checks if the next character in typing is valid
func (wm *WordManager) IsValidChar(word, typed string, char rune) bool {
	if len(typed) >= len(word) {
		return false
	}

	expectedChar := rune(strings.ToLower(word)[len(typed)])
	return strings.ToLower(string(char)) == string(expectedChar)
}

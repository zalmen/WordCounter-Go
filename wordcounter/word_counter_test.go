package wordcounter

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWordCounter_wordValidity(t *testing.T) {
	wc := New()

	wc.CountWord("")
	assert.Zero(t, wc.LowestFrequency())

	wc.CountWord("   ")
	assert.Zero(t, wc.LowestFrequency())

	wc.CountWord("cyolo")
	wc.CountWord("cyolo   ")
	wc.CountWord("   cyolo")
	wc.CountWord("   cyolo  ")
	assert.Equal(t, 4, wc.LowestFrequency())
	assert.Equal(t, map[string]int{"cyolo": 4}, wc.TopFiveFrequentWords())
}

func TestWordCounter_singleWord(t *testing.T) {
	wc := New()
	word := "hello"
	freq := 20
	for i := 0; i < freq; i++ {
		wc.CountWord(word)
	}

	assert.Equal(t, freq, wc.LowestFrequency())
	assert.Equal(t, freq, wc.MedianFrequency())
	assert.Equal(t, map[string]int{word: 20}, wc.TopFiveFrequentWords())
}

func isTop5Distinct(top5 map[string]int) bool {
	saw := make(map[string]bool)
	for w, _ := range top5 {
		if saw[w] {
			return false
		}
		saw[w] = true
	}
	return true
}

func TestWordCounter_multipleWords(t *testing.T) {
	wc := New()

	words := [7]string{}
	for i := 0; i < 7; i++ {
		word, _ := uuid.NewUUID()
		words[i] = word.String()
	}
	// frequencies (idx to frequency):
	// 0 -> 1, 1 -> 2, 2 -> 1, 3 -> 2, 4 -> 2, 5 -> 2, 6 -> 2

	wc.CountWord(words[1])
	wc.CountWord(words[3])
	wc.CountWord(words[6])
	wc.CountWord(words[4])
	wc.CountWord(words[4])
	wc.CountWord(words[5])
	wc.CountWord(words[0])
	wc.CountWord(words[3])
	wc.CountWord(words[6])
	wc.CountWord(words[2])
	wc.CountWord(words[5])
	wc.CountWord(words[1])

	assert.Equal(t, 1, wc.LowestFrequency())
	assert.Equal(t, 2, wc.MedianFrequency())
	top5 := wc.TopFiveFrequentWords()
	assert.Equal(t, 5, len(top5))
	assert.True(t, isTop5Distinct(top5))
}

func TestWordCounter_multipleDistinctWords(t *testing.T) {
	wc := New()
	count := 200
	for i := 0; i < count; i++ {
		word, _ := uuid.NewUUID()
		wc.CountWord(word.String())
	}

	assert.Equal(t, 1, wc.LowestFrequency())
	assert.Equal(t, 1, wc.MedianFrequency())
	top5 := wc.TopFiveFrequentWords()
	assert.Equal(t, 5, len(top5))
	assert.True(t, isTop5Distinct(top5))
}

func TestWordCounter_multipleWordsDistinctFrequencies(t *testing.T) {
	wc := New()

	words := [6]string{}
	for i := 0; i < 6; i++ {
		word, _ := uuid.NewUUID()
		words[i] = word.String()
	}
	// frequencies (idx to frequency):
	// 0 -> 7, 1 -> 4, 2 -> 1, 3 -> 2, 4 -> 6, 5 -> 3

	wc.CountWord(words[0])
	wc.CountWord(words[0])
	wc.CountWord(words[1])
	wc.CountWord(words[3])
	wc.CountWord(words[2])
	wc.CountWord(words[0])
	wc.CountWord(words[4])
	wc.CountWord(words[4])
	wc.CountWord(words[5])
	wc.CountWord(words[0])
	wc.CountWord(words[3])
	wc.CountWord(words[4])
	wc.CountWord(words[0])
	wc.CountWord(words[0])
	wc.CountWord(words[4])
	wc.CountWord(words[1])
	wc.CountWord(words[5])
	wc.CountWord(words[4])
	wc.CountWord(words[1])
	wc.CountWord(words[0])
	wc.CountWord(words[4])
	wc.CountWord(words[5])
	wc.CountWord(words[1])

	assert.Equal(t, 1, wc.LowestFrequency())
	assert.Equal(t, 4, wc.MedianFrequency())
	top5 := wc.TopFiveFrequentWords()
	assert.Equal(t, 5, len(top5))
	assert.True(t, isTop5Distinct(top5))
}

func TestWordCounter_defaultStats(t *testing.T) {
	wc := New()
	assert.Zero(t, wc.MedianFrequency())
	assert.Zero(t, wc.LowestFrequency())
	assert.Zero(t, len(wc.TopFiveFrequentWords()))
}

func TestWordCounter_exampleTestcase(t *testing.T) {
	wc := New()
	wc.CountWord("ball")
	wc.CountWord("eggs")
	wc.CountWord("pool")
	wc.CountWord("dart")
	wc.CountWord("ball")
	wc.CountWord("ball")
	wc.CountWord("table")
	wc.CountWord("eggs")
	wc.CountWord("pool")
	wc.CountWord("mouse")
	wc.CountWord("ball")
	wc.CountWord("eggs")
	wc.CountWord("table")
	wc.CountWord("mouse")

	assert.Equal(t, 1, wc.LowestFrequency())
	assert.Equal(t, 2, wc.MedianFrequency())
	assert.Equal(t, map[string]int{"ball": 4, "eggs": 3, "mouse": 2, "pool": 2, "table": 2}, wc.TopFiveFrequentWords())
}

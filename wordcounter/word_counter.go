package wordcounter

import (
	"container/list"
	"strings"
)

type wordCounter struct {
	// Map between words and the number of times they appeared
	wordsHistogram map[string]int
	// List of all the distinct frequencies, kept sorted
	sortedFrequencies *list.List
	// Map between frequencies and their element in `sortedFrequencies`, allowing const time access to the list elements
	markers map[int]*Marker
	// Top 5 recurring words
	top5 [5]string
}

type Marker struct {
	position  *list.Element // the position of the frequency in `sortedFrequencies`
	batchSize int           // number of words that appeared this number of times, needed since each distinct frequency
	// has a single element in `sortedFrequencies`
}

/*
Constructor
*/
func New() *wordCounter {
	return &wordCounter{
		wordsHistogram:    make(map[string]int),
		sortedFrequencies: list.New(),
		markers:           make(map[int]*Marker),
	}
}

/*
updateTop5 is used for maintaining the `top5` state whenever `CountWord` is called
*/
func (wc *wordCounter) updateTop5(word string, frequency int) {
	curWordIdx := 5
	minWordIdx := 5
	minWordFreq := 0
	for i := 0; i < 5; i++ {
		if wc.top5[i] == word {
			// `word` is already in top5
			curWordIdx = i
		}
		iWordFreq := wc.wordsHistogram[wc.top5[i]]
		if minWordIdx == 5 || iWordFreq < minWordFreq {
			minWordIdx = i
			minWordFreq = iWordFreq
		}
	}
	if curWordIdx == 5 && frequency > minWordFreq {
		// we need to update `wc.top5` only if `word` isn't already part of top5 and if it appears more times
		// than the word that appear the least in top5 (`minWordFreq`)
		wc.top5[minWordIdx] = word
	}
}

/*
addFrequency updates `sortedFrequencies` and `markers` on a new distinct frequency
*/
func (wc *wordCounter) addFrequency(frequency int) {
	var elem *list.Element
	if frequency == 1 {
		elem = wc.sortedFrequencies.PushFront(frequency)
	} else {
		prevFreq := frequency - 1
		prevMarker := wc.markers[prevFreq] // must exist since word's frequency was `prevFreq` until now
		elem = wc.sortedFrequencies.InsertAfter(frequency, prevMarker.position)
	}
	wc.markers[frequency] = &Marker{
		position:  elem,
		batchSize: 1,
	}
}

/*
updateMarkers is used for maintaining `sortedFrequencies` and `markers` whenever `CountWord` is called. If there are
currently other words that appeared `curFreq` times it will only maintain the relevant markers' `batchSize`, otherwise
it will add `curFreq` to `sortedFrequencies` and `markers`
*/
func (wc *wordCounter) updateMarkers(curFreq int) {
	// update the current frequency marker
	if m, ok := wc.markers[curFreq]; !ok {
		// no marker for `curFreq`, adding one
		wc.addFrequency(curFreq)
	} else {
		// there is already a marker for `curFreq`, incrementing its `batchSize`
		m.batchSize++
	}
	// update the previous frequency marker
	if curFreq > 1 {
		// if it is a known word we need to decrement its previous (`curFreq` - 1) marker's `batchSize`
		prevFreq := curFreq - 1
		prevMarker := wc.markers[prevFreq] // must exist since `word`'s frequency was `prevFreq` until now
		if prevMarker.batchSize == 1 {
			// the previous marker was only pointing on this word's frequency
			wc.sortedFrequencies.Remove(prevMarker.position)
			delete(wc.markers, prevFreq)
		} else {
			prevMarker.batchSize--
		}
	}
}

/*
updateState updates wc's state when `CountWord` is called
*/
func (wc *wordCounter) updateState(word string, newFreq int) {
	wc.wordsHistogram[word] = newFreq
	wc.updateTop5(word, newFreq)
	wc.updateMarkers(newFreq) //must be called after updating `wordsHistogram`
}

/*
CountWord is an exported function used for reporting words
*/
func (wc *wordCounter) CountWord(word string) {
	trimmed := strings.TrimSpace(word)
	if len(trimmed) == 0 {
		return
	}
	newFreq := wc.wordsHistogram[trimmed] + 1
	wc.updateState(trimmed, newFreq)
}

/*
TopFiveFrequentWords returns the top 5 recurring words with their frequency
*/
func (wc *wordCounter) TopFiveFrequentWords() map[string]int {
	res := make(map[string]int)
	for _, w := range wc.top5 {
		if h := wc.wordsHistogram[w]; h != 0 {
			res[w] = h
		}
	}
	return res
}

/*
LowestFrequency returns the frequency of the words the appeared the least so far
*/
func (wc *wordCounter) LowestFrequency() int {
	if wc.sortedFrequencies.Len() == 0 {
		return 0
	}
	return wc.sortedFrequencies.Front().Value.(int)
}

/*
MedianFrequency returns the median of the frequencies of all the words that appeared so far
*/
func (wc *wordCounter) MedianFrequency() int {
	if wc.sortedFrequencies.Len() == 0 {
		return 0
	}
	totalWordsCount := len(wc.wordsHistogram)
	var medianIndex int = totalWordsCount / 2

	marker := wc.markers[wc.LowestFrequency()]
	idx := marker.batchSize - 1
	iter := marker.position

	for idx < medianIndex {
		nextFreq := iter.Next().Value.(int)
		marker = wc.markers[nextFreq]
		idx += marker.batchSize
		iter = marker.position
	}
	return iter.Value.(int)
}

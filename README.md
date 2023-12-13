The code includes three packages: 
- httpserver which implements a REST api with two endpoints - one for providing the words ("words") and one for fetching the statistics ("stats). 
- wordcounter which implements the core functionality needed to serve the two endpoints.
- main

When designing the solution for this assignment I had to balance between the time it takes to handle a word submission, the time it takes to retrieve the statistics and code complexity. The assignment mentioned that words are provided more frequently than requests to retrieve the statistics, hence I only considered solutions that had low time complexity when processing a word on the expense of the time complexity of retrieving the statistics (e.g. I ruled out a solution that maintains a median-heap since it takes O(log(n)) to insert values). 

I considered three options:
1. Maintain a histogram for all the words, update it on every word submitted and sort all the word frequencies when handling a request to retrieve the statistics.

Word submission - O(1)  
Statistics retrieval - O(nlog(n)) where n is the total number of distinct frequencies
Code complexity - low
Memory consumption - low

2. Maintaining a histogram for all the words which is updated on every word submitted and extracting the median w/o sorting the entire list of frequencies (utilizing median of medians or quick select algorithms), when handling a request to retrieve the statistics.

Word submission - O(1)  
Statistics retrieval - O(n) where n is the total number of distinct frequencies
Code complexity - high
Memory consumption - low

3. Maintaining a histogram for all the words and a sorted list of all distinct frequencies with O(1) access, when handling a word submission. When handling a request to retrieve the statistics we only need to extract the median from the sorted list.

Word submission - O(1). It is not trivial that we can maintain the list of frequency sorted in O(1), it can be achieved due to the fact that there are only two cases where we modify this list - when a new word is submitted and we need to insert the frequency of 1 at the beginning of the list (O(1)) or when an existing word is submitted and we need to update its frequency from `a` to `a+1` where it is guaranteed that `a` is already in the list and can be accessed in O(1) time.

Statistics retrieval - O(n) where n is the number of distinct frequencies.

Code complexity - medium
Memory consumption - medium

For the sake of this assignment I assumed that the use-case is time critical and therefore I chose to implement option 3 since it has the best performance with a reasonable code complexity. If the use-case wasn't time critical I might have chosen a simpler solution which prioritizes code complexity over performance.

Possible improvements (weren't implemented):
1. The histogram of words is implemented as a `map` which stores the actual words (strings). Since I don't need to iterate over the words there is no reason to store the actual strings and the histogram can be keyed by a hash of each word, reducing the memory consumption significantly.
2. It is possible to add to the state a median member which will be maintained when handling a word submission (in a nutshell, in a list of size one, the median points at the single element. Whenever a new element is inserted the median can remain untouched or shifted one position to the left or to the right depending on the list length parity and whether the new value was larger or smaller than the current median). Adding this will allow us to retrieve the statistics in O(1) time while keeping the word submission at O(1) as well. I didn't implement it in order to avoid making the code more complex.

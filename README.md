# Log-based key-value store

The goal is to implement an embedded log-based key-value store as described in 
Designing Data-Intesive Applications - Chapter 3 - Storage and Retrieval. This is not
the SSTable implementation but rather the simple log-based one described at the very
beginning

This is not a huge undertaking but to guarantee it's completion I'll be breaking down
the task into several distinct stages - each one with a clear goal. That way I'll be
working my way up to a fully functional implementation with (decent throughput).

## Stages

1. [x] Key-value store using a map (basic API)
2. [x] Store data by writing and reading to/from a single file (a Log)
3. [ ] Split the file into multiple smaller files when a certain size is reached (Segments)
4. [ ] Use an index to store the position of each key in the appropriate segment
5. [ ] Add a background process that compacts the segments after some time
6. [ ] Allow the segment size and time to be configured


## Things that can be improved

- Using a binary format to store the data, that means it is no longer human-readable
but has better performance (potentially less storage used and faster reads as we know
beforehand the length of each key or value)

## Learning experience

### Mocking
As it's compiled language, there's no real "mocking" as you have in other languages like
Python or Javascript. This is don using DI (Dependency Injection) and injecting a "mock"
for the interface.

There are several sophisticated "mock" libraries but as I didn't really need anything
fancy, I just rolled a very basic one on my own.

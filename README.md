# Log-based key-value store

The goal is to implement an embedded log-based key-value store as described in 
Designing Data-Intesive Applications - Chapter 3 - Storage and Retrieval. This is not
the SSTable implementation but rather the simple log-based one described at the very
beginning

This is not a huge undertaking but to guarantee it's completion I'll be breaking down
the task into several distinct stages - each one with a clear goal. That way I'll be
working my way up to a fully functional implementation with decent throughput.

Although not initially the idea, having a TCP server and a client to load/retrieve data
from the database has proven to be particularly useful (and fun).

## Stages

- [x] Key-value store using a map (basic API)
- [x] Store data by writing and reading to/from a single file (a Log)
- [x] Turn the Key-value store into a TCP server that accepts connections using a 
(very) simplified version of the REDIS protocol
- [x] Split the file into multiple smaller files when a certain size is reached (Segments)
- [x] Allow the segment size and time to be configured
- [ ] Use an index to store the position of each key in the appropriate segment
- [ ] Add a background process that compacts the segments after some time


## Things that can be improved

- Using a binary format to store the data, that means it is no longer human-readable
but has better performance (potentially less storage used and faster reads as we know
beforehand the length of each key or value)
- Improving the way the server handles connections. Right now it will happily
create any number of goroutines to respond to the load, potentially overloading a
system
- Improving the client to allow us to issue arbitrary commands instead of 
predefined ones


## Learning experience

### Mocking
As it's compiled language, there's no real "mocking" as you have in other languages like
Python or Javascript. This is achieved using DI (Dependency Injection) and injecting 
a "mock" for the interface.

There are several sophisticated "mock" libraries but as I didn't really need anything
fancy, I just rolled a very basic one of my own.

### Testing
Pretty straightforward. You test using `go test`, in this case it's `go test ./store`.
To have a coverage report:
```bash
# This allows you to get a percentage of coverage along the tests that passed/failed
# and stores the coverage report in "coverage.txt"
$ go test ./store -v -coverprofile=coverage.txt

# This has a very simple but useful UI to check what lines are covered in each file
# (opens an html file in the browser)
$ go tool cover -html=coverage.txt 
```

### CI/CD
Once again pretty straightforward using Github Actions. There's a default action to
trigger the test suite and it required minimal changes to work properly
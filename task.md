The task is to write an open-source Go library package that provides APIs to run external processes and put it on GitHub.
It should:
* start processes with given arguments and environment variables;
* stop them;
* restart them when they crash;
* relay termination signals;
* read their stdout and stderr;
* compile and work on Linux and macOS.

Optionally, you can provide (in random order):
* ability to stop processes when main processes are SIGKILL'ed;
* comments and documentation in code;
* configurable backoff strategy for restarts;
* README file;
* continuous integration configuration;
* integration tests;
* command (package main) that demonstrates the usage;
* unit tests.
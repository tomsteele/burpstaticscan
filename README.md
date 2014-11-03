burpstaticscan
==============

Use burp's JS static code analysis on code from your local system. Here's generally how the process works:
* Go static file server is started to host the specified directory
* Add file server URL to burp's scope
* Walk the directory
* For every file make a request to the file server
* Send the raw request and response to burp's passive scanner through burpbuddy
* Issues can be seen in burp

```
$ ./burpstaticscan -dir ./foo
2014/11/02 18:11:22 Static file server listening on http://localhost:9999, serving /foo
2014/11/02 18:11:22 Adding http://localhost:9999 to scope
2014/11/02 18:11:22 Walking directory, each file sent to burp's passive scan
2014/11/02 18:11:22 1 file sent to burp
2014/11/02 18:11:22 http://localhost:9999 removed from scope
```

## Installation
Depends on the [burpbuddy](https://github.com/liftsecurity/burpbuddy) extension. Binary packages for most operating systems are available [here](https://github.com/tomsteele/burpstaticscan/releases/latest). There are no external dependencies, just extract and run with `./burpstaticscan`.

## Usage
```
$ ./burpstaticscan -h
Usage of ./burpstaticscan:
  -burpbuddy="http://localhost:8001": HTTP API URL for burpbuddy
  -dir="": directory with code to scan
  -host="localhost": host for the file server to listen on
  -port="9999": port for the file server to listen on
```

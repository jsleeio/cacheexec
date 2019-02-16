# cacheexec

## tl;dr

Speeds up stuff in your shell initialization, like loading completion rules for
many commands, each of which requires executing that command...

```
source <(cacheexec -- kubectl completion bash)
source <(cacheexec -- helm completion bash)
...
source <(cacheexec -- eleventycommand completion bash)
```

## what is it?

`cacheexec`, per its name:

0. initializes a cache directory (by default, `$HOME/.cacheexec`)

1. checks to see if the specified command has been cached recently
 
2a. if cached, and the cached output is fresh, emits it to stdout

2b. if cached, but the output is stale, runs the command, saves output in the cache and emits to stdout

2c. if uncached, runs the command, saves the output in the cache and emits to stdout

## usage

```
Usage of ./cacheexec:
  -storepath string
    	location of cache store (default "/Users/jslee/.cacheexec")
  -ttl duration
    	maximum time before cache invalidation (default 8h0m0s)
```

Use a `--` to separate `cacheexec` options from the command to be executed.

## license

Copyright 2019 John Slee.  Released under the terms of the MIT license
[as included in this repository](LICENSE).

This project is no longer maintained.

Because I want to make a terminal multiplexer by Rust or Go,
but I decided to make it with C, technical reason.

But this project is good as exercise to handle ptys,
and in the web, it is few that information to handle ptys from Go,
so I leave this.

`go run sample.go`, and spawn a bash in other process.

TODO:

- This code supports only darwin, does not support linux, bsd, ...
- `lib/pty/pty_test.go` will fail when the test process does not connected stdio

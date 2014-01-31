yamatrix
========

Yet Another Matrix program

This is just for fun.

This one has low code volume, but not as tiny as http://bruxy.regnet.cz/web/linux/EN/matrix-sh/

I did not (yet) find a good way to get the terminal size,
so you have to start it with
```bash
lines_t=$(tput lines) cols_t=$(tput cols) matrix
```

If you run it with ˝go run matrix˝ it does not catch the parent process id,
so you need to start it like
```bash
ppid_t=$$ lines_t=$(tput lines) cols_t=$(tput cols) go run matrix.go
```

Stop it with Ctrl-C

Todo
----

* get the terminal size (see termbox)
* get it to run in a non-utf8 (i.e. latin1) terminal
* add a lot if funny utf-8 chars



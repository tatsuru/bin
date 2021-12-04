package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	var tmux = flag.String("tmux-command", "tmux", "tmux command name")
	var ssh = flag.String("ssh-command", "ssh", "ssh command")
	var windowTitle = flag.String("t", fmt.Sprintf("tssh-%d", os.Getpid()), "window title")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Error: server name[s] required.\n")
		flag.Usage()
		os.Exit(1)
	}

	servers := flag.Args()
	exec.Command(*tmux, "new-window", "-n", *windowTitle, fmt.Sprintf("%s %s", *ssh, servers[0])).Start()

	for i := 1; i < len(servers); i++ {
		exec.Command(*tmux, "split-window", "-t", *windowTitle, fmt.Sprintf("%s %s", *ssh, servers[i])).Start()
	}

	exec.Command(*tmux, "select-layout", "-t", *windowTitle, "tiled").Start()
	exec.Command(*tmux, "set-window", "-t", *windowTitle, "synchronize-panes").Start()

	fmt.Println(servers, *tmux, *ssh, *windowTitle)
}

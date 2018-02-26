package main

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"github.com/omakoto/hl2/src/hl/highlighter"
	"github.com/omakoto/hl2/src/hl/matcher"
	"github.com/omakoto/hl2/src/hl/term"
	"github.com/omakoto/hl2/src/hl/util"
	"github.com/pborman/getopt/v2"
	"io"
	"os"
	"os/exec"
	"runtime/pprof"
)

const (
	Name              = "hl"
	CommandTerminator = ","
	RangeSeparator    = "-"
)

var (
	ruleFile = getopt.StringLong("rule", 'r', "", "Specify TOML rule file.")

	after          = getopt.IntLong("after", 'A', 0, "Specify number of 'after' context lines.")
	before         = getopt.IntLong("before", 'B', 0, "Specify number of 'before' context lines.")
	context        = getopt.IntLong("context", 'C', 0, "Specify number of context lines.")
	ignoreCase     = getopt.BoolLong("ignore-case", 'i', "Perform case insensitive match.")
	defaultHide    = getopt.BoolLong("hide", 'n', "Hide all lines by default.")
	execute        = getopt.StringLong("command", 'c', CommandTerminator, "Execute a command and apply to output.\nOptionally specify command line terminator. (default="+CommandTerminator+")")
	eatStderr      = getopt.BoolLong("stderr", '2', "Use with -c; process stderr from command too.")
	width          = getopt.IntLong("width", 'w', term.GetTermWidth(), "Set terminal width, used for pre and post lines.")
	cpuprofile     = getopt.StringLong("cpuprofile", 'P', "", "Write cpu profile to file.")
	help           = getopt.BoolLong("help", 'h', "Show this help.")
	noTtyWarning   = getopt.BoolLong("no-tty-warning", 'q', "Don't show warning even when stdin is tty.")
	inFile         = getopt.StringLong("input", 'f', "", "Read input from specified file instead of stdins.")
	rangeSeparator = getopt.StringLong("range-separator", 's', RangeSeparator, "Specify range separator. (default="+RangeSeparator+")")

	executeOption = getopt.Lookup('c')
)

func init() {
	getopt.FlagLong(&util.Debug, "debug", 'd', "Enable debug output.")
	getopt.FlagLong(&matcher.NoPcre, "no-pcre", 'N', "Disable PCRE and use Go's regexp engine instead.")

	executeOption.SetOptional()
}

func getCommandTerminator() string {
	if !executeOption.Seen() {
		return ""
	}
	if len(*execute) > 0 {
		return *execute
	}
	return CommandTerminator
}

func main() {
	getopt.Parse()

	if *help {
		getopt.Usage()
		return
	}

	if *width > 0 {
		term.TermWidth = *width
	}
	if *context > 0 {
		*after = *context
		*before = *context
	}

	h := highlighter.NewHighlighter(term.NewTerm(), *ignoreCase, *defaultHide, *before, *after)
	util.Dump("Highlighter (start): ", h)

	if *ruleFile != "" {
		err := h.LoadToml(*ruleFile)
		if err != nil {
			Fatalf("Unable to read rule file: %s", err)
		}
	}
	err := parseArgs(h, getopt.Args(), executeOption.Seen(), getCommandTerminator(), *rangeSeparator)
	if err != nil {
		Fatalf("Invalid options: %s", err)
	}

	util.Dump("Highlighter (all built up): ", h)

	// Maybe start the profiler.
	if cleaner := mayStartProfiler(*cpuprofile); cleaner != nil {
		defer cleaner()
	}

	// Execute the command if one is passed.
	var in io.ReadCloser = os.Stdin

	if len(h.CommandLine()) > 0 {
		var cleaner func()
		in, cleaner = startCommand(h.CommandLine())
		defer cleaner()
	} else if *inFile != "" {
		in, err = os.Open(*inFile)
		if err != nil {
			Fatalf("Cannot open file %s: %s", *inFile, err)
		}
	}

	if !*noTtyWarning && in == os.Stdin && isatty.IsTerminal(os.Stdin.Fd()) {
		fmt.Fprint(os.Stderr, "Waiting for input from stdin. (Use -q to suppress this message.)\n")
	}

	err = h.NewRuntime(os.Stdout).ColorReader(in)
	if err != nil {
		Fatalf("Unknown failure: %s", err)
	}
}

func mayStartProfiler(outfile string) func() {
	if outfile == "" {
		return nil
	}
	f, err := os.Create(outfile)
	if err != nil {
		Fatalf(fmt.Sprintf("Unable to create %s: %s", outfile, err))
	}
	pprof.StartCPUProfile(f)

	return func() {
		pprof.StopCPUProfile()
	}
}

func startCommand(commandLine []string) (io.ReadCloser, func()) {
	cmd := exec.Command(commandLine[0], commandLine[1:]...)

	// Set up stdin and stdout.
	cmd.Stdin = os.Stdin

	in, err := cmd.StdoutPipe()
	if err != nil {
		Fatalf("Unable to obtain stdout pipe: %s", err)
	}
	if *eatStderr {
		cmd.Stderr = cmd.Stdout
	} else {
		cmd.Stderr = os.Stderr
	}

	// Then start it.
	err = cmd.Start()
	if err != nil {
		Fatalf("Unable to start command \"%v\": %s", commandLine, err)
	}
	return in, func() {
		cmd.Wait()
	}
}

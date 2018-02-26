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
	readFiles      = getopt.StringLong("files", 'f', "", "Similar to -c but specify input files.")
	rangeSeparator = getopt.StringLong("range-separator", 's', RangeSeparator, "Specify range separator. (default="+RangeSeparator+")")

	executeOption   = getopt.Lookup('c')
	readFilesOption = getopt.Lookup('f')
)

func init() {
	getopt.FlagLong(&util.Debug, "debug", 'd', "Enable debug output.")
	getopt.FlagLong(&matcher.NoPcre, "no-pcre", 'N', "Disable PCRE and use Go's regexp engine instead.")

	executeOption.SetOptional()
	readFilesOption.SetOptional()

	getopt.SetUsage(usage)
}

func usage() {
	os.Stderr.WriteString(`
hl2: Versatile coloring filter

Usage:
  hl2 -r RULE_TOML [OPTIONS]    (Read rules from RULE_TOML)
  hl2 [OPTIONS] COLOR-SPEC...   (Give color spec from command line)
  hl2 -c    [OPTIONS] COMMAND [ARGS...] [, FILTER-SPEC...]   (Apply to command output; -r can be used too)
  hl2 -cSEP [OPTIONS] COMMAND [ARGS...] [SEP FILTER-SPEC...] (Same as above but use arbitrary separator)
  hl2 -f    [OPTIONS] FILES... [, FILTER-SPEC...]   (Apply to FILES; -r can be used too)
  hl2 -fSEP [OPTIONS] FILES... [SEP FILTER-SPEC...] (Same as above but use arbitrary separator)

  FILTER-SPEC is a list of:
    PATTERN [ COLOR-SPEC ]
    PATTERN [ COLOR-SPEC ] '-' PATTERN [ COLOR-SPEC ] (Implies -n. '-' can be changed with -s)

  COLOR-SPEC is:
    '@' [ATTRS] [FG-COLOR] [/BG-COLOR] [ '@' [ATTRS] [LINE-FG-COLOR] [/LINE-BG-COLOR] ]

  ATTRS is a set of:
    b: Bold / intense
    i: Italic
    f: Faint
    u: Underline
    s: Strike-through

  COLOR is any of:
    black | red | green | yellow | blue | magenta | cyan | white
    [0-5][0-5][0-5]     (RGB: 216 colors)
    [0-9a-f]{6}         (RRGGBB: 24bit colors)

Options:
`)
	getopt.CommandLine.PrintOptions(os.Stderr)
	os.Stderr.WriteString("\n")
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

	// Initialize highlighter.
	h := highlighter.NewHighlighter(term.NewTerm(), *ignoreCase, *defaultHide, *before, *after)
	util.Dump("Highlighter (start): ", h)

	// Process -c and -f, and also extract simple (inline) rules.
	doExecute := executeOption.Seen()
	doReadFiles := readFilesOption.Seen()
	if doExecute && doReadFiles {
		Fatalf("Cannot use -c and -f at the same time.\n")
	}
	term := ""
	if doExecute {
		term = *execute
	} else if doReadFiles {
		term = *readFiles
	}
	if term == "" {
		term = CommandTerminator
	}

	inputArgs, err := parseArgs(h, getopt.Args(), doExecute || doReadFiles, term, *rangeSeparator)
	if err != nil {
		Fatalf("Invalid options: %s", err)
	}

	// Load TOML
	if *ruleFile != "" {
		err := h.LoadToml(*ruleFile)
		if err != nil {
			Fatalf("Unable to read rule file: %s", err)
		}
	}

	util.Dump("Highlighter (all built up): ", h)

	// Maybe start the profiler.
	if cleaner := mayStartProfiler(*cpuprofile); cleaner != nil {
		defer cleaner()
	}

	if doReadFiles {
		for _, f := range inputArgs {
			in, err := os.Open(f)
			if err != nil {
				Fatalf("Cannot open file %s: %s", f, err)
			}
			doOnReader(h, in)
		}
	} else {
		// Execute the command if one is passed.
		var in io.ReadCloser = os.Stdin

		if doExecute {
			var cleaner func()
			in, cleaner = startCommand(inputArgs)
			defer cleaner()
		}

		if !*noTtyWarning && in == os.Stdin && isatty.IsTerminal(os.Stdin.Fd()) {
			fmt.Fprint(os.Stderr, "Waiting for input from stdin. (Use -q to suppress this message.)\n")
		}
		doOnReader(h, in)
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

func doOnReader(h *highlighter.Highlighter, rd io.ReadCloser) {
	defer rd.Close()

	err := h.NewRuntime(os.Stdout).ColorReader(rd)
	if err != nil {
		Fatalf("Unknown failure: %s", err)
	}
}

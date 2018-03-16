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
	ArgumentSeparator = ","
)

var (
	ruleFile = getopt.StringLong("rule", 'r', "", "Specify TOML rule file.")

	after             = getopt.IntLong("after", 'A', 0, "Specify number of 'after' context lines.")
	before            = getopt.IntLong("before", 'B', 0, "Specify number of 'before' context lines.")
	context           = getopt.IntLong("context", 'C', 0, "Specify number of context lines.")
	ignoreCase        = getopt.BoolLong("ignore-case", 'i', "Perform case insensitive match.")
	defaultHide       = getopt.BoolLong("hide", 'n', "Hide all lines by default.")
	noSkipMarker      = getopt.BoolLong("no-skip-marker", 'S', "Suppress skip markers.")
	execute           = getopt.BoolLong("command", 'c', "TODO Doc") // "Treat arguments as command line instead of input filese, execute it and apply to output.\nOptionally specify command line terminator.")
	eatStderr         = getopt.BoolLong("stderr", '2', "Use with -c; process stderr from command too.")
	width             = getopt.IntLong("width", 'w', term.GetTermWidth(), "Set terminal width, used for pre and post lines.")
	cpuprofile        = getopt.StringLong("cpuprofile", 'P', "", "Write cpu profile to file.")
	help              = getopt.BoolLong("help", 'h', "Show this help.")
	noTtyWarning      = getopt.BoolLong("no-tty-warning", 'q', "Don't show warning even when stdin is tty.")
	readFiles         = getopt.BoolLong("files", 'f', "TODO Doc")
	argumentSeparator = getopt.StringLong("range-separator", 's', ArgumentSeparator, "Specify argument separator. (default="+ArgumentSeparator+")")
)

func init() {
	getopt.FlagLong(&util.Debug, "debug", 'd', "Enable debug output.")
	getopt.FlagLong(&matcher.NoPcre, "no-pcre", 'N', "Disable PCRE and use Go's regexp engine instead.")

	getopt.SetUsage(usage)
}

func usage() {
	os.Stderr.WriteString(`
hl: Versatile coloring filter

Basic usage: 
  hl [ -r RULE_TOML ] [OPTIONS] [ FILTER-SPEC... ] <FILE
    Read FILE and apply filters/colors.
    
    If FILTER-SPECs contain a range, or a -n option is given,
    it'll only print lines matching given PATTERNs, or
    lines in the ranges given by FILER-SPECs.
    When -A, -B or -C is given, "context" lines will also be
    printed.

  hl -f [ -r RULE_TOML ] [OPTIONS] FILE... , [ FILTER-SPEC... ]
    Read FILE(s) and apply filers/colors.
		
  hl -c [-2] [ -r RULE_TOML ] [OPTIONS] COMMAND [ARG...] , [ FILTER-SPEC... ]
    Execute COMMAND and apply filters/colors to its stdout output.
    If a -2 option is given, the stderr output will be processed too.

  FILTER-SPEC is a list of:
    PATTERN [ COLOR-SPEC ]
    PATTERN [ COLOR-SPEC ] ',' PATTERN [ COLOR-SPEC ]

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

func preprocessOptions() {
	if *help {
		getopt.Usage()
		os.Exit(0)
	}

	if *width > 0 {
		term.TermWidth = *width
	}
	if *context > 0 {
		*after = *context
		*before = *context
	}

	if *execute && *readFiles {
		Fatalf("Cannot use -c and -f at the same time.\n")
	}
}

func main() {
	getopt.Parse()

	preprocessOptions()

	// Initialize highlighter.
	h := highlighter.NewHighlighter()
	h.SetIgnoreCase(*ignoreCase)
	h.SetDefaultHide(*defaultHide)
	h.SetDefaultBefore(*before)
	h.SetDefaultAfter(*after)
	h.SetNoSkipMarker(*noSkipMarker)
	util.Dump("Highlighter (start): ", h)

	// Process -c and -f, and also extract simple (inline) rules.

	inputArgs, err := parseArgs(h, getopt.Args(), *execute || *readFiles, *argumentSeparator)
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

	// Main.
	if *readFiles {
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

		if *execute {
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

	err := h.NewRuntime(os.Stdout).ColorReader(rd /*callFinish*/, true)
	if err != nil {
		Fatalf("Unknown failure: %s", err)
	}
}

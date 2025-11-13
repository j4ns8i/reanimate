package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/j4ns8i/reanimate/pkg/horde"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// reanimateCmd represents the reanimate command
var reanimateCmd = &cobra.Command{
	Use:   "reanimate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: runReanimateCmd,
}

// PersistentFlags
var (
	oLog string
)

// Local flags
var (
	oFile string
	oList bool
)

func init() {
	// rootCmd.AddCommand(reanimateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	reanimateCmd.PersistentFlags().StringVar(&oLog, "log", "info", "Set logging level (trace, debug, info, warn, error)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	reanimateCmd.Flags().StringVarP(&oFile, "file", "f", "", "File for reading instructions")
	reanimateCmd.Flags().BoolVarP(&oList, "list", "l", false, "List hordes")
	cobra.OnInitialize(func() {
		if term.IsTerminal(int(os.Stdout.Fd())) {
			log.Logger = log.Output(zerolog.NewConsoleWriter())
		}
		if oLog != "" {
			level, err := zerolog.ParseLevel(oLog)
			if err != nil {
				log.Fatal().Err(err).Msg("parsing verbosity level")
			}
			zerolog.SetGlobalLevel(level)
		}
	})
}

func runReanimateCmd(cmd *cobra.Command, args []string) error {
	if oList {
		return runReanimateListCmd()
	}
	var inputInstructionsIO io.ReadCloser
	if oFile != "" {
		if oFile == "-" {
			inputInstructionsIO = os.Stdin
		} else {
			f, err := os.Open(oFile)
			if err != nil {
				return fmt.Errorf("opening input: %w", err)
			}
			inputInstructionsIO = f
		}
	}

	if oFile == "" && len(args) == 0 {
		return fmt.Errorf("no input instructions provided, use -f or provide as arguments")
	}

	err := start(inputInstructionsIO)
	if err != nil {
		return fmt.Errorf("starting reanimation process: %w", err)
	}

	return nil
}

func runReanimateListCmd() error {
	hordeReanimator := horde.NewTmux()
	hordes, err := hordeReanimator.List()
	if err != nil {
		return err
	}
	for _, h := range hordes {
		fmt.Println(h.Name())
	}
	return nil
}

func logErr(err error) {
	if err != nil {
		log.Warn().Err(err).Msg("received non-critical error")
	}
}

// TODO: test these with mocks for things like tmux sessions
// TODO: add listing sessions (tmux ls -f '#{m:<prefix>*,#S}')
func start(input io.ReadCloser) error {
	defer func() { logErr(input.Close()) }()

	horde, err := horde.NewTmux().Reanimate()
	if err != nil {
		return fmt.Errorf("reanimating horde: %w", err)
	}

	log.Info().Str("name", horde.Name()).Msg("reanimated horde")

	return nil
}

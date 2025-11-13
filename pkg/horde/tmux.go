package horde

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

const tmuxSessionPrefix = "reanimate-"

func NewTmux() Reanimator {
	return &TmuxReanimator{}
}

type TmuxReanimator struct {
}

type TmuxHorde struct {
	r               *TmuxReanimator
	tmuxSessionName string
}

func (t *TmuxReanimator) generateSessionName() string {
	id := ulid.MustNewDefault(time.Now())
	return tmuxSessionPrefix + id.String()
}

func (t *TmuxReanimator) Reanimate() (Horde, error) {
	sessionName := t.generateSessionName()
	cmd := exec.Command("tmux", "new-session", "-s", sessionName, "-d")
	log.Debug().Strs("cmd", cmd.Args).Msg("running command")
	_, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("creating new session: %w", err)
	}

	log.Info().Str("session_name", sessionName).Msg("created session")

	horde := &TmuxHorde{
		r:               t,
		tmuxSessionName: sessionName,
	}
	return horde, nil
}

func (t *TmuxReanimator) List() ([]Horde, error) {
	cmd := exec.Command("tmux", "list-session",
		"-f", fmt.Sprintf(`#{m:%s*,#S}`, tmuxSessionPrefix),
		"-F", "#S")

	log.Debug().Stringer("cmd", cmd).Msg("running command")
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("listing sessions: %w", err)
	}

	var hordes []Horde
	for sessionName := range strings.SplitSeq(string(b), "\n") {
		horde := &TmuxHorde{
			r:               t,
			tmuxSessionName: sessionName,
		}
		hordes = append(hordes, horde)
	}

	return hordes, nil
}

func (h *TmuxHorde) Destroy() error {
	cmd := exec.Command("tmux", "kill-session", "-t", h.tmuxSessionName)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("destroying session: %w", err)
	}
	return nil
}

func (h *TmuxHorde) Summon() error {
	cmd := exec.Command("tmux", "attach-session", "-t", h.tmuxSessionName)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("destroying session: %w", err)
	}
	return nil
}

func (h *TmuxHorde) Name() string {
	return h.tmuxSessionName
}

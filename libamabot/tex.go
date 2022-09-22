package libamabot

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gw31415/amabot/libamabot/ghostscript"
	"star-tex.org/x/tex"
)

func init() {
	slashCmd(&discordgo.ApplicationCommand{
		Name:        "tex",
		Description: "Render mathematical expressions using plain.tex",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "tex",
				Description: "TeX representation of the equation to be rendered",
				Required:    true,
			},
		},
	}, func(ctx context.Context, opts AmabotOptions, s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}
		tex := optionMap["tex"].StringValue()
		tex = fmt.Sprintf("\\nopagenumbers$$%s$$", tex)
		childCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		tmp_png, err := os.CreateTemp("", "tmp_png")
		if err != nil {
			panic(err)
		}
		defer tmp_png.Close()
		defer os.Remove(tmp_png.Name())
		reader := strings.NewReader(tex)
		err = tex2ps(tmp_png.Name(), os.Stderr, reader)
		if err != nil {
			panic(err)
		}
		select {
		case <-childCtx.Done():
			return
		default:
			name := filepath.Base(tmp_png.Name()) + ".png"
			e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "TeX output",
							Image: &discordgo.MessageEmbedImage{
								URL: "attachment://" + name,
							},
						},
					},
					Files: []*discordgo.File{
						{
							ContentType: "image/png",
							Reader:      tmp_png,
							Name:        name,
						},
					},
				},
			})
			if e != nil {
				panic(e)
			} else {
				return
			}
		}
	})
}

func tex2ps(file_path string, stderr io.Writer, stdin io.Reader) error {
	// [dvips] The standard input must be seekable, so it cannot be a pipe. If
	// you must use a pipe, write a shell script that copies the pipe output to
	// a temporary file and then points dvips at this file.
	tmp_dvi, err := os.CreateTemp("", "tmp_dvips")
	if err != nil {
		return err
	}
	defer tmp_dvi.Close()
	defer os.Remove(tmp_dvi.Name())
	dummy_r, _ := io.Pipe()
	tex2dvi := tex.NewEngine(stderr, dummy_r)
	err = tex2dvi.Process(tmp_dvi, stdin)
	if err != nil {
		return err
	}

	tmp_ps, err := os.CreateTemp("", "tmp_dvips")
	if err != nil {
		return err
	}
	defer tmp_ps.Close()
	defer os.Remove(tmp_ps.Name())
	dvips := exec.Command("dvips", "-f")
	dvips.Stdin = tmp_dvi
	dvips.Stdout = tmp_ps
	dvips.Stderr = stderr
	err = dvips.Run()
	if err != nil {
		return err
	}

	tmp_eps, err := os.CreateTemp("", "tmp_gs_ps2eps")
	if err != nil {
		return err
	}
	defer tmp_eps.Close()
	defer os.Remove(tmp_eps.Name())
	gs_ps2eps, _ := ghostscript.NewInstance()
	gs_ps2eps.Init([]string{
		"gs",
		"-q",
		"-dBATCH",
		"-dSAFER",
		"-dNOPAUSE",
		"-dEPSCrop",
		"-sDEVICE=eps2write",
		"-sOutputFile=" + tmp_eps.Name(),
		tmp_ps.Name(),
	})
	gs_ps2eps.Exit()
	gs_ps2eps.Destroy()

	gs_eps2png, _ := ghostscript.NewInstance()
	gs_eps2png.Init([]string{
		"gs",
		"-q",
		"-dBATCH",
		"-dSAFER",
		"-dNOPAUSE",
		"-dGraphicsAlphaBits=2",
		"-dTextAlphaBits=2",
		"-dDownScaleFactor=2",
		"-r750",
		"-dEPSCrop",
		"-sDEVICE=pnggray",
		"-sOutputFile=" + file_path,
		tmp_eps.Name(),
	})
	gs_eps2png.Exit()
	gs_eps2png.Destroy()

	return nil
}

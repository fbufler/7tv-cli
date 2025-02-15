package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/fbufler/7tv-cli/package/imagecat"
	"github.com/fbufler/7tv-cli/package/seventv"
	"github.com/spf13/cobra"
)

var supportedFormats = []seventv.FileFormat{seventv.GIF, seventv.PNG}

var cmd = &cobra.Command{
	Use:   "7tv [emote]",
	Short: "7tv cli",
	Args:  cobra.ExactArgs(1),
	Run:   exec,
}

func exec(cmd *cobra.Command, args []string) {
	emoteName := args[0]
	format := cmd.Flag("format").Value.String()
	scaler := cmd.Flag("scaler").Value.String()

	// validate format
	if _, ok := seventv.FileFormatReverseMap[format]; !ok {
		slog.Error(fmt.Sprintf("invalid format: %s", format))
		return
	}
	parsedFormat := seventv.FileFormatReverseMap[format]
	if err := supportedFormat(parsedFormat); err != nil {
		slog.Error(err.Error())
		return
	}

	if _, ok := seventv.FileScalerReverseMap[scaler]; !ok {
		slog.Error(fmt.Sprintf("invalid scaler: %s", scaler))
		return
	}
	parsedScaler := seventv.FileScalerReverseMap[scaler]

	u, err := url.Parse("https://7tv.io/v3")
	if err != nil {
		slog.Error(fmt.Sprintf("failed to parse url: %s", err))
		return
	}

	client := seventv.New(u)
	emoteUrl, err := client.GetEmoteURLByQuery(emoteName, parsedFormat, parsedScaler)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get emote url: %s", err))
		return
	}

	slog.Info(emoteUrl)
	// get emote image as bytes
	resp, err := http.Get(emoteUrl)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get emote image: %s", err))
		return
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to read response body: %s", err))
		return
	}

	imagecat.New().PrintImage(imageData)
}

func supportedFormat(format seventv.FileFormat) error {
	for _, t := range supportedFormats {
		if t == format {
			return nil
		}
	}
	return fmt.Errorf("unsupported format: %s, supported formats are %s", seventv.FileFormatMap[format], supportedFormats)
}

func main() {
	// add flags for format and scaler
	cmd.Flags().String("format", "GIF", "emote format, default GIF")
	cmd.Flags().String("scaler", "2x", "emote scaler, default 2x")
	cmd.Execute()
}

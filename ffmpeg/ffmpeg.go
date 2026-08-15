package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const FFMPEG = "ffmpeg"

// resolveFFmpeg returns the ffmpeg binary path. It prefers the FFMPEG_BIN
// environment variable (set in the Docker image) and falls back to resolving
// `ffmpeg` on PATH.
func resolveFFmpeg() string {
	if bin := os.Getenv("FFMPEG_BIN"); bin != "" {
		return bin
	}
	return FFMPEG
}

func GenerateVideo(imagePaths []string, outputPath string) error {
	if len(imagePaths) == 0 {
		return fmt.Errorf("at least one image is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}

	args := []string{"-y"}
	for _, imagePath := range imagePaths {
		args = append(args, "-loop", "1", "-framerate", "1", "-t", "2", "-i", imagePath)
	}

	filters := make([]string, 0, len(imagePaths)+1)
	inputs := make([]string, 0, len(imagePaths))
	for index := range imagePaths {
		filters = append(filters, fmt.Sprintf(
			"[%d:v]scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black,setsar=1,trim=duration=1,setpts=PTS-STARTPTS[v%d]",
			index,
			index,
		))
		inputs = append(inputs, fmt.Sprintf("[v%d]", index))
	}
	filters = append(filters, fmt.Sprintf("%sconcat=n=%d:v=1:a=0,setpts=N/TB[outv]", strings.Join(inputs, ""), len(imagePaths)))

	args = append(args,
		"-filter_complex", strings.Join(filters, ";"),
		"-map", "[outv]",
		"-c:v", "libx264",
		"-threads", "1",
		"-preset", "medium",
		"-crf", "18",
		"-pix_fmt", "yuv420p",
		"-r", "1",
		"-fps_mode", "cfr",
		"-g", "1",
		"-keyint_min", "1",
		"-sc_threshold", "0",
		"-bf", "0",
		"-frames:v", fmt.Sprint(len(imagePaths)),
		"-an",
		"-movflags", "+faststart",
		"-f", "mp4",
		outputPath,
	)

	command := exec.Command(resolveFFmpeg(), args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("run ffmpeg: %w", err)
	}
	return nil
}

package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const FFMPEG = "ffmpeg"
const FFMPEG_THREADS = "1"
const BatchSize = 4

// GenerateVideo 将图片序列编码为视频。图片分批编码为片段后无损合并。
func GenerateVideo(imagePaths []string, outputPath string) error {
	if len(imagePaths) == 0 {
		return fmt.Errorf("at least one image is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}

	if BatchSize <= 0 {
		return encodeImageBatch(imagePaths, outputPath)
	}

	var partPaths []string
	for start := 0; start < len(imagePaths); start += BatchSize {
		end := min(start+BatchSize, len(imagePaths))
		partPath := fmt.Sprintf("%s.part%d.mp4", outputPath, start/BatchSize)
		if err := encodeImageBatch(imagePaths[start:end], partPath); err != nil {
			return err
		}
		partPaths = append(partPaths, partPath)
	}
	defer func() {
		for _, p := range partPaths {
			os.Remove(p)
		}
	}()

	// 无损合并片段
	return concatParts(partPaths, outputPath)
}

// encodeImageBatch 把一批图片编码为一个片段(参数与原逻辑一致)。
func encodeImageBatch(imagePaths []string, outputPath string) error {
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
		"-threads", FFMPEG_THREADS,
		"-preset", "veryfast",
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

	command := exec.Command(FFMPEG, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("run ffmpeg: %w", err)
	}
	return nil
}

// concatParts 用 concat demuxer 无损合并片段(各片段参数一致)。
func concatParts(partPaths []string, outputPath string) error {
	listPath := outputPath + ".list"
	var b strings.Builder
	for _, p := range partPaths {
		abs, err := filepath.Abs(p)
		if err != nil {
			return err
		}
		escaped := strings.ReplaceAll(abs, "\\", "/")
		escaped = strings.ReplaceAll(escaped, "'", "'\\''")
		fmt.Fprintf(&b, "file '%s'\n", escaped)
	}
	if err := os.WriteFile(listPath, []byte(b.String()), 0644); err != nil {
		return err
	}
	defer os.Remove(listPath)

	args := []string{
		"-y",
		"-fflags", "+genpts",
		"-f", "concat",
		"-safe", "0",
		"-i", listPath,
		"-c", "copy",
		outputPath,
	}
	command := exec.Command(FFMPEG, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("concat parts: %w", err)
	}
	return nil
}

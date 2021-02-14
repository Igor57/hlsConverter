package hlsConverter

import (
	//"log"

	"os/exec"

	"github.com/Igor57/transcoder"
	ffmpeg "github.com/Igor57/transcoder/ffmpeg"
)

func Convert(input string, outputPath string, outputName string, partName string) (<-chan transcoder.Progress, *exec.Cmd, error) {

	ffmpegConf := &ffmpeg.Config{
		FfmpegBinPath:   "/usr/bin/ffmpeg",
		FfprobeBinPath:  "/usr/bin/ffprobe",
		ProgressEnabled: true,
	}
	flagUp := 1
	optsBeforeInput := ffmpeg.Options{
		Listen: &flagUp,
	}

	filterComplex := "[v:0]split=3[original][temp720p][temp540p];[temp720p]scale=w=1280:h=720[720p];[temp540p]scale=w=960:h=540[540p]"
	preset := "veryfast"
	keyframeInterval := 50
	scThreshold := 0
	maps := []string{
		"[original] -c:v:0 libx264 -b:v:0 6000k -maxrate:v:0 6600k -bufsize:v:0 9000k",
		"[720p] -c:v:1 libx264 -b:v:1 4000k -maxrate:v:1 4400k -bufsize:v:1 6000k",
		"[540p] -c:v:2 libx264 -b:v:2 2000k -maxrate:v:2 2200k -bufsize:v:2 3000k",
		"a:0",
		"a:0",
		"a:0",
	}
	outputFormat := "hls"
	hlsFlags := "append_list+omit_endlist+discont_start"
	hlsSegmentDuration := 4
	hlsPlaylistType := "event"
	hlsMasterPlaylistName := "index.m3u8"
	hlsSegmentFilename := outputPath + "stream_%v/data%06d.ts"
	useLocaltimeMkdir := 1
	varStreamMap := "v:0,a:0 v:1,a:1 v:2,a:2"
	optsAfterInput := ffmpeg.Options{
		FilterComplex:         &filterComplex,
		Preset:                &preset,
		KeyframeInterval:      &keyframeInterval,
		ScThreshold:           &scThreshold,
		Maps:                  maps,
		OutputFormat:          &outputFormat,
		HlsFlags:              &hlsFlags,
		HlsSegmentDuration:    &hlsSegmentDuration,
		HlsPlaylistType:       &hlsPlaylistType,
		HlsMasterPlaylistName: &hlsMasterPlaylistName,
		HlsSegmentFilename:    &hlsSegmentFilename,
		UseLocaltimeMkdir:     &useLocaltimeMkdir,
		VarStreamMap:          &varStreamMap,
	}

	progress, cmd, err := ffmpeg.
		New(ffmpegConf).
		Input(input).
		Output(outputPath+outputName).
		StartAndReturnCmd(optsBeforeInput, optsAfterInput)

	return progress, cmd, err
}

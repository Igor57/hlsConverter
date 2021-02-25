package hlsConverter

import (
	//"log"

	"os/exec"

	"github.com/Igor57/transcoder"
	ffmpeg "github.com/Igor57/transcoder/ffmpeg"
)

func Convert(input string, outputPath string, outputName string) (<-chan transcoder.Progress, *exec.Cmd, error) {

	ffmpegConf := &ffmpeg.Config{
		FfmpegBinPath:   "/usr/bin/ffmpeg",
		FfprobeBinPath:  "/usr/bin/ffprobe",
		ProgressEnabled: true,
		Verbose:         true,
		DirtyCMD:        true,
	}
	flagUp := 1
	optsBeforeInput := ffmpeg.Options{
		Listen: &flagUp,
	}

	filterComplex := "'[v:0]split=3[temp1080p][temp720p][temp480p];[temp1080p]scale=w=1920:h=1080[1080p];[temp720p]scale=w=1280:h=720[720p];[temp480p]scale=w=854:h=480[480p]'"
	preset := "veryfast"
	keyframeInterval := 50
	scThreshold := 0
	//videoCodec := "libx264"
	// maps := []string{
	// 	"v:0 -s:v:0 854x480 -b:v:0 800k -maxrate:v:0 900k -bufsize:v:0 1.8M",
	// 	"v:0 -s:v:1 1280x720 -b:v:1 1.2M -maxrate:v:1 1.8M -bufsize:v:1 3.6M",
	// 	"v:0 -s:v:2 1920x1080 -b:v:2 2.5M -maxrate:v:2 4M -bufsize:v:2 8M",
	// 	"a:0",
	// 	"a:0",
	// 	"a:0 -c:a aac -ac 2",
	// }
	maps := []string{
		"[1080p] -c:v:0 libx264 -b:v:0 3000k -maxrate:v:0 4500k -bufsize:v:0 9000k",
		"[720p] -c:v:1 libx264 -b:v:1 1200k -maxrate:v:1 1400k -bufsize:v:1 2800k",
		"[480p] -c:v:2 libx264 -b:v:2 700k -maxrate:v:2 900k -bufsize:v:2 1800k",
		"a:0 -c:a:0 aac -b:a:0 192k -ac 2",
		"a:0 -c:a:1 aac -b:a:1 192k -ac 2",
		"a:0 -c:a:2 aac -b:a:2 192k -ac 2",
	}
	outputFormat := "hls"
	hlsFlags := "append_list+omit_endlist+discont_start"
	hlsSegmentDuration := 4
	hlsPlaylistType := "event"
	hlsMasterPlaylistName := "index.m3u8"
	hlsSegmentFilename := outputPath + "stream_%v_data%06d.ts"
	//useLocaltimeMkdir := 1
	varStreamMap := "'v:0,a:0 v:1,a:1 v:2,a:2'"
	optsAfterInput := ffmpeg.Options{
		FilterComplex: &filterComplex,
		//VideoCodec:            &videoCodec,
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
		//UseLocaltimeMkdir:     &useLocaltimeMkdir,
		VarStreamMap: &varStreamMap,
	}

	progress, cmd, err := ffmpeg.
		New(ffmpegConf).
		Input(input).
		Output(outputPath+outputName).
		StartAndReturnCmd(optsBeforeInput, optsAfterInput)

	return progress, cmd, err
}

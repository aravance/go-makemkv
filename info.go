package mkv

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type InfoJob struct {
	device  Device
	options MkvOptions
}

func Info(device Device, opts MkvOptions) (*InfoJob, error) {
	return &InfoJob{
		device:  device,
		options: opts,
	}, nil
}

type DiscInfo struct {
	Titles []TitleInfo

	DiscType   string
	Name       string
	LangCode   string
	LangName   string
	VolumeName string
}

type TitleInfo struct {
	VideoStreams []VideoStreamInfo
	AudioStreams []AudioStreamInfo

	Name             string
	ChapterCount     int
	Duration         time.Duration
	FileSize         int64
	SourceFileName   string
	Segments         []int
	FileName         string
	MetadataLangCode string
	MetadataLangName string
}

type VideoStreamInfo struct {
	Name             string
	CodecId          string
	CodecShort       string
	CodecLong        string
	VideoSize        string
	AspectRatio      string
	FrameRate        string
	StreamFlags      int
	MetadataLangCode string
	MetadataLangName string
	ConversionType   string
}

type AudioStreamInfo struct {
	Name             string
	LangCode         string
	LangName         string
	CodecId          string
	CodecShort       string
	CodecLong        string
	BitRate          string
	ChannelCount     int
	SampleRate       int
	SampleSize       int
	StreamFlags      int
	MetadataLangCode string
	MetadataLangName string
	ConversionType   string
}

func (j *InfoJob) Run() (*DiscInfo, error) {
	dev := j.device.Type() + ":" + j.device.Device()
	options := append(j.options.toStrings(), []string{"info", dev}...)
	cmd := exec.Command("makemkvcon", options...)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	if discInfo, err := parseDiscInfo(scanner); err != nil {
		return nil, err
	} else {
		return &discInfo, nil
	}
}

func parseDiscInfo(scanner *bufio.Scanner) (DiscInfo, error) {
	// since SINFO contains both video and audio, we use these to keep track
	// of the index offset while parsing, so we can put them in separate slices
	streamIndices := make([]streamIndex, 0, 5)

	var discInfo DiscInfo
	for scanner.Scan() {
		line := scanner.Text()
		prefix, content, found := strings.Cut(line, ":")
		if !found {
			continue
		}

		switch prefix {
		case "DRV":
			continue
		case "MSG":
			continue

		case "TCOUNT":
			size, _ := strconv.Atoi(content)
			discInfo.Titles = make([]TitleInfo, 0, size)

		case "CINFO":
			attrId, _, value, ok := parseCinfo(content)
			if !ok {
				continue
			}
			switch attrId {
			case ap_iaType:
				discInfo.DiscType = value
			case ap_iaName:
				discInfo.Name = value
			case ap_iaMetadataLanguageCode:
				discInfo.LangCode = value
			case ap_iaMetadataLanguageName:
				discInfo.LangName = value
			case ap_iaVolumeName:
				discInfo.VolumeName = value
			}

		case "TINFO":
			titleId, attrId, _, value, ok := parseTinfo(content)
			if !ok {
				continue
			}
			switch attrId {
			case ap_iaName:
				discInfo.Titles[titleId].Name = value
			case ap_iaChapterCount:
				discInfo.Titles[titleId].ChapterCount, _ = strconv.Atoi(value)
			case ap_iaDuration:
				discInfo.Titles[titleId].Duration, _ = parseDuration(value)
			case ap_iaDiskSizeBytes:
				discInfo.Titles[titleId].FileSize, _ = strconv.ParseInt(value, 10, 64)
			case ap_iaSourceFileName:
				discInfo.Titles[titleId].SourceFileName = value
			case ap_iaSegmentsCount:
				if count, err := strconv.Atoi(value); err != nil {
					discInfo.Titles[titleId].Segments = make([]int, 0, count)
				}
			case ap_iaSegmentsMap:
				if discInfo.Titles[titleId].Segments != nil {
					for i, s := range strings.Split(value, ",") {
						discInfo.Titles[titleId].Segments[i], _ = strconv.Atoi(s)
					}
				}
			case ap_iaOutputFileName:
				discInfo.Titles[titleId].FileName = value
			case ap_iaMetadataLanguageCode:
				discInfo.Titles[titleId].MetadataLangCode = value
			case ap_iaMetadataLanguageName:
				discInfo.Titles[titleId].MetadataLangName = value
			}

		case "SINFO":
			titleId, streamId, attrId, _, value, ok := parseSinfo(content)
			if !ok {
				continue
			}
			var index streamIndex
			if len(streamIndices) >= streamId {
				index = streamIndices[streamId]
			}
			switch attrId {
			case ap_iaType:
				var i int
				if value == "Video" {
					i = len(discInfo.Titles[titleId].VideoStreams)
					discInfo.Titles[titleId].VideoStreams = append(discInfo.Titles[titleId].VideoStreams, VideoStreamInfo{})
				} else {
					i = len(discInfo.Titles[titleId].AudioStreams)
					discInfo.Titles[titleId].AudioStreams = append(discInfo.Titles[titleId].AudioStreams, AudioStreamInfo{})
				}
				streamIndices[streamId] = streamIndex{value, i}
			case ap_iaName:
				discInfo.Titles[titleId].getStream(index).SetName(value)
			case ap_iaLangCode:
				discInfo.Titles[titleId].getStream(index).SetLangCode(value)
			case ap_iaLangName:
				discInfo.Titles[titleId].getStream(index).SetLangName(value)
			case ap_iaCodecId:
				discInfo.Titles[titleId].getStream(index).SetCodecId(value)
			case ap_iaCodecShort:
				discInfo.Titles[titleId].getStream(index).SetCodecShort(value)
			case ap_iaCodecLong:
				discInfo.Titles[titleId].getStream(index).SetCodecLong(value)
			case ap_iaBitrate:
				discInfo.Titles[titleId].getStream(index).SetBitRate(value)
			case ap_iaAudioChannelsCount:
				i, _ := strconv.Atoi(value)
				discInfo.Titles[titleId].getStream(index).SetChannelCount(i)
			case ap_iaAudioSampleRate:
				i, _ := strconv.Atoi(value)
				discInfo.Titles[titleId].getStream(index).SetSampleRate(i)
			case ap_iaAudioSampleSize:
				i, _ := strconv.Atoi(value)
				discInfo.Titles[titleId].getStream(index).SetSampleSize(i)
			case ap_iaVideoSize:
				discInfo.Titles[titleId].getStream(index).SetVideoSize(value)
			case ap_iaVideoAspectRatio:
				discInfo.Titles[titleId].getStream(index).SetAspectRatio(value)
			case ap_iaVideoFrameRate:
				discInfo.Titles[titleId].getStream(index).SetFrameRate(value)
			case ap_iaStreamFlags:
				i, _ := strconv.Atoi(value)
				discInfo.Titles[titleId].getStream(index).SetStreamFlags(i)
			case ap_iaMetadataLanguageCode:
				discInfo.Titles[titleId].getStream(index).SetMetadataLangCode(value)
			case ap_iaMetadataLanguageName:
				discInfo.Titles[titleId].getStream(index).SetMetadataLangName(value)
			case ap_iaOutputConversionType:
				discInfo.Titles[titleId].getStream(index).SetConversionType(value)
			}
		}
	}

	return discInfo, nil
}

func parseDuration(value string) (time.Duration, error) {
	var h, m, s int
	if _, err := fmt.Sscanf(value, "%d:%d:%d", &h, &m, &s); err != nil {
		return 0, err
	}
	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second, nil
}

func cutInt(s string, sep string) (int, string, bool) {
	str, out, found := strings.Cut(s, sep)
	if !found {
		return 0, str, false
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return i, out, false
	}
	return i, out, true
}

func parseCinfo(content string) (attrId int, code int, value string, ok bool) {
	attrId, content, ok = cutInt(content, ",")
	if !ok {
		return attrId, code, value, ok
	}

	code, value, ok = cutInt(content, ",")
	if !ok {
		return attrId, code, value, ok
	}

	value = strings.Trim(value, `"`)
	return attrId, code, value, ok
}

func parseTinfo(content string) (titleId int, attrId int, code int, value string, ok bool) {
	titleId, content, ok = cutInt(content, ",")
	if !ok {
		return titleId, attrId, code, value, ok
	}

	attrId, content, ok = cutInt(content, ",")
	if !ok {
		return titleId, attrId, code, value, ok
	}

	code, value, ok = cutInt(content, ",")
	if !ok {
		return titleId, attrId, code, value, ok
	}

	value = strings.Trim(value, `"`)
	return titleId, attrId, code, value, ok
}

func parseSinfo(content string) (titleId int, streamId int, attrId int, code int, value string, ok bool) {
	titleId, content, ok = cutInt(content, ",")
	if !ok {
		return titleId, streamId, attrId, code, value, ok
	}

	streamId, content, ok = cutInt(content, ",")
	if !ok {
		return titleId, streamId, attrId, code, value, ok
	}

	attrId, content, ok = cutInt(content, ",")
	if !ok {
		return titleId, streamId, attrId, code, value, ok
	}

	code, value, ok = cutInt(content, ",")
	if !ok {
		return titleId, streamId, attrId, code, value, ok
	}

	value = strings.Trim(value, `"`)
	return titleId, streamId, attrId, code, value, ok
}

////////////////////////// apdefs.h //////////////////////////

const (
	ap_iaUnknown                      int = 0
	ap_iaType                             = 1
	ap_iaName                             = 2
	ap_iaLangCode                         = 3
	ap_iaLangName                         = 4
	ap_iaCodecId                          = 5
	ap_iaCodecShort                       = 6
	ap_iaCodecLong                        = 7
	ap_iaChapterCount                     = 8
	ap_iaDuration                         = 9
	ap_iaDiskSize                         = 10
	ap_iaDiskSizeBytes                    = 11
	ap_iaStreamTypeExtension              = 12
	ap_iaBitrate                          = 13
	ap_iaAudioChannelsCount               = 14
	ap_iaAngleInfo                        = 15
	ap_iaSourceFileName                   = 16
	ap_iaAudioSampleRate                  = 17
	ap_iaAudioSampleSize                  = 18
	ap_iaVideoSize                        = 19
	ap_iaVideoAspectRatio                 = 20
	ap_iaVideoFrameRate                   = 21
	ap_iaStreamFlags                      = 22
	ap_iaDateTime                         = 23
	ap_iaOriginalTitleId                  = 24
	ap_iaSegmentsCount                    = 25
	ap_iaSegmentsMap                      = 26
	ap_iaOutputFileName                   = 27
	ap_iaMetadataLanguageCode             = 28
	ap_iaMetadataLanguageName             = 29
	ap_iaTreeInfo                         = 30
	ap_iaPanelTitle                       = 31
	ap_iaVolumeName                       = 32
	ap_iaOrderWeight                      = 33
	ap_iaOutputFormat                     = 34
	ap_iaOutputFormatDescription          = 35
	ap_iaSeamlessInfo                     = 36
	ap_iaPanelText                        = 37
	ap_iaMkvFlags                         = 38
	ap_iaMkvFlagsText                     = 39
	ap_iaAudioChannelLayoutName           = 40
	ap_iaOutputCodecShort                 = 41
	ap_iaOutputConversionType             = 42
	ap_iaOutputAudioSampleRate            = 43
	ap_iaOutputAudioSampleSize            = 44
	ap_iaOutputAudioChannelsCount         = 45
	ap_iaOutputAudioChannelLayoutName     = 46
	ap_iaOutputAudioChannelLayout         = 47
	ap_iaOutputAudioMixDescription        = 48
	ap_iaComment                          = 49
	ap_iaOffsetSequenceId                 = 50
	ap_iaMaxValue
)

//////////////////////////// hack ////////////////////////////
// janky abstraction to simplify video/audio stream parsing //

type streamIndex struct {
	t string
	i int
}

func (t *TitleInfo) getStream(index streamIndex) iStreamInfo {
	if index.t == "Video" {
		return &t.VideoStreams[index.i]
	} else {
		return &t.AudioStreams[index.i]
	}
}

type iStreamInfo interface {
	SetName(string)
	SetLangCode(string)
	SetLangName(string)
	SetCodecId(string)
	SetCodecShort(string)
	SetCodecLong(string)
	SetBitRate(string)
	SetChannelCount(int)
	SetSampleRate(int)
	SetSampleSize(int)
	SetVideoSize(string)
	SetAspectRatio(string)
	SetFrameRate(string)
	SetStreamFlags(int)
	SetMetadataLangCode(string)
	SetMetadataLangName(string)
	SetConversionType(string)
}

func (v *VideoStreamInfo) SetName(name string) {
	v.Name = name
}

func (v *VideoStreamInfo) SetLangCode(langCode string) {
	// nop
}

func (v *VideoStreamInfo) SetLangName(langName string) {
	// nop
}

func (v *VideoStreamInfo) SetCodecId(codecId string) {
	v.CodecId = codecId
}

func (v *VideoStreamInfo) SetCodecShort(codecShort string) {
	v.CodecShort = codecShort
}

func (v *VideoStreamInfo) SetCodecLong(codecLong string) {
	v.CodecLong = codecLong
}

func (v *VideoStreamInfo) SetBitRate(bitRate string) {
	// nop
}

func (v *VideoStreamInfo) SetChannelCount(channelCount int) {
	// nop
}

func (v *VideoStreamInfo) SetSampleRate(sampleRate int) {
	// nop
}

func (v *VideoStreamInfo) SetSampleSize(sampleSize int) {
	// nop
}

func (v *VideoStreamInfo) SetVideoSize(videoSize string) {
	v.VideoSize = videoSize
}

func (v *VideoStreamInfo) SetAspectRatio(aspectRatio string) {
	v.AspectRatio = aspectRatio
}

func (v *VideoStreamInfo) SetFrameRate(frameRate string) {
	v.FrameRate = frameRate
}

func (v *VideoStreamInfo) SetStreamFlags(streamFlags int) {
	v.StreamFlags = streamFlags
}

func (v *VideoStreamInfo) SetMetadataLangCode(metadataLangCode string) {
	v.MetadataLangCode = metadataLangCode
}

func (v *VideoStreamInfo) SetMetadataLangName(metadataLangName string) {
	v.MetadataLangName = metadataLangName
}

func (v *VideoStreamInfo) SetConversionType(conversionType string) {
	v.ConversionType = conversionType
}

func (a *AudioStreamInfo) SetName(name string) {
	a.Name = name
}

func (a *AudioStreamInfo) SetLangCode(langCode string) {
	a.LangCode = langCode
}

func (a *AudioStreamInfo) SetLangName(langName string) {
	a.LangName = langName
}

func (a *AudioStreamInfo) SetCodecId(codecId string) {
	a.CodecId = codecId
}

func (a *AudioStreamInfo) SetCodecShort(codecShort string) {
	a.CodecShort = codecShort
}

func (a *AudioStreamInfo) SetCodecLong(codecLong string) {
	a.CodecLong = codecLong
}

func (a *AudioStreamInfo) SetBitRate(bitRate string) {
	a.BitRate = bitRate
}

func (a *AudioStreamInfo) SetChannelCount(channelCount int) {
	a.ChannelCount = channelCount
}

func (a *AudioStreamInfo) SetSampleRate(sampleRate int) {
	a.SampleRate = sampleRate
}

func (a *AudioStreamInfo) SetSampleSize(sampleSize int) {
	a.SampleSize = sampleSize
}

func (a *AudioStreamInfo) SetVideoSize(videoSize string) {
	// nop
}

func (a *AudioStreamInfo) SetAspectRatio(aspectRatio string) {
	// nop
}

func (a *AudioStreamInfo) SetFrameRate(frameRate string) {
	// nop
}

func (a *AudioStreamInfo) SetStreamFlags(streamFlags int) {
	a.StreamFlags = streamFlags
}

func (a *AudioStreamInfo) SetMetadataLangCode(metadataLangCode string) {
	a.MetadataLangCode = metadataLangCode
}

func (a *AudioStreamInfo) SetMetadataLangName(metadataLangName string) {
	a.MetadataLangName = metadataLangName
}

func (a *AudioStreamInfo) SetConversionType(conversionType string) {
	a.ConversionType = conversionType
}

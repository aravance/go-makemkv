package mkv

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDiscInfo(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	result, err := parseDiscInfo(scanner)
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, "DiscType", result.DiscType)
	assert.Equal(t, "DiscName", result.Name)
	assert.Equal(t, "LangCode", result.LangCode)
	assert.Equal(t, "LangName", result.LangName)
	assert.Equal(t, "VolumeName", result.VolumeName)
	assert.Equal(t, 3, len(result.Titles), "Titles length does not match")
	assertTitle(t, TitleInfo{
		VideoStreams:     make([]VideoStreamInfo, 1),
		AudioStreams:     make([]AudioStreamInfo, 2),
		Name:             "TitleName0",
		ChapterCount:     42,
		Duration:         1*time.Hour + 32*time.Minute + 31*time.Second,
		FileSize:         12345,
		SourceFileName:   "00000.mpls",
		Segments:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		FileName:         "TitleName0_t00.mkv",
		MetadataLangCode: "TitleLangCode0",
		MetadataLangName: "TitleLangName0",
	}, result.Titles[0])
	assertTitle(t, TitleInfo{
		VideoStreams:     make([]VideoStreamInfo, 1),
		AudioStreams:     make([]AudioStreamInfo, 1),
		Name:             "TitleName1",
		ChapterCount:     42,
		Duration:         1*time.Hour + 32*time.Minute + 31*time.Second,
		FileSize:         23456,
		SourceFileName:   "00002.mpls",
		Segments:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		FileName:         "TitleName1_t01.mkv",
		MetadataLangCode: "TitleLangCode1",
		MetadataLangName: "TitleLangName1",
	}, result.Titles[1])
	assertTitle(t, TitleInfo{
		VideoStreams:     make([]VideoStreamInfo, 1),
		AudioStreams:     make([]AudioStreamInfo, 1),
		Name:             "TitleName2",
		ChapterCount:     42,
		Duration:         1*time.Hour + 32*time.Minute + 31*time.Second,
		FileSize:         34567,
		SourceFileName:   "00001.mpls",
		Segments:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		FileName:         "TitleName2_t02.mkv",
		MetadataLangCode: "TitleLangCode2",
		MetadataLangName: "TitleLangName2",
	}, result.Titles[2])
}

func assertTitle(t *testing.T, expected TitleInfo, actual TitleInfo) {
	assert.Equal(t, len(expected.AudioStreams), len(actual.AudioStreams), "AudioStream length does not match")
	assert.Equal(t, len(expected.VideoStreams), len(actual.VideoStreams), "VideoStream length does not match")
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.ChapterCount, actual.ChapterCount)
	assert.Equal(t, expected.Duration, actual.Duration)
	assert.Equal(t, expected.FileSize, actual.FileSize)
	assert.Equal(t, expected.SourceFileName, actual.SourceFileName)
	// assert.Equal(t, expected.Segments, actual.Segments)
	assert.Equal(t, expected.FileName, actual.FileName)
	assert.Equal(t, expected.MetadataLangCode, actual.MetadataLangCode)
	assert.Equal(t, expected.MetadataLangName, actual.MetadataLangName)
}

const input = `
MSG:1005,0,1,"MakeMKV v1.17.6 linux(x64-release) started","%1 started","MakeMKV v1.17.6 linux(x64-release)"
DRV:0,0,999,0,"MyBluRayDrive","VolumeName","/dev/sr0"
DRV:1,256,999,0,"","",""
DRV:0,2,999,12,"MyBluRayDrive","DiscLabel","/dev/sr0"
MSG:3007,0,0,"Using direct disc access mode","Using direct disc access mode"
MSG:5085,0,0,"Loaded content hash table, will verify integrity of M2TS files.","Loaded content hash table, will verify integrity of M2TS files."
MSG:3025,16777216,3,"Title #00003.mpls has length of 8 seconds which is less than minimum title length of 3600 seconds and was therefore skipped","Title #%1 has length of %2 seconds which is less than minimum title length of %3 seconds and was therefore skipped","00003.mpls","8","3600"
MSG:3025,16777216,3,"Title #00004.m2ts has length of 5 seconds which is less than minimum title length of 3600 seconds and was therefore skipped","Title #%1 has length of %2 seconds which is less than minimum title length of %3 seconds and was therefore skipped","00004.m2ts","5","3600"
MSG:5011,0,0,"Operation successfully completed","Operation successfully completed"
TCOUNT:3
CINFO:1,6209,"DiscType"
CINFO:2,0,"DiscName"
CINFO:28,0,"LangCode"
CINFO:29,0,"LangName"
CINFO:30,0,"DiscTreeInfo"
CINFO:31,6119,"<b>Source information</b><br>"
CINFO:32,0,"VolumeName"
CINFO:33,0,"0"
TINFO:0,2,0,"TitleName0"
TINFO:0,8,0,"42"
TINFO:0,9,0,"1:32:31"
TINFO:0,10,0,"40.4 GB"
TINFO:0,11,0,"12345"
TINFO:0,16,0,"00000.mpls"
TINFO:0,25,0,"10"
TINFO:0,26,0,"1,2,3,4,5,6,7,8,9,10"
TINFO:0,27,0,"TitleName0_t00.mkv"
TINFO:0,28,0,"TitleLangCode0"
TINFO:0,29,0,"TitleLangName0"
TINFO:0,30,0,"TitleName0 - 42 chapter(s) , 40.4 GB"
TINFO:0,31,6120,"<b>Title information</b><br>"
TINFO:0,33,0,"0"
SINFO:0,0,1,6201,"Video"
SINFO:0,0,5,0,"V_MPEGH/ISO/HEVC"
SINFO:0,0,6,0,"MpegH"
SINFO:0,0,7,0,"MpegH HEVC Main10@L5.1"
SINFO:0,0,19,0,"3840x2160"
SINFO:0,0,20,0,"16:9"
SINFO:0,0,21,0,"23.976 (24000/1001)"
SINFO:0,0,22,0,"0"
SINFO:0,0,28,0,"StreamLangCode0"
SINFO:0,0,29,0,"StreamLangName0"
SINFO:0,0,30,0,"MpegH HEVC Main10@L5.1"
SINFO:0,0,31,6121,"<b>Track information</b><br>"
SINFO:0,0,33,0,"0"
SINFO:0,0,38,0,""
SINFO:0,0,42,5088,"ConversionType"
SINFO:0,1,1,6202,"Audio"
SINFO:0,1,2,0,"Surround 7.1"
SINFO:0,1,3,0,"eng"
SINFO:0,1,4,0,"English"
SINFO:0,1,5,0,"A_TRUEHD"
SINFO:0,1,6,0,"TrueHD"
SINFO:0,1,7,0,"TrueHD Atmos"
SINFO:0,1,13,0,"128 Kb/s"
SINFO:0,1,14,0,"8"
SINFO:0,1,17,0,"48000"
SINFO:0,1,18,0,"24"
SINFO:0,1,22,0,"1024"
SINFO:0,1,28,0,"StreamLangCode1"
SINFO:0,1,29,0,"StreamLangName1"
SINFO:0,1,30,0,"TrueHD Surround 7.1 English"
SINFO:0,1,31,6121,"<b>Track information</b><br>"
SINFO:0,1,33,0,"90"
SINFO:0,1,38,0,"d"
SINFO:0,1,39,0,"Default"
SINFO:0,1,40,0,"7.1"
SINFO:0,1,42,5088,"ConversionType"
SINFO:0,2,1,6202,"Audio"
SINFO:0,2,2,0,"Surround 5.1"
SINFO:0,2,3,0,"eng"
SINFO:0,2,4,0,"English"
SINFO:0,2,5,0,"A_AC3"
SINFO:0,2,6,0,"DD"
SINFO:0,2,7,0,"Dolby Digital"
SINFO:0,2,13,0,"640 Kb/s"
SINFO:0,2,14,0,"6"
SINFO:0,2,17,0,"48000"
SINFO:0,2,22,0,"2304"
SINFO:0,2,28,0,"StreamLangCode2"
SINFO:0,2,29,0,"StreamLangName2"
SINFO:0,2,30,0,"DD Surround 5.1 English"
SINFO:0,2,31,6121,"<b>Track information</b><br>"
SINFO:0,2,33,0,"90"
SINFO:0,2,38,0,""
SINFO:0,2,40,0,"5.1(side)"
SINFO:0,2,42,5088,"ConversionType"
SINFO:0,3,1,6203,"Subtitles"
SINFO:0,3,3,0,"eng"
SINFO:0,3,4,0,"English"
SINFO:0,3,5,0,"S_HDMV/PGS"
SINFO:0,3,6,0,"PGS"
SINFO:0,3,7,0,"HDMV PGS Subtitles"
SINFO:0,3,22,0,"0"
SINFO:0,3,28,0,"eng"
SINFO:0,3,29,0,"English"
SINFO:0,3,30,0,"PGS English"
SINFO:0,3,31,6121,"<b>Track information</b><br>"
SINFO:0,3,33,0,"90"
SINFO:0,3,38,0,""
SINFO:0,3,42,5088,"ConversionType"
SINFO:0,4,1,6203,"Subtitles"
SINFO:0,4,3,0,"eng"
SINFO:0,4,4,0,"English"
SINFO:0,4,5,0,"S_HDMV/PGS"
SINFO:0,4,6,0,"PGS"
SINFO:0,4,7,0,"HDMV PGS Subtitles"
SINFO:0,4,22,0,"6144"
SINFO:0,4,28,0,"eng"
SINFO:0,4,29,0,"English"
SINFO:0,4,30,0,"PGS English  (forced only)"
SINFO:0,4,31,6121,"<b>Track information</b><br>"
SINFO:0,4,33,0,"90"
SINFO:0,4,38,0,"d"
SINFO:0,4,39,0,"Default"
SINFO:0,4,42,5088,"ConversionType"
TINFO:1,2,0,"TitleName1"
TINFO:1,8,0,"42"
TINFO:1,9,0,"1:32:31"
TINFO:1,10,0,"40.3 GB"
TINFO:1,11,0,"23456"
TINFO:1,16,0,"00002.mpls"
TINFO:1,25,0,"9"
TINFO:1,26,0,"1,2,3,4,5,6,7,8,9"
TINFO:1,27,0,"TitleName1_t01.mkv"
TINFO:1,28,0,"TitleLangCode1"
TINFO:1,29,0,"TitleLangName1"
TINFO:1,30,0,"TitleName1 - 42 chapter(s) , 40.3 GB"
TINFO:1,31,6120,"<b>Title information</b><br>"
TINFO:1,33,0,"0"
SINFO:1,0,1,6201,"Video"
SINFO:1,0,5,0,"V_MPEGH/ISO/HEVC"
SINFO:1,0,6,0,"MpegH"
SINFO:1,0,7,0,"MpegH HEVC Main10@L5.1"
SINFO:1,0,19,0,"3840x2160"
SINFO:1,0,20,0,"16:9"
SINFO:1,0,21,0,"23.976 (24000/1001)"
SINFO:1,0,22,0,"0"
SINFO:1,0,28,0,"StreamLangCode0"
SINFO:1,0,29,0,"StreamLangName0"
SINFO:1,0,30,0,"MpegH HEVC Main10@L5.1"
SINFO:1,0,31,6121,"<b>Track information</b><br>"
SINFO:1,0,33,0,"0"
SINFO:1,0,38,0,""
SINFO:1,0,42,5088,"ConversionType"
SINFO:1,1,1,6202,"Audio"
SINFO:1,1,2,0,"Surround 7.1"
SINFO:1,1,3,0,"eng"
SINFO:1,1,4,0,"English"
SINFO:1,1,5,0,"A_TRUEHD"
SINFO:1,1,6,0,"TrueHD"
SINFO:1,1,7,0,"TrueHD Atmos"
SINFO:1,1,13,0,"128 Kb/s"
SINFO:1,1,14,0,"8"
SINFO:1,1,17,0,"48000"
SINFO:1,1,18,0,"24"
SINFO:1,1,22,0,"1024"
SINFO:1,1,28,0,"StreamLangCode1"
SINFO:1,1,29,0,"StreamLangName1"
SINFO:1,1,30,0,"TrueHD Surround 7.1 English"
SINFO:1,1,31,6121,"<b>Track information</b><br>"
SINFO:1,1,33,0,"90"
SINFO:1,1,38,0,"d"
SINFO:1,1,39,0,"Default"
SINFO:1,1,40,0,"7.1"
SINFO:1,1,42,5088,"ConversionType"
SINFO:1,2,1,6203,"Subtitles"
SINFO:1,2,3,0,"eng"
SINFO:1,2,4,0,"English"
SINFO:1,2,5,0,"S_HDMV/PGS"
SINFO:1,2,6,0,"PGS"
SINFO:1,2,7,0,"HDMV PGS Subtitles"
SINFO:1,2,22,0,"0"
SINFO:1,2,28,0,"eng"
SINFO:1,2,29,0,"English"
SINFO:1,2,30,0,"PGS English"
SINFO:1,2,31,6121,"<b>Track information</b><br>"
SINFO:1,2,33,0,"90"
SINFO:1,2,38,0,""
SINFO:1,2,42,5088,"ConversionType"
TINFO:2,2,0,"TitleName2"
TINFO:2,8,0,"42"
TINFO:2,9,0,"1:32:31"
TINFO:2,10,0,"40.3 GB"
TINFO:2,11,0,"34567"
TINFO:2,16,0,"00001.mpls"
TINFO:2,25,0,"8"
TINFO:2,26,0,"3,4,5,6,7,8,9,10"
TINFO:2,27,0,"TitleName2_t02.mkv"
TINFO:2,28,0,"TitleLangCode2"
TINFO:2,29,0,"TitleLangName2"
TINFO:2,30,0,"TitleName2 - 42 chapter(s) , 40.3 GB"
TINFO:2,31,6120,"<b>Title information</b><br>"
TINFO:2,33,0,"0"
SINFO:2,0,1,6201,"Video"
SINFO:2,0,5,0,"V_MPEGH/ISO/HEVC"
SINFO:2,0,6,0,"MpegH"
SINFO:2,0,7,0,"MpegH HEVC Main10@L5.1"
SINFO:2,0,19,0,"3840x2160"
SINFO:2,0,20,0,"16:9"
SINFO:2,0,21,0,"23.976 (24000/1001)"
SINFO:2,0,22,0,"0"
SINFO:2,0,28,0,"StreamLangCode0"
SINFO:2,0,29,0,"StreamLangName0"
SINFO:2,0,30,0,"MpegH HEVC Main10@L5.1"
SINFO:2,0,31,6121,"<b>Track information</b><br>"
SINFO:2,0,33,0,"0"
SINFO:2,0,38,0,""
SINFO:2,0,42,5088,"ConversionType"
SINFO:2,1,1,6202,"Audio"
SINFO:2,1,2,0,"Surround 7.1"
SINFO:2,1,3,0,"eng"
SINFO:2,1,4,0,"English"
SINFO:2,1,5,0,"A_TRUEHD"
SINFO:2,1,6,0,"TrueHD"
SINFO:2,1,7,0,"TrueHD Atmos"
SINFO:2,1,13,0,"128 Kb/s"
SINFO:2,1,14,0,"8"
SINFO:2,1,17,0,"48000"
SINFO:2,1,18,0,"24"
SINFO:2,1,22,0,"1024"
SINFO:2,1,28,0,"StreamLangCode1"
SINFO:2,1,29,0,"StreamLangName1"
SINFO:2,1,30,0,"TrueHD Surround 7.1 English"
SINFO:2,1,31,6121,"<b>Track information</b><br>"
SINFO:2,1,33,0,"90"
SINFO:2,1,38,0,"d"
SINFO:2,1,39,0,"Default"
SINFO:2,1,40,0,"7.1"
SINFO:2,1,42,5088,"ConversionType"
`

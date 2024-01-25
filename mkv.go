package makemkv

import (
	"bufio"
	"os/exec"
	"strconv"
	"strings"
)

type MkvJob struct {
	Statuschan  chan Status
	device      Device
	titleId     string
	destination string
	options     MkvOptions
}

func Mkv(device Device, titleId int, destination string, opts MkvOptions) *MkvJob {
	return &MkvJob{
		Statuschan:  nil,
		device:      device,
		titleId:     strconv.Itoa(titleId),
		destination: destination,
		options:     opts,
	}
}

func MkvAll(device Device, titleId int, destination string, opts MkvOptions) *MkvJob {
	return &MkvJob{
		Statuschan:  nil,
		device:      device,
		titleId:     "all",
		destination: destination,
		options:     opts,
	}
}

func (j *MkvJob) Run() error {
	dev := j.device.Type() + ":" + j.device.Device()
	options := append(j.options.toStrings(), []string{"mkv", dev, j.titleId, j.destination}...)
	cmd := exec.Command("makemkvcon", options...)

	var scanner bufio.Scanner
	if out, err := cmd.StdoutPipe(); err != nil {
		return err
	} else {
		scanner = *bufio.NewScanner(out)
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	var title string
	var channel string
	var total int
	var current int
	var max int

	for scanner.Scan() {
		line := scanner.Text()
		prefix, content, found := strings.Cut(line, ":")
		if !found {
			continue
		}

		parts := strings.Split(content, ",")
		switch prefix {
		case "PRGT":
			title = parts[2]
		case "PRGC":
			channel = parts[2]
		case "PRGV":
			current, _ = strconv.Atoi(parts[0])
			total, _ = strconv.Atoi(parts[1])
			max, _ = strconv.Atoi(parts[2])
			if j.Statuschan != nil {
				select {
				case j.Statuschan <- Status{
					Title:   title,
					Channel: channel,
					Current: current,
					Total:   total,
					Max:     max,
				}:
				}
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

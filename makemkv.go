package mkv

import (
	"strconv"
)

type Status struct {
	Title   string
	Channel string
	Current int
	Total   int
	Max     int
}

type MkvOptions struct {
	Messages  *string
	Progress  *string
	Debug     *string
	Directio  *bool
	Cache     *int
	Minlength *int
	Noscan    bool
	Decrypt   bool
}

func (m MkvOptions) toStrings() []string {
	result := []string{"-r"}
	if m.Messages != nil {
		result = append(result, "--messages="+*m.Messages)
	}
	if m.Progress != nil {
		result = append(result, "--progress="+*m.Progress)
	}
	if m.Debug != nil {
		result = append(result, "--debug="+*m.Debug)
	}
	if m.Directio != nil {
		result = append(result, "--directio="+strconv.FormatBool(*m.Directio))
	}
	if m.Cache != nil {
		result = append(result, "--cache="+strconv.Itoa(*m.Cache))
	}
	if m.Minlength != nil {
		result = append(result, "--minlength="+strconv.Itoa(*m.Minlength))
	}
	if m.Noscan {
		result = append(result, "--noscan")
	}
	if m.Decrypt {
		result = append(result, "--decrypt")
	}
	return result
}

func Stropt(s string) *string {
	return &s
}

func Intopt(i int) *int {
	return &i
}

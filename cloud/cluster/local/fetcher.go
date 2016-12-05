package local

import (
	"errors"
	"os/exec"
	"regexp"
	"strings"

	"github.com/scootdev/scoot/cloud/cluster"
)

// Poor man's dynamic localhost cluster nodes.
// Note: lsof is slow to return on osx so we just use 'ps' and regex match the port.
func MakeFetcher() cluster.Fetcher {
	return &localFetcher{}
}

type localFetcher struct{}

func (f *localFetcher) Fetch() (nodes []cluster.Node, err error) {
	var data []byte
	if data, err = f.fetchData(); err != nil {
		return nil, err
	} else if nodes, err = parseData(data); err != nil {
		return nil, err
	} else {
		return nodes, nil
	}
}

func (f *localFetcher) fetchData() ([]byte, error) {
	cmd := exec.Command("ps", "x")
	data, err := cmd.Output()
	return data, err
}

func parseData(data []byte) ([]cluster.Node, error) {
	nodes := []cluster.Node{}
	lines := string(data)
	for _, line := range strings.Split(lines, "\n") {
		thrift, err := parseFlag("thrift_addr", line)
		http, err2 := parseFlag("http_addr", line)
		if err == nil {
			if err2 == nil {
				http = "http://" + http
			}
			nodes = append(nodes, cluster.NewIdStatusNode(thrift, http))
		}
	}
	return nodes, nil
}

func parseFlag(flag, line string) (string, error) {
	// This is ugly but it works for now.
	re := regexp.MustCompile("workerserver.*" + flag + "(?: +|=)([^ ]*)")
	matches := re.FindStringSubmatch(line)
	if len(matches) == 2 {
		return matches[1], nil
	}
	return "", errors.New("Could not parse flag:'" + flag + "', from line:'" + line + "'")

}

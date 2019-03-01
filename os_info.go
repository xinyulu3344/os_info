package system

import (
	"vmstat/command"
	"fmt"
	"os"
	"bufio"
	"io"
	"regexp"
	"strings"
)

var _UNIXCONFDIR = "/etc"
var _distributor_id_file_re = regexp.MustCompile(`(?:DISTRIB_ID\s*=)\s*(.*)`)
var _release_file_re = regexp.MustCompile(`(?:DISTRIB_RELEASE\s*=)\s*(.*)`)
var _codename_file_re = regexp.MustCompile(`(?:DISTRIB_CODENAME\s*=)\s*(.*)`)
var _release_filename = regexp.MustCompile(`(\w+)[-_](release|version)`)
var _supported_dists = []string{"SuSE", "debian", "fedora", "redhat", "centos", "mandrake", "mandriva", "rocks", "slackware", "yellowdog", "gentoo", "UnitedLinux", "turbolinux", "arch", "mageia", "Ubuntu"}
var _lsb_release_version = regexp.MustCompile(`(.+) release ([\d.]+)[^(]*(?:\((.+)\))?`)
var _release_version = regexp.MustCompile(`([^0-9]+)(?: release )?([\d.]+)[^(]*(?:\((.+)\))?`)

type OsInfo struct {
	KernelVersion string
	PlatformType string
	PlatformVersion string
}

func NewOsInfo() *OsInfo {
	return &OsInfo{}
}

// 获取Linux内核版本
func (o *OsInfo) GetKernelVersion() string {
	content, ok := command.ExecCommand("uname", []string{"-r"})
	if !ok {
		fmt.Println("uname err")
	}
	o.KernelVersion = content[0]
	return o.KernelVersion
}

// 获取Linux发行版
func (o *OsInfo) GetLinuxDistribution() []string {
	result := make([]string, 0)
	var _u_distname string
	var _u_version string
	var _u_id string
	var element string

	filename := "/etc/lsb-release"
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err == nil {
		buf := bufio.NewReader(file)

		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				} else {
					fmt.Println("Read file error!", err)
				}
			}
			m := _distributor_id_file_re.FindString(line)
			if m != "" {
				_u_distname = strings.TrimSpace(strings.Split(m, "=")[1])
			}
			m = _release_file_re.FindString(line)
			if m != "" {
				_u_version = strings.TrimSpace(strings.Split(m, "=")[1])
			}
			m = _codename_file_re.FindString(line)
			if m != "" {
				_u_id = strings.TrimSpace(strings.Split(m, "=")[1])
			}
		}
		result = append(result, _u_distname, _u_version, _u_id)
		return result
	}
	defer file.Close()

	etc, err := command.ListDir(_UNIXCONFDIR, "")
	if err != nil {
		fmt.Println(err)
	}
	for _, element = range etc {
		m := _release_filename.FindStringSubmatch(element)
		if m != nil {
			if command.CheckStringInSlice(_supported_dists, m[1]) {
				break
			}
		}
	}

	firstline := command.ReadFirstLine(element)
	_u_distname, _u_version, _u_id = o._parse_release_file(firstline)

	if _u_distname != "" {
		result = append(result, _u_distname)
	}
	if _u_version != "" {
		result = append(result, _u_version)
	}
	if _u_id != "" {
		result = append(result, _u_id)
	}

	return result
}

func (o *OsInfo)_parse_release_file(firstline string) (string, string, string){
	version := ""
	id := ""

	// LSB format: "distro release x.x (codename)"
	m := _lsb_release_version.FindStringSubmatch(firstline)
	if m != nil {
		return m[1], m[2], m[3]
	}

	// Pre-LSB format: "distro x.x (codename)"
	m = _release_version.FindStringSubmatch(firstline)
	if m != nil {
		return m[1], m[2], m[3]
	}

	// Unknown format... take the first two words
	l := strings.Fields(firstline)
	if l != nil {
		version = l[0]
		if len(l) > 1 {
			id = l[1]
		}
	}
	return "", version, id
}

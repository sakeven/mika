package utils

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	// "log"
	// "os"
)

type ServerConf struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Method   string `json:"method"`
}

type Conf struct {
	Server      []*ServerConf `json:"server"`
	LocalAddr   string        `json:"local_addr"`
	LocalPort   int           `json:"local_port"`
	Timeout     int64         `json:"timeout"`
	TcpFastOpen bool          `json:"tcp_fastopen"`
}

func newConf() *Conf {
	c := &Conf{}
	c.Server = make([]*ServerConf, 1)
	c.Server[0] = new(ServerConf)
	return c
}

func ParseSeverConf() *Conf {
	var confFile string
	var conf = newConf()

	flag.StringVar(&confFile, "c", "", "path to config file")
	flag.StringVar(&conf.Server[0].Address, "s", "", "server address")
	flag.IntVar(&conf.Server[0].Port, "p", 8388, "server port")
	flag.StringVar(&conf.Server[0].Password, "k", "", "password")
	flag.StringVar(&conf.Server[0].Method, "m", "aes-256-cfb", "encryption method")
	flag.StringVar(&conf.LocalAddr, "b", "127.0.0.1", "local binding address")
	flag.IntVar(&conf.LocalPort, "l", 1080, "local port")
	flag.Int64Var(&conf.Timeout, "t", 300, "timeout in seconds")
	flag.BoolVar(&conf.TcpFastOpen, "-fast-open", false, "use TCP_FASTOPEN, requires Linux 3.7+")

	c, err := parseConf(confFile)
	if err != nil {
		return conf
	}
	return c
}

func parseConf(confFile string) (*Conf, error) {
	rawConf, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, err
	}
	v := &Conf{}

	json.Unmarshal(rawConf, v)

	return v, nil
}

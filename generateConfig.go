package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
)

func GenerateConfig(link string) string {
	parseUrl, err := url.Parse(link)
	if err != nil {
		//log.Println(err)

		return ""
	}
	var rc rawConfig

	switch parseUrl.Scheme {
	case "vless":
		parseVless(parseUrl, &rc)
	case "vmess":
		parseVmess(parseUrl, &rc)
	case "trojan":
		parseTrojan(parseUrl, &rc)
	case "ss":
		parseSS(parseUrl, &rc)
	case "ssr":
		parseSSR(parseUrl, &rc)
	default:
		return ""
	}

	var port, _ = freePorts.ExtractMin()
	jsonConf := buildConfig(port, &rc)
	hash := md5.Sum([]byte(jsonConf))
	hashString := hex.EncodeToString(hash[:])
	createConfigFile(hashString, jsonConf)
	return hashString
}

func parseVless(link *url.URL, rc *rawConfig) {
	// Implementation here
}

func parseVmess(link *url.URL, rc *rawConfig) {
	decoded, err := base64.StdEncoding.DecodeString(link.Host)
	if err != nil {
		return
	}
	var c VmessConfig
	err = json.Unmarshal(decoded, &c)
	rc.protocol = link.Scheme
	rc.ip = c.Add
	rc.id = c.ID
	rc.host = c.Host
	rc.port = c.Port
	rc.path = c.Path
	rc.net = c.Net
	rc.Type = c.Type
	rc.fp = c.Fp
	rc.security = c.Scy
	rc.sni = c.Sni
	rc.tls = c.Tls
}

func parseTrojan(link *url.URL, rc *rawConfig) {
	// Implementation here
}

func parseSS(link *url.URL, rc *rawConfig) {
	// Implementation here
}

func parseSSR(link *url.URL, rc *rawConfig) {
	// Implementation here
}

func buildConfig(port uint16, rc *rawConfig) string {
	conf := Config{
		Log: &LOG{Loglevel: "warning"},
		Dns: &DNS{
			DisableFallback: true,
			Servers: []Server{
				{Address: "https://8.8.8.8/dns-query"},
				{Address: "localhost", Domains: []string{"full:" + rc.host}},
			},
			Tag: "dns",
		},
		Stats: &STATS{},
		Routing: &Routing{
			DomainStrategy: "AsIs",
			Rules: []Rule{
				{InboundTag: []string{"socks-in"}, OutboundTag: "dns-out", Port: "53", Type: "field"},
				{OutboundTag: "proxy", Port: "0-65535", Type: "field"},
			},
		},
		Policy: &Policy{
			Levels: map[int]level{1: {30}},
			System: SYSTEM{
				StatsOutboundDownlink: true,
				StatsOutboundUplink:   true,
			},
		},
		Inbounds: &[]Inbound{
			{
				Listen:   AppConfig.Ip,
				Port:     0,
				Protocol: "socks",
				Settings: map[string]bool{"udp": true},
				Tag:      "socks-in",
				Sniffing: sniffing{
					Enabled:      true,
					MetadataOnly: false,
					RouteOnly:    true,
				},
			},
		},
		Outbounds: buildOutbounds(rc),
	}

	confJson, _ := json.Marshal(conf)
	//log.Println(string(confJson))
	return string(confJson)
}

func buildOutbounds(rc *rawConfig) *[]Outbound {
	outbounds := []Outbound{
		{
			DomainStrategy: "AsIs",
			Protocol:       rc.protocol,
			Settings: &OUTBOUNDSETTING{
				Vnext: []vnext{
					{
						Address: rc.ip,
						Port:    getPort(rc.port),
						Users: []USER{
							{Id: rc.id, AlterId: 0, Security: "auto"},
						},
					},
				},
			},
			StreamSettings: buildStreamSettings(rc),
			Tag:            "proxy",
		},
		{Protocol: "freedom", Tag: "direct"},
		{Protocol: "freedom", Tag: "bypass"},
		{Protocol: "blackhole", Tag: "block"},
		{
			Protocol:      "dns",
			ProxySettings: &PROXYSETTINGS{Tag: "proxy", TransportLayer: true},
			Settings:      &OUTBOUNDSETTING{Address: "8.8.8.8", Network: "tcp", Port: 53, UserLevel: 1},
			Tag:           "dns-out",
		},
	}
	return &outbounds
}

func buildStreamSettings(rc *rawConfig) *StreamSettingsObject {
	ss := &StreamSettingsObject{Network: rc.net}
	switch rc.net {
	case "ws":
		ss.WsSettings = &wsSettings{
			Path: rc.path,
			Host: rc.sni,
			Headers: map[string]string{
				"path": rc.path,
				"host": rc.host,
			},
		}
		getSecurity(rc, ss)
	default:
		return nil
	}
	return ss
}

func getPort(port uint16) uint16 {
	if port == 0 {
		return 443
	}
	return port
}

func isZeroValue(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func getSecurity(rc *rawConfig, s *StreamSettingsObject) {
	if rc.tls == "none" || rc.tls == "" {
		return
	}
	s.Security = rc.tls
	if !isZeroValue(s) {
		s.TlsSettings = &tlsSettings{ServerName: rc.sni}
	}
}

func createConfigFile(filename string, content string) {
	dir := "configs"

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return
		}
	}

	filePath := filepath.Join(dir, filepath.Base(filename))

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			return
		}
	} else {
		return
	}
}

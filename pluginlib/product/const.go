package product

import "github.com/hashicorp/go-plugin"

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

const PluginsMapHash = "product"

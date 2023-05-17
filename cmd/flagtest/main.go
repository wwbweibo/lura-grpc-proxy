package main

import "plugin"

func main() {
	p, _ := plugin.Open("plugin.so")
	f, _ := p.Lookup("PluginMain")
	f.(func())()
}

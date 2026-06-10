package modules

type Registry struct {
	modules []Module
}

func NewRegistry(modules []Module) *Registry {
	items := append([]Module(nil), modules...)
	for index := range items {
		if items[index].Route == "" {
			items[index].Route = "/" + items[index].ID
		}
		if items[index].Phase == "" {
			if items[index].Stub {
				items[index].Phase = "Stub"
			} else {
				items[index].Phase = "MVP"
			}
		}
		items[index].Implemented = !items[index].Stub
	}
	return &Registry{modules: items}
}

func DefaultRegistry() *Registry {
	return NewRegistry([]Module{
		{ID: "status", Name: "Status", Description: "Runtime host and module overview.", Category: "core", State: StateReady, Enabled: true, Endpoints: []string{"/api/status/host-overview", "/api/status/module-overview"}},
		{ID: "logs", Name: "Logs", Description: "Bounded in-memory application log view.", Category: "core", State: StateReady, Enabled: true, Endpoints: []string{"/api/logscenter/query"}},
		{ID: "settings", Name: "Settings", Description: "OpenLucky configuration storage.", Category: "core", State: StateReady, Enabled: true, Endpoints: []string{"/api/baseconfigure"}},
		{ID: "ddns", Name: "DDNS", Description: "MVP read-only task list placeholder.", Category: "network", State: StateReady, Enabled: true, Endpoints: []string{"/api/ddnstasklist"}},
		{ID: "web", Name: "Web Service", Description: "MVP read-only rule list placeholder.", Category: "network", State: StateReady, Enabled: true, Endpoints: []string{"/api/webservice/rules"}},
		{ID: "portforward", Name: "Port Forward", Description: "MVP read-only port-forward list placeholder.", Category: "network", State: StateReady, Enabled: true, Endpoints: []string{"/api/portforwards"}},
		{ID: "ssl", Name: "SSL", Description: "MVP read-only certificate list placeholder.", Category: "security", State: StateReady, Enabled: true, Endpoints: []string{"/api/ssl"}},
		{ID: "cron", Name: "Cron", Description: "MVP read-only schedule list placeholder.", Category: "automation", State: StateReady, Enabled: true, Endpoints: []string{"/api/cron/list"}},
		{ID: "docker", Name: "Docker", Description: "High-risk container controls are intentionally stubbed.", Category: "stub", State: StateStub, Stub: true},
		{ID: "webterminal", Name: "Web Terminal", Description: "Interactive terminal and SFTP are intentionally stubbed.", Category: "stub", State: StateStub, Stub: true},
		{ID: "webdav", Name: "WebDAV", Description: "File service implementation requires a separate threat model.", Category: "stub", State: StateStub, Stub: true},
		{ID: "smb", Name: "SMB", Description: "SMB implementation requires a separate threat model.", Category: "stub", State: StateStub, Stub: true},
		{ID: "ftpserver", Name: "FTP Server", Description: "FTP implementation requires a separate threat model.", Category: "stub", State: StateStub, Stub: true},
		{ID: "rclone", Name: "Rclone", Description: "Cloud storage integration is intentionally stubbed.", Category: "stub", State: StateStub, Stub: true},
		{ID: "cloudflared", Name: "Cloudflared", Description: "Tunnel control is intentionally stubbed.", Category: "stub", State: StateStub, Stub: true},
		{ID: "frp", Name: "FRP", Description: "Tunnel control is intentionally stubbed.", Category: "stub", State: StateStub, Stub: true},
		{ID: "coraza", Name: "Coraza WAF", Description: "WAF policy controls require a separate plan.", Category: "stub", State: StateStub, Stub: true},
	})
}

func (r *Registry) List() []Module {
	return append([]Module(nil), r.modules...)
}

func (r *Registry) Snapshot() Snapshot {
	return Snapshot{Modules: r.List()}
}

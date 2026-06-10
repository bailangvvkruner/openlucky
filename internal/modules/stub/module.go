package stub

func Response(module string, path string) map[string]any {
	return map[string]any{
		"module":  module,
		"path":    path,
		"status":  "not_implemented",
		"message": "This high-risk module is intentionally stubbed in the OpenLucky MVP.",
	}
}

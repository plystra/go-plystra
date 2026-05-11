package plystra

import "context"

func (s DataService) Tables(ctx context.Context) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/data/tables", nil)
}
func (s DataService) ListRows(ctx context.Context, resourceType string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/data/rows/"+esc(resourceType), query)
}
func (s DataService) GetRow(ctx context.Context, resourceType, resourceID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/data/rows/"+esc(resourceType)+"/"+esc(resourceID), nil)
}
func (s DataService) CreateRow(ctx context.Context, resourceType string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/data/rows/"+esc(resourceType), input)
}
func (s DataService) UpdateRow(ctx context.Context, resourceType, resourceID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/data/rows/"+esc(resourceType)+"/"+esc(resourceID), input)
}
func (s DataService) DeleteRow(ctx context.Context, resourceType, resourceID string, input Map) (Map, error) {
	return s.client.deleteMap(ctx, "/api/v1/data/rows/"+esc(resourceType)+"/"+esc(resourceID), input)
}

func (s PluginsService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/plugins", query)
}
func (s PluginsService) Get(ctx context.Context, pluginID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/plugins/"+esc(pluginID), nil)
}
func (s PluginsService) ValidateManifest(ctx context.Context, manifest any) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/plugins/validate-manifest", manifest)
}
func (s PluginsService) Install(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/plugins/install", input)
}
func (s PluginsService) Enable(ctx context.Context, pluginID string) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/plugins/"+esc(pluginID)+"/enable", Map{})
}
func (s PluginsService) Disable(ctx context.Context, pluginID string) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/plugins/"+esc(pluginID)+"/disable", Map{})
}
func (s PluginsService) Uninstall(ctx context.Context, pluginID string) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/plugins/"+esc(pluginID)+"/uninstall", Map{})
}
func (s PluginsService) Resources(ctx context.Context, pluginID string) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/plugins/"+esc(pluginID)+"/resources", nil)
}
func (s PluginsService) Permissions(ctx context.Context, pluginID string) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/plugins/"+esc(pluginID)+"/permissions", nil)
}
func (s PluginsService) AuditEvents(ctx context.Context, pluginID string) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/plugins/"+esc(pluginID)+"/audit-events", nil)
}
func (s PluginsService) AdminMenu(ctx context.Context, pluginID string) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/plugins/"+esc(pluginID)+"/admin-menu", nil)
}
func (s PluginsService) Settings(ctx context.Context, pluginID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/plugins/"+esc(pluginID)+"/settings", query)
}
func (s PluginsService) UpdateSettings(ctx context.Context, pluginID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/plugins/"+esc(pluginID)+"/settings", input)
}

func (s TemplatesService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/templates", query)
}
func (s TemplatesService) Get(ctx context.Context, templateID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/templates/"+esc(templateID), nil)
}
func (s TemplatesService) PreviewInstall(ctx context.Context, templateID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/templates/"+esc(templateID)+"/preview-install", input)
}
func (s TemplatesService) Install(ctx context.Context, templateID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/templates/"+esc(templateID)+"/install", input)
}

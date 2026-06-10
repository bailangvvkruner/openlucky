package app

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	hertzapp "github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/openlucky/openlucky/internal/auth"
	"github.com/openlucky/openlucky/internal/config"
	"github.com/openlucky/openlucky/internal/httpx"
	"github.com/openlucky/openlucky/internal/modules"
	logsmod "github.com/openlucky/openlucky/internal/modules/logs"
	statusmod "github.com/openlucky/openlucky/internal/modules/status"
	stubmod "github.com/openlucky/openlucky/internal/modules/stub"
	"github.com/openlucky/openlucky/web"
)

type Options struct {
	Addr          string
	ConfigPath    string
	AdminUsername string
	AdminPassword string
	TokenTTL      time.Duration
}

type App struct {
	Engine   *server.Hertz
	config   *config.Store
	auth     *auth.Service
	registry *modules.Registry
	logs     *logsmod.Store
	started  time.Time
}

func New(options Options) (*App, error) {
	if options.Addr == "" {
		options.Addr = "127.0.0.1:16601"
	}
	if options.ConfigPath == "" {
		options.ConfigPath = "openlucky.json"
	}
	store := config.NewStore(options.ConfigPath)
	if _, err := store.Load(); err != nil {
		return nil, err
	}

	openLucky := &App{
		Engine:   server.Default(server.WithHostPorts(options.Addr)),
		config:   store,
		auth:     auth.New(options.AdminUsername, options.AdminPassword, options.TokenTTL),
		registry: modules.DefaultRegistry(),
		logs:     logsmod.NewStore(512),
		started:  time.Now(),
	}
	openLucky.logs.Add("info", "app", "OpenLucky initialized")
	openLucky.registerRoutes()
	return openLucky, nil
}

func (a *App) Spin() {
	a.Engine.Spin()
}

func (a *App) registerRoutes() {
	a.Engine.GET("/healthz", a.handleHealth)
	a.Engine.POST("/api/login/challenge", a.handleLoginChallenge)
	a.Engine.POST("/api/login", a.handleLogin)
	a.registerStaticRoutes()

	api := a.Engine.Group("/api", a.authMiddleware)
	api.POST("/logout", a.handleLogout)
	api.GET("/status/host-overview", a.handleHostOverview)
	api.GET("/status/module-overview", a.handleModuleOverview)
	api.GET("/logscenter/query", a.handleLogs)
	api.GET("/baseconfigure", a.handleGetConfig)
	api.PUT("/baseconfigure", a.handlePutConfig)
	api.GET("/modules/list", a.handleModules)
	api.GET("/ddnstasklist", a.handleEmptyList("ddns"))
	api.GET("/webservice/rules", a.handleEmptyList("webservice"))
	api.GET("/portforwards", a.handleEmptyList("portforward"))
	api.GET("/ssl", a.handleEmptyList("ssl"))
	api.GET("/cron/list", a.handleEmptyList("cron"))

	for _, module := range []string{"docker", "webterminal", "webdav", "smb", "ftpserver", "filebrowser", "dlnaservice", "storagemanagement", "rclone", "cloudflared", "frp", "stun", "coraza", "ipdb", "thirdPartyAuthManager", "dbbackup", "wanjiadmin"} {
		api.GET("/"+module, a.handleStub(module))
		api.GET("/"+module+"/*path", a.handleStub(module))
	}
}

func (a *App) registerStaticRoutes() {
	a.Engine.GET("/lucky/", serveStatic)
	a.Engine.GET("/lucky/*filepath", serveStatic)
}

func serveStatic(ctx context.Context, c *hertzapp.RequestContext) {
	filePath := strings.TrimPrefix(c.Param("filepath"), "/")
	data, contentType, err := web.Read(filePath)
	if err != nil {
		c.JSON(consts.StatusNotFound, httpx.Error("not_found", "static asset not found"))
		return
	}
	c.Data(consts.StatusOK, contentType, data)
}

func (a *App) authMiddleware(ctx context.Context, c *hertzapp.RequestContext) {
	token := tokenFromRequest(c)
	if _, err := a.auth.Validate(token); err != nil {
		c.AbortWithStatusJSON(consts.StatusUnauthorized, httpx.Error("unauthorized", "login required"))
		return
	}
	c.Next(ctx)
}

func tokenFromRequest(c *hertzapp.RequestContext) string {
	for _, header := range []string{"OpenLucky-Admin-Token", "Lucky-Admin-Token"} {
		if value := strings.TrimSpace(string(c.GetHeader(header))); value != "" {
			return value
		}
	}
	authorization := strings.TrimSpace(string(c.GetHeader("Authorization")))
	if strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
		return strings.TrimSpace(authorization[len("bearer "):])
	}
	return ""
}

func (a *App) handleHealth(ctx context.Context, c *hertzapp.RequestContext) {
	c.JSON(consts.StatusOK, httpx.OK(map[string]any{"service": "openlucky", "status": "ok"}))
}

func (a *App) handleLoginChallenge(ctx context.Context, c *hertzapp.RequestContext) {
	c.JSON(consts.StatusOK, httpx.OK(a.auth.Challenge()))
}

func (a *App) handleLogin(ctx context.Context, c *hertzapp.RequestContext) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := readJSON(c, &request); err != nil {
		c.JSON(consts.StatusBadRequest, httpx.Error("bad_request", "invalid login payload"))
		return
	}
	session, err := a.auth.Login(request.Username, request.Password)
	if errors.Is(err, auth.ErrInvalidCredentials) {
		c.JSON(consts.StatusUnauthorized, httpx.Error("invalid_credentials", "invalid username or password"))
		return
	}
	if err != nil {
		c.JSON(consts.StatusInternalServerError, httpx.Error("login_failed", "could not create session"))
		return
	}
	a.logs.Add("info", "auth", "admin login succeeded")
	c.JSON(consts.StatusOK, httpx.OK(session))
}

func (a *App) handleLogout(ctx context.Context, c *hertzapp.RequestContext) {
	a.auth.Logout(tokenFromRequest(c))
	c.JSON(consts.StatusOK, httpx.OK(map[string]bool{"loggedOut": true}))
}

func (a *App) handleHostOverview(ctx context.Context, c *hertzapp.RequestContext) {
	c.JSON(consts.StatusOK, httpx.OK(statusmod.Host(a.started)))
}

func (a *App) handleModuleOverview(ctx context.Context, c *hertzapp.RequestContext) {
	c.JSON(consts.StatusOK, httpx.OK(a.registry.Snapshot()))
}

func (a *App) handleLogs(ctx context.Context, c *hertzapp.RequestContext) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	c.JSON(consts.StatusOK, httpx.OK(map[string]any{"entries": a.logs.Query(limit)}))
}

func (a *App) handleGetConfig(ctx context.Context, c *hertzapp.RequestContext) {
	cfg, err := a.config.Load()
	if err != nil {
		c.JSON(consts.StatusInternalServerError, httpx.Error("config_load_failed", "could not load config"))
		return
	}
	c.JSON(consts.StatusOK, httpx.OK(cfg))
}

func (a *App) handlePutConfig(ctx context.Context, c *hertzapp.RequestContext) {
	var cfg config.Config
	if err := readJSON(c, &cfg); err != nil {
		c.JSON(consts.StatusBadRequest, httpx.Error("bad_request", "invalid config payload"))
		return
	}
	if err := a.config.Save(cfg); err != nil {
		c.JSON(consts.StatusInternalServerError, httpx.Error("config_save_failed", "could not save config"))
		return
	}
	a.logs.Add("info", "settings", "configuration updated")
	c.JSON(consts.StatusOK, httpx.OK(cfg))
}

func (a *App) handleModules(ctx context.Context, c *hertzapp.RequestContext) {
	c.JSON(consts.StatusOK, httpx.OK(a.registry.Snapshot()))
}

func (a *App) handleEmptyList(module string) hertzapp.HandlerFunc {
	return func(ctx context.Context, c *hertzapp.RequestContext) {
		c.JSON(consts.StatusOK, httpx.OK(map[string]any{"module": module, "items": []any{}}))
	}
}

func (a *App) handleStub(module string) hertzapp.HandlerFunc {
	return func(ctx context.Context, c *hertzapp.RequestContext) {
		c.JSON(consts.StatusNotImplemented, httpx.Error("not_implemented", stubmod.Response(module, c.Param("path"))["message"].(string)))
	}
}

func readJSON(c *hertzapp.RequestContext, target any) error {
	body := c.Request.Body()
	if len(body) == 0 {
		return errors.New("empty body")
	}
	return json.Unmarshal(body, target)
}

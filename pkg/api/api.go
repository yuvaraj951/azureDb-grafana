package api

import (
  "github.com/go-macaron/binding"
  "github.com/grafana/grafana/pkg/api/avatar"
  "github.com/grafana/grafana/pkg/api/dtos"
  "github.com/grafana/grafana/pkg/api/live"
  "github.com/grafana/grafana/pkg/middleware"
  m "github.com/grafana/grafana/pkg/models"
  "gopkg.in/macaron.v1"
)

// Register adds http routes
func Register(r *macaron.Macaron) {
  reqSignedIn := middleware.Auth(&middleware.AuthOptions{ReqSignedIn: true})
  reqGrafanaAdmin := middleware.Auth(&middleware.AuthOptions{ReqSignedIn: true, ReqGrafanaAdmin: true})
  reqEditorRole := middleware.RoleAuth(m.ROLE_EDITOR, m.ROLE_ADMIN)
  reqOrgAdmin := middleware.RoleAuth(m.ROLE_ADMIN)
  quota := middleware.Quota
  bind := binding.Bind

  // not logged in views
  r.Get("/", reqSignedIn, Index)
  r.Get("/logout", Logout)
  r.Post("/login", quota("session"), bind(dtos.LoginCommand{}), wrap(LoginPost))
  r.Get("/login/:name", quota("session"), OAuthLogin)
  r.Get("/login", LoginView)
  r.Get("/invite/:code", Index)

  // authed views
  r.Get("/profile/", reqSignedIn, Index)
  r.Get("/profile/password", reqSignedIn, Index)
  r.Get("/profile/switch-org/:id", reqSignedIn, ChangeActiveOrgAndRedirectToHome)
  r.Get("/org/", reqSignedIn, Index)
  r.Get("/org/new", reqSignedIn, Index)
  r.Get("/datasources/", reqSignedIn, Index)
  r.Get("/datasources/new", reqSignedIn, Index)
  r.Get("/datasources/edit/*", reqSignedIn, Index)
  r.Get("/org/users/", reqSignedIn, Index)
  r.Get("/org/process/", reqSignedIn, Index)
  r.Get("/org/alerts/", reqSignedIn, Index)
  r.Get("/org/process/edit/*", reqSignedIn, Index)
  r.Get("/org/machine/", reqSignedIn, Index)
  r.Get("/org/apikeys/", reqSignedIn, Index)
  r.Get("/dashboard/import/", reqSignedIn, Index)
  r.Get("/admin", reqGrafanaAdmin, Index)
  r.Get("/admin/settings", reqGrafanaAdmin, Index)
  r.Get("/admin/users", reqGrafanaAdmin, Index)
  r.Get("/admin/users/create", reqGrafanaAdmin, Index)
  r.Get("/admin/users/edit/:id", reqGrafanaAdmin, Index)
  r.Get("/admin/orgs", reqGrafanaAdmin, Index)
  r.Get("/admin/orgs/edit/:id", reqGrafanaAdmin, Index)
  r.Get("/admin/stats", reqGrafanaAdmin, Index)

  r.Get("/styleguide", reqSignedIn, Index)

  r.Get("/plugins", reqSignedIn, Index)
  r.Get("/plugins/:id/edit", reqSignedIn, Index)
  r.Get("/plugins/:id/page/:page", reqSignedIn, Index)

  r.Get("/dashboard/*", reqSignedIn, Index)
  r.Get("/dashboard-solo/*", reqSignedIn, Index)
  r.Get("/import/dashboard", reqSignedIn, Index)
  r.Get("/dashboards/*", reqSignedIn, Index)

  r.Get("/playlists/", reqSignedIn, Index)
  r.Get("/playlists/*", reqSignedIn, Index)

  // sign up
  r.Get("/signup", Index)
  r.Get("/api/user/signup/options", wrap(GetSignUpOptions))
  r.Post("/api/user/signup", quota("user"), bind(dtos.SignUpForm{}), wrap(SignUp))
  r.Post("/api/user/signup/step2", bind(dtos.SignUpStep2Form{}), wrap(SignUpStep2))

  // invited
  r.Get("/api/user/invite/:code", wrap(GetInviteInfoByCode))
  r.Post("/api/user/invite/complete", bind(dtos.CompleteInviteForm{}), wrap(CompleteInvite))

  // reset password
  r.Get("/user/password/send-reset-email", Index)
  r.Get("/user/password/reset", Index)

  r.Post("/api/user/password/send-reset-email", bind(dtos.SendResetPasswordEmailForm{}), wrap(SendResetPasswordEmail))
  r.Post("/api/user/password/reset", bind(dtos.ResetUserPasswordForm{}), wrap(ResetPassword))

  // dashboard snapshots
  r.Get("/dashboard/snapshot/*", Index)
  r.Get("/dashboard/snapshots/", reqSignedIn, Index)

  // api for dashboard snapshots
  r.Post("/api/snapshots/", bind(m.CreateDashboardSnapshotCommand{}), CreateDashboardSnapshot)
  r.Get("/api/snapshot/shared-options/", GetSharingOptions)
  r.Get("/api/snapshots/:key", GetDashboardSnapshot)
  r.Get("/api/snapshots-delete/:key", reqEditorRole, DeleteDashboardSnapshot)

  // api renew session based on remember cookie
  r.Get("/api/login/ping", quota("session"), LoginApiPing)

  // authed api
  r.Group("/api", func() {

    // user (signed in)
    r.Group("/user", func() {
      r.Get("/", wrap(GetSignedInUser))
      r.Put("/", bind(m.UpdateUserCommand{}), wrap(UpdateSignedInUser))
      r.Post("/using/:id", wrap(UserSetUsingOrg))
      r.Get("/orgs", wrap(GetSignedInUserOrgList))

      r.Post("/stars/dashboard/:id", wrap(StarDashboard))
      r.Delete("/stars/dashboard/:id", wrap(UnstarDashboard))

      r.Put("/password", bind(m.ChangeUserPasswordCommand{}), wrap(ChangeUserPassword))
      r.Get("/quotas", wrap(GetUserQuotas))

      r.Get("/preferences", wrap(GetUserPreferences))
      r.Put("/preferences", bind(dtos.UpdatePrefsCmd{}), wrap(UpdateUserPreferences))
    })
    //
    //r.Get("api/org/process", wrap(GetProcessForCurrentOrg))




    // users (admin permission required)
    r.Group("/users", func() {
      r.Get("/", wrap(SearchUsers))
      r.Get("/:id", wrap(GetUserById))
      r.Get("/:id/orgs", wrap(GetUserOrgList))
      r.Put("/:id", bind(m.UpdateUserCommand{}), wrap(UpdateUser))
      r.Post("/:id/using/:orgId", wrap(UpdateUserActiveOrg))
    }, reqGrafanaAdmin)

    // org information available to all users.
    r.Group("/org", func() {
      r.Get("/", wrap(GetOrgCurrent))
      r.Get("/quotas", wrap(GetOrgQuotas))
    })

    // current org
    r.Group("/org", func() {
      r.Put("/", bind(dtos.UpdateOrgForm{}), wrap(UpdateOrgCurrent))
      r.Put("/address", bind(dtos.UpdateOrgAddressForm{}), wrap(UpdateOrgAddressCurrent))
      r.Post("/users", quota("user"), bind(m.AddOrgUserCommand{}), wrap(AddOrgUserToCurrentOrg))
      r.Get("/users", wrap(GetOrgUsersForCurrentOrg))
      r.Patch("/users/:userId", bind(m.UpdateOrgUserCommand{}), wrap(UpdateOrgUserForCurrentOrg))
      r.Delete("/users/:userId", wrap(RemoveOrgUserForCurrentOrg))

      // invites
      r.Get("/invites", wrap(GetPendingOrgInvites))
      r.Post("/invites", quota("user"), bind(dtos.AddInviteForm{}), wrap(AddOrgInvite))
      r.Patch("/invites/:code/revoke", wrap(RevokeInvite))

      // prefs
      r.Get("/preferences", wrap(GetOrgPreferences))
      r.Put("/preferences", bind(dtos.UpdatePrefsCmd{}), wrap(UpdateOrgPreferences))


      // process
      //r.Get("/process", wrap(GetProcess))
      r.Get("/process", wrap(GetProcessForCurrentOrg))
      r.Post("/process", quota("process"),bind(m.AddProcessCommand{}), wrap(AddProcessToCurrentOrg))
      r.Post("/parent", quota("process"),bind(m.AddProcessCommand{}), wrap(AddProcess))
      r.Delete("/process/:processId",wrap(RemoveProcessCurrentOrg))
      //r.Get("/process/:processId",wrap(GetProcessById))
      r.Get("/process/edit/:processId",wrap(GetProcessById))
      r.Get("/parent",wrap(GetProcessByParentName))
      //r.Patch("/process/:processId", bind(dtos.UpdateProcessForm{}), wrap(UpdateProcess))
      //r.Put("/process",bind(dtos))
      r.Patch("/process/:processId", bind(dtos.UpdateProcessForm{}), wrap(UpdateProcess))

      //Sub Process
      r.Get("/subprocess", wrap(GetSubProcessForCurrentOrg))
      r.Post("/subprocess", quota("subprocess"),bind(m.AddSubProcessCommand{}), wrap(AddSubProcessToCurrentOrg))
      r.Delete("/subprocess/:subProcessId",wrap(RemoveSubProcessCurrentOrg))
      r.Patch("/subprocess/:subProcessId", bind(dtos.UpdateSubProcessForm{}), wrap(UpdateSubProcess))
      r.Get("/subprocess/edit/:subProcessId",wrap(GetSubProcessById))
      r.Get("/subprocess/get/:processName",wrap(GetSubProcessByName))

      // machine
      r.Get("/machine",wrap(GetMachineForCurrentOrg))
      r.Post("/machine", quota("machines"),bind(m.AddMachineCommand{}), wrap(AddMachineToCurrentOrg))
      r.Delete("/machine/:machineId",wrap(RemoveMachineCurrentOrg))
      r.Patch("/machine/:machineId", bind(dtos.UpdateMachineForm{}), wrap(UpdateMachine))
      r.Get("/machine/:machineId",wrap(GetMachineById))
      r.Get("/machine/edit/:machineId",wrap(GetMachineById))
      //maintenance plan
      r.Get("/maintenance",wrap(GetMaintenanceForCurrentOrg))
      r.Post("/maintenance", quota("maintenance"),bind(m.AddMaintenanceCommand{}), wrap(AddMaintenanceToCurrentOrg))
      r.Delete("/maintenance/:Id",wrap(RemoveMaintenanceCurrentOrg))
      r.Get("/maintenance/edit/:Id",wrap(GetMaintenacneById))
      r.Get("/maintenance/:Id",wrap(GetMaintenacneById))
      r.Patch("/maintenance/:Id", bind(dtos.UpdateMaintenanceForm{}), wrap(UpdateMaintenance))


      //Alert History
      r.Get("/alerts/pending",wrap(GetPendingAlertHistory))
      r.Get("/alerts/pending/action/:id",wrap(GetPendingAlertActionHistory))
      r.Get("/alerts/completed",wrap(GetCompletedAlertHistory))
      r.Patch("/alerts/pending/:id", bind(dtos.UpdateAlertActionForm{}), wrap(UpdateAlertAction))
      //maintenance update
      r.Get("/maintenanceAlerts",wrap(GetMaintenanceUpdateForCurrentOrg))
      r.Get("/maintenanceAlerts/get/:interval",wrap(GetMaintenanceAlertsByInterval))
      r.Delete("/maintenanceAlerts/:id",wrap(RemoveMaintenanceUpdateCurrentOrg))
      r.Patch("/maintenanceAlerts/:id",wrap(UpdateMaintenanceCurrentOrg))
      //user action
      r.Post("/maintenanceAlertsUser", quota("maintenance_updated"),bind(m.AddMalfunalertActivity{}), wrap(AddMaintenanceAlertToCurrentOrg))

      r.Get("/maintenanceHistory",wrap(GetMaintenanceHistoryForCurrentOrg))
      r.Get("/maintenanceHistory/get/:interval",wrap(GetMaintenanceHistoryByInterval))
      r.Patch("/maintenanceHistory/:id",wrap(UpdateMaintenanceHistoryCurrentOrg))
      r.Get("/maintenanceActivity",wrap(GetMaintenanceActivitesForCurrentOrg))
      r.Post("/maintenanceActivity", quota("maintenance_activity"),bind(m.AddMaintenanceActivity{}), wrap(AddMaintenanceActivityToCurrentOrg))
    }, reqOrgAdmin)

    // create new org
    r.Post("/orgs", quota("org"), bind(m.CreateOrgCommand{}), wrap(CreateOrg))

    // search all orgs
    r.Get("/orgs", reqGrafanaAdmin, wrap(SearchOrgs))

    // orgs (admin routes)
    r.Group("/orgs/:orgId", func() {
      r.Get("/", wrap(GetOrgById))
      r.Put("/", bind(dtos.UpdateOrgForm{}), wrap(UpdateOrg))
      r.Put("/address", bind(dtos.UpdateOrgAddressForm{}), wrap(UpdateOrgAddress))
      r.Delete("/", wrap(DeleteOrgById))
      r.Get("/users", wrap(GetOrgUsers))
      r.Post("/users", bind(m.AddOrgUserCommand{}), wrap(AddOrgUser))
      r.Patch("/users/:userId", bind(m.UpdateOrgUserCommand{}), wrap(UpdateOrgUser))
      r.Delete("/users/:userId", wrap(RemoveOrgUser))
      r.Get("/quotas", wrap(GetOrgQuotas))
      r.Put("/quotas/:target", bind(m.UpdateOrgQuotaCmd{}), wrap(UpdateOrgQuota))
      //r.Post("/process", bind(m.AddProcessCommand{}), wrap(AddProcess))
      r.Post("/process", quota("process"),bind(m.AddProcessCommand{}), wrap(AddProcessToCurrentOrg))
    }, reqGrafanaAdmin)









    // orgs (admin routes)
    r.Group("/orgs/name/:name", func() {
      r.Get("/", wrap(GetOrgByName))
    }, reqGrafanaAdmin)

    // auth api keys
    r.Group("/auth/keys", func() {
      r.Get("/", wrap(GetApiKeys))
      r.Post("/", quota("api_key"), bind(m.AddApiKeyCommand{}), wrap(AddApiKey))
      r.Delete("/:id", wrap(DeleteApiKey))
    }, reqOrgAdmin)

    // Preferences
    r.Group("/preferences", func() {
      r.Post("/set-home-dash", bind(m.SavePreferencesCommand{}), wrap(SetHomeDashboard))
    })

    // Data sources
    r.Group("/datasources", func() {
      r.Get("/", GetDataSources)
      r.Post("/", quota("data_source"), bind(m.AddDataSourceCommand{}), AddDataSource)
      r.Put("/:id", bind(m.UpdateDataSourceCommand{}), UpdateDataSource)
      r.Delete("/:id", DeleteDataSource)
      r.Get("/:id", wrap(GetDataSourceById))
      r.Get("/name/:name", wrap(GetDataSourceByName))
    }, reqOrgAdmin)

    r.Get("/datasources/id/:name", wrap(GetDataSourceIdByName), reqSignedIn)

    r.Get("/plugins", wrap(GetPluginList))
    r.Get("/plugins/:pluginId/settings", wrap(GetPluginSettingById))

    r.Group("/plugins", func() {
      r.Get("/:pluginId/readme", wrap(GetPluginReadme))
      r.Get("/:pluginId/dashboards/", wrap(GetPluginDashboards))
      r.Post("/:pluginId/settings", bind(m.UpdatePluginSettingCmd{}), wrap(UpdatePluginSetting))
    }, reqOrgAdmin)

    r.Get("/frontend/settings/", GetFrontendSettings)
    r.Any("/datasources/proxy/:id/*", reqSignedIn, ProxyDataSourceRequest)
    r.Any("/datasources/proxy/:id", reqSignedIn, ProxyDataSourceRequest)

    // Dashboard
    r.Group("/dashboards", func() {
      r.Combo("/db/:slug").Get(GetDashboard).Delete(DeleteDashboard)
      r.Post("/db", reqEditorRole, bind(m.SaveDashboardCommand{}), PostDashboard)
      r.Get("/file/:file", GetDashboardFromJsonFile)
      r.Get("/home", wrap(GetHomeDashboard))
      r.Get("/tags", GetDashboardTags)
      r.Post("/import", bind(dtos.ImportDashboardCommand{}), wrap(ImportDashboard))
    })

    // Dashboard snapshots
    r.Group("/dashboard/snapshots", func() {
      r.Get("/", wrap(SearchDashboardSnapshots))
    })

    // Playlist
    r.Group("/playlists", func() {
      r.Get("/", wrap(SearchPlaylists))
      r.Get("/:id", ValidateOrgPlaylist, wrap(GetPlaylist))
      r.Get("/:id/items", ValidateOrgPlaylist, wrap(GetPlaylistItems))
      r.Get("/:id/dashboards", ValidateOrgPlaylist, wrap(GetPlaylistDashboards))
      r.Delete("/:id", reqEditorRole, ValidateOrgPlaylist, wrap(DeletePlaylist))
      r.Put("/:id", reqEditorRole, bind(m.UpdatePlaylistCommand{}), ValidateOrgPlaylist, wrap(UpdatePlaylist))
      r.Post("/", reqEditorRole, bind(m.CreatePlaylistCommand{}), wrap(CreatePlaylist))
    })

    // Search
    r.Get("/search/", Search)

    // metrics
    r.Get("/metrics/test", wrap(GetTestMetrics))

    // metrics
    r.Get("/metrics", wrap(GetInternalMetrics))

    // error test
    r.Get("/metrics/error", wrap(GenerateError))

  }, reqSignedIn)

  // admin api
  r.Group("/api/admin", func() {
    r.Get("/settings", AdminGetSettings)
    r.Post("/users", bind(dtos.AdminCreateUserForm{}), AdminCreateUser)
    r.Put("/users/:id/password", bind(dtos.AdminUpdateUserPasswordForm{}), AdminUpdateUserPassword)
    r.Put("/users/:id/permissions", bind(dtos.AdminUpdateUserPermissionsForm{}), AdminUpdateUserPermissions)
    r.Delete("/users/:id", AdminDeleteUser)
    r.Get("/users/:id/quotas", wrap(GetUserQuotas))
    r.Put("/users/:id/quotas/:target", bind(m.UpdateUserQuotaCmd{}), wrap(UpdateUserQuota))
    r.Get("/stats", AdminGetStats)
  }, reqGrafanaAdmin)

  // rendering
  r.Get("/render/*", reqSignedIn, RenderToPng)

  // grafana.net proxy
  r.Any("/api/gnet/*", reqSignedIn, ProxyGnetRequest)

  // Gravatar service.
  avt := avatar.CacheServer()
  r.Get("/avatar/:hash", avt.ServeHTTP)

  // Websocket
  liveConn := live.New()
  r.Any("/ws", liveConn.Serve)

  // streams
  r.Post("/api/streams/push", reqSignedIn, bind(dtos.StreamMessage{}), liveConn.PushToStream)

  InitAppPluginRoutes(r)

}

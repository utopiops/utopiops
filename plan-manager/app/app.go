package app

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/utopiops-water/framework/middlewares"
	"gitlab.com/utopiops-water/plan-manager/config"
	"gitlab.com/utopiops-water/plan-manager/controllers/healthcheck"
	"gitlab.com/utopiops-water/plan-manager/controllers/paymentController"
	"gitlab.com/utopiops-water/plan-manager/controllers/planController"
	"gitlab.com/utopiops-water/plan-manager/controllers/usageController"
)

type App struct {
	r *gin.Engine
}

func New(healthcheck *healthcheck.HealthCheck, planController *planController.PlanController, usageController *usageController.UsageController, paymentController *paymentController.PaymentController) *App {
	r := gin.Default()
	routing(r, healthcheck, planController, usageController, paymentController)
	return &App{
		r: r,
	}
}

func (a *App) Start(addr ...string) error {
	return a.r.Run(addr...)
}

func routing(r *gin.Engine, healthcheck *healthcheck.HealthCheck, planController *planController.PlanController, usageController *usageController.UsageController, paymentController *paymentController.PaymentController) {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.CORSMiddleware(config.Configs.App.AllowedOrigins))
	r.GET("/health", healthcheck.Healthcheck)
	tempIdsMiddleware := middlewares.NewTempIdsMiddleware(
		config.Configs.Endpoints.IdsPublic,
		config.Configs.Endpoints.IdsAdmin,
		config.Configs.Endpoints.IdsJwks,
		config.Configs.Secrets.AuthServerJwtSecret,
	)
	user := r.Group("/user", tempIdsMiddleware.IdsAuthorizationMiddleware())
	{
		user.POST("/enforce", usageController.Enforce)
	}
	internal := r.Group("", tempIdsMiddleware.IdsAuthorizationMiddleware(), middlewares.Internal())
	{
		//internal.POST("/user/plan/add", usageController.AddPlanToUser)
		internal.POST("/user/usage/buildtime", usageController.UpdateBuildTime)
		internal.POST("/user/usage/bandwidth", func(ctx *gin.Context) { ctx.Status(200) })
		internal.POST("/user/usage/application/create", usageController.CreateApplication)
		internal.POST("/user/usage/application/delete", usageController.DeleteApplication)
		internal.POST("/user/usage/user/create", usageController.AddUser)
		internal.POST("/user/usage/user/delete", usageController.EndUser)
		internal.POST("/user/usage/function/call", usageController.CallFunction)
		internal.POST("/user/access/:resource", usageController.CheckAccess)
		//internal.POST("/user/current/plan", usageController.GetCurrentPlanInternal)
		//internal.POST("/user/plan/change", usageController.ChangeUserPlan)
		internal.POST("/user/init", usageController.InitUser)
	}
}

package todosapi

import (
	"github.com/gin-gonic/gin"

	"vikishptra/shared/gogen"
	"vikishptra/shared/infrastructure/config"
	"vikishptra/shared/infrastructure/logger"
	"vikishptra/shared/infrastructure/token"
)

type selectedRouter = gin.IRouter

type ginController struct {
	*gogen.BaseController
	log      logger.Logger
	cfg      *config.Config
	jwtToken token.JWTToken
}

func NewGinController(log logger.Logger, cfg *config.Config, tk token.JWTToken) gogen.RegisterRouterHandler[selectedRouter] {
	return &ginController{
		BaseController: gogen.NewBaseController(),
		log:            log,
		cfg:            cfg,
		jwtToken:       tk,
	}
}

func (r *ginController) RegisterRouter(router selectedRouter) {

	router.POST("/activity-groups", r.runCreateActivitieHandler())

	router.GET("/activity-groups", r.getAllActivitiesHandler())
	router.GET("/activity-groups/:id", r.findoneactivitiesHandler())
	router.PATCH("/activity-groups/:id", r.runupdateactivitiesHandler())
	router.DELETE("/activity-groups/:id", r.rundeleteactivitieHandler())

	router.POST("/todo-items", r.runcreatetodoHandler())
	router.GET("/todo-items", r.getalltodoHandler())
	router.GET("/todo-items/:id", r.findonetodoHandler())
	router.PATCH("/todo-items/:id", r.runupdatetodoHandler())
	router.DELETE("/todo-items/:id", r.rundeletetodoHandler())
}

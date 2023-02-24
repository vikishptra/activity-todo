package application

import (
	"vikishptra/domain_todos/controller/todosapi"
	"vikishptra/domain_todos/gateway/withgorm"
	"vikishptra/domain_todos/usecase/findoneactivities"
	"vikishptra/domain_todos/usecase/findonetodo"
	"vikishptra/domain_todos/usecase/getallactivities"
	"vikishptra/domain_todos/usecase/getalltodo"
	"vikishptra/domain_todos/usecase/runcreateactivitie"
	"vikishptra/domain_todos/usecase/runcreatetodo"
	"vikishptra/domain_todos/usecase/rundeleteactivitie"
	"vikishptra/domain_todos/usecase/rundeletetodo"
	"vikishptra/domain_todos/usecase/runupdateactivities"
	"vikishptra/domain_todos/usecase/runupdatetodo"
	"vikishptra/shared/gogen"
	"vikishptra/shared/infrastructure/config"
	"vikishptra/shared/infrastructure/logger"
	"vikishptra/shared/infrastructure/server"
	"vikishptra/shared/infrastructure/token"
)

type todos struct{}

func NewTodos() gogen.Runner {
	return &todos{}
}

func (todos) Run() error {

	const appName = "todos"

	cfg := config.ReadConfig()

	appData := gogen.NewApplicationData(appName)

	log := logger.NewSimpleJSONLogger(appData)

	jwtToken := token.NewJWTToken(cfg.JWTSecretKey)

	datasource := withgorm.NewGateway(log, appData, cfg)

	httpHandler := server.NewGinHTTPHandler(log, cfg.Servers[appName].Address, appData)

	x := todosapi.NewGinController(log, cfg, jwtToken)
	x.AddUsecase(
		//
		rundeletetodo.NewUsecase(datasource),
		runupdatetodo.NewUsecase(datasource),
		findonetodo.NewUsecase(datasource),
		getalltodo.NewUsecase(datasource),
		runcreatetodo.NewUsecase(datasource),
		rundeleteactivitie.NewUsecase(datasource),
		runupdateactivities.NewUsecase(datasource),
		findoneactivities.NewUsecase(datasource),
		getallactivities.NewUsecase(datasource),
		runcreateactivitie.NewUsecase(datasource),
	)
	x.RegisterRouter(httpHandler.Router)

	httpHandler.RunWithGracefullyShutdown()

	return nil
}

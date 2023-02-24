package todosapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"vikishptra/domain_todos/usecase/runcreateactivitie"
	"vikishptra/shared/gogen"
	"vikishptra/shared/infrastructure/logger"
	"vikishptra/shared/model/payload"
	"vikishptra/shared/util"
)

func (r *ginController) runCreateActivitieHandler() gin.HandlerFunc {

	type InportRequest = runcreateactivitie.InportRequest
	type InportResponse = runcreateactivitie.InportResponse

	inport := gogen.GetInport[InportRequest, InportResponse](r.GetUsecase(InportRequest{}))

	type request struct {
		InportRequest
	}

	return func(c *gin.Context) {

		traceID := util.GenerateID()

		ctx := logger.SetTraceID(context.Background(), traceID)

		var jsonReq request
		if err := c.BindJSON(&jsonReq); err != nil {
			r.log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, "Bad Request"))
			return
		}

		var req InportRequest
		req.Email = jsonReq.Email
		req.Title = jsonReq.Title
		req.Now = time.Now()

		r.log.Info(ctx, util.MustJSON(req))

		res, err := inport.Execute(ctx, req)
		if err != nil {
			r.log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, "Bad Request"))
			return
		}

		jsonRes := res.Res

		r.log.Info(ctx, util.MustJSON(jsonRes))
		c.JSON(http.StatusCreated, payload.NewSuccessResponse(jsonRes))

	}
}

package todosapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"vikishptra/domain_todos/usecase/getallactivities"
	"vikishptra/shared/gogen"
	"vikishptra/shared/infrastructure/logger"
	"vikishptra/shared/model/payload"
	"vikishptra/shared/util"
)

func (r *ginController) getAllActivitiesHandler() gin.HandlerFunc {

	type InportRequest = getallactivities.InportRequest
	type InportResponse = getallactivities.InportResponse

	inport := gogen.GetInport[InportRequest, InportResponse](r.GetUsecase(InportRequest{}))

	type request struct {
	}

	return func(c *gin.Context) {

		traceID := util.GenerateID()

		ctx := logger.SetTraceID(context.Background(), traceID)

		var jsonReq request
		if err := c.Bind(&jsonReq); err != nil {
			r.log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, "Bad request"))
			return
		}

		var req InportRequest

		r.log.Info(ctx, util.MustJSON(req))

		res, err := inport.Execute(ctx, req)
		if err != nil {
			r.log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, "Bad Request"))
			return
		}

		jsonRes := res.Res

		r.log.Info(ctx, util.MustJSON(jsonRes))
		c.JSON(http.StatusOK, payload.NewSuccessResponse(jsonRes))

	}
}

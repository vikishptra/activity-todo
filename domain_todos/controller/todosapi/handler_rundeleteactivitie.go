package todosapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"vikishptra/domain_todos/model/errorenum"
	"vikishptra/domain_todos/usecase/rundeleteactivitie"
	"vikishptra/shared/gogen"
	"vikishptra/shared/infrastructure/logger"
	"vikishptra/shared/model/payload"
	"vikishptra/shared/util"
)

func (r *ginController) rundeleteactivitieHandler() gin.HandlerFunc {

	type InportRequest = rundeleteactivitie.InportRequest
	type InportResponse = rundeleteactivitie.InportResponse

	inport := gogen.GetInport[InportRequest, InportResponse](r.GetUsecase(InportRequest{}))

	type request struct {
		InportRequest
	}

	return func(c *gin.Context) {

		traceID := util.GenerateID()

		ctx := logger.SetTraceID(context.Background(), traceID)

		var jsonReq request
		if err := c.BindUri(&jsonReq); err != nil {
			r.log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, "Bad Request"))
			return
		}

		var req InportRequest
		req.ID = jsonReq.ID

		r.log.Info(ctx, util.MustJSON(req))

		res, err := inport.Execute(ctx, req)
		if err != nil {
			if err == errorenum.DataNotFound {
				r.log.Error(ctx, err.Error())
				textError := "Activity with ID"
				c.JSON(http.StatusNotFound, payload.NewErrorResponse(fmt.Errorf("%s %d Not Found", textError, req.ID), "Not Found"))
				return
			}
			r.log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, "Bad Request"))
			return
		}

		jsonRes := res

		r.log.Info(ctx, util.MustJSON(jsonRes))
		c.JSON(http.StatusOK, payload.NewSuccessResponse(jsonRes))

	}
}

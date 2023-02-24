package todosapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"vikishptra/domain_todos/model/errorenum"
	"vikishptra/domain_todos/usecase/runupdatetodo"
	"vikishptra/shared/gogen"
	"vikishptra/shared/infrastructure/logger"
	"vikishptra/shared/model/payload"
	"vikishptra/shared/util"
)

func (r *ginController) runupdatetodoHandler() gin.HandlerFunc {

	type InportRequest = runupdatetodo.InportRequest
	type InportResponse = runupdatetodo.InportResponse

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
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, traceID))
			return
		}

		var jsonURI request
		if err := c.ShouldBindUri(&jsonURI); err != nil {

			r.log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, "Bad Request"))
			return
		}

		var req InportRequest
		req.ID = jsonURI.ID
		req.ActivityGroupId = jsonReq.ActivityGroupId
		req.IsActive = jsonReq.IsActive
		req.Now = time.Now()
		req.Title = jsonReq.Title
		req.Priority = jsonReq.Priority
		req.Status = jsonReq.Status

		r.log.Info(ctx, util.MustJSON(req))
		res, err := inport.Execute(ctx, req)
		if err != nil {
			if err == errorenum.DataNotFound {
				r.log.Error(ctx, err.Error())
				textError := "Todo with ID"
				c.JSON(http.StatusNotFound, payload.NewErrorResponse(fmt.Errorf("%s %d Not Found", textError, req.ID), "Not Found"))
				return
			}
			r.log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, "Bad Request"))
			return
		}
		jsonRes := res.Res

		r.log.Info(ctx, util.MustJSON(jsonRes))
		c.JSON(http.StatusOK, payload.NewSuccessResponse(jsonRes))

	}
}

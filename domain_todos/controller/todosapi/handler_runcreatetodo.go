package todosapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"vikishptra/domain_todos/model/errorenum"
	"vikishptra/domain_todos/usecase/runcreatetodo"
	"vikishptra/shared/gogen"
	"vikishptra/shared/infrastructure/logger"
	"vikishptra/shared/model/payload"
	"vikishptra/shared/util"
)

func (r *ginController) runcreatetodoHandler() gin.HandlerFunc {

	type InportRequest = runcreatetodo.InportRequest
	type InportResponse = runcreatetodo.InportResponse

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
		req.Title = jsonReq.Title
		req.ActivityGroupId = jsonReq.ActivityGroupId
		req.Now = time.Now()
		req.Priority = jsonReq.Priority

		r.log.Info(ctx, util.MustJSON(req))

		res, err := inport.Execute(ctx, req)
		if err != nil {
			if err == errorenum.ActivityGroupIdcannotbenull {
				r.log.Error(ctx, err.Error())
				c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, "Bad Request"))
				return
			} else if err == errorenum.DataNotFound {
				r.log.Error(ctx, err.Error())
				textError := "Activity with ID"
				c.JSON(http.StatusNotFound, payload.NewErrorResponse(fmt.Errorf("%s %d Not Found", textError, req.ActivityGroupId), "Not Found"))
				return
			}
			r.log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, payload.NewErrorResponse(err, "Bad Request"))
			return
		}

		jsonRes := res.Res

		r.log.Info(ctx, util.MustJSON(jsonRes))
		c.JSON(http.StatusCreated, payload.NewSuccessResponse(jsonRes))

	}
}

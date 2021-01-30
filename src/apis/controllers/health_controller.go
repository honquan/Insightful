package controllers

import (
	"insightful/src/apis/dtos"
	"net/http"
)

type HealthController struct {
	BaseController
}

func (s *HealthController) HealthCheck(w http.ResponseWriter, req *http.Request) {
	resp := &dtos.HttpResponse{
		Meta: &dtos.MetaResp{
			Code:    http.StatusOK,
			Message: "I'm ok",
		},
	}
	s.ServeJSONWithCode(w, http.StatusOK, resp)
}

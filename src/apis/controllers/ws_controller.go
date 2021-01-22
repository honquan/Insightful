package controllers

import (
	"insightful/src/apis/dtos"
	"net/http"
)

type WsController struct {
	BaseController
}

func (s *WsController) WsGet(w http.ResponseWriter, req *http.Request) {
	//logger.Warnf
	resp := &dtos.HttpResponse{
		Meta: &dtos.MetaResp{
			Code:    http.StatusOK,
			Message: "Ok",
		},
	}
	s.ServeJSONWithCode(w, http.StatusOK, resp)
}

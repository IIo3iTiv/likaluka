package server

import "time"

func (s *Server) SetMainHandlers() {
	s.rgm = s.Router.Group("")
	s.rgm.GET("/ping", MWTimeout(time.Second*60), hPing)
}

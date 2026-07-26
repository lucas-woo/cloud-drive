package gateway



func (s *Server) InitializeRouter() {
	s.gin.GET("/") 
}
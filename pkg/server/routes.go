package server

func (s *Server) routes() {
	s.router.Get("/", s.handleHome())
	s.router.Get("/health", s.handleCheckHealth())
	s.router.Post("/encrypt", s.handleEncrypt())
}

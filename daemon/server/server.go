package server

type Server struct {
	Port    int
	Running bool
}

func NewServer(port int) *Server {
	return &Server{port, false}
}

func (s *Server) Start() error {

}

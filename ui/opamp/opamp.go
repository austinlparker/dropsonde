package opamp

import (
	"fmt"
	"github.com/austinlparker/dropsonde/internal/opampsrv"
	data "github.com/austinlparker/dropsonde/internal/opampsrv/agent"
	"log"
	"os"
	"strings"
)

type Server struct {
	server *opampsrv.Server
}

func NewServer(s *Server) {
	file, err := os.OpenFile("opamp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	var logger = log.New(file, "", log.LstdFlags|log.Lmicroseconds)
	server := opampsrv.NewServer(&data.AllAgents)
	server.Start()
	s.server = server
	logger.Println("Server started")
}

func (s *Server) GetAgents() string {
	d := s.server.GetAllAgents()
	str := strings.Builder{}
	for id, agent := range d {
		str.WriteString(fmt.Sprintf("Instance ID: %v\n", string(id)))
		str.WriteString(fmt.Sprintf("Effective Config: %v\n", agent.EffectiveConfig))
		str.WriteString(agent.Status.String())
	}
	return str.String()
}

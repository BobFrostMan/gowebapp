package server

type Server struct {
	Hostname string
	Port int
	AllowedHosts []string //list of hosts to allow requests from
}

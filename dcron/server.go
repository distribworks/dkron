package dcron

import (
	"encoding/json"
	"flag"
	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	"github.com/mitchellh/cli"
	"github.com/tylerb/graceful"
	"net/http"
	"strings"
	"time"
)

// ServerCommand run dcron server
type ServerCommand struct {
	Ui cli.Ui
}

func (s *ServerCommand) Help() string {
	helpText := `
Usage: dcron server [options]
	Provides debugging information for operators
Options:
  -format                  If provided, output is returned in the specified
                           format. Valid formats are 'json', and 'text' (default)
  -rpc-addr=127.0.0.1:7373 RPC address of the Serf agent.
  -rpc-auth=""             RPC auth token of the Serf agent.
`
	return strings.TrimSpace(helpText)
}

func (s *ServerCommand) Run(args []string) int {
	var format string
	cmdFlags := flag.NewFlagSet("server", flag.ContinueOnError)
	cmdFlags.Usage = func() { s.Ui.Output(s.Help()) }
	cmdFlags.StringVar(&format, "format", "text", "output format")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	return 0
}

func (s *ServerCommand) Synopsis() string {
	return "Run dcron server"
}

func ServerInit() {
	loadConfig()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)

	middle := interpose.New()
	middle.UseHandler(router)

	srv := &graceful.Server{
		Timeout: 1 * time.Second,
		Server:  &http.Server{Addr: ":8080", Handler: middle},
	}

	log.Infoln("Running HTTP server on 8080")

	certFile := config.GetString("certFile")
	keyFile := config.GetString("keyFile")
	if certFile != "" && keyFile != "" {
		srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		srv.ListenAndServe()
	}
	log.Debug("Exiting")
}

func Index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(jobs)
}

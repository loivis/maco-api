package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/TV4/env"
	"github.com/loivis/maco-api/api"
	"github.com/rs/zerolog/log"
)

func main() {
	conf := readConfig()
	fmt.Fprintln(os.Stderr, conf)

	mux := http.NewServeMux()

	api.New().Register(mux)

	addr := ":" + conf.port
	log.Info().Str("port", conf.port).Msg("starting server ...")

	if err := http.ListenAndServe(addr, mux); err != http.ErrServerClosed {
		log.Fatal().Msgf("failed to start server: %v", err)
	}
}

type config struct {
	port string
}

func readConfig() *config {
	return &config{
		port: env.String("PORT", "8080"),
	}
}

func (c *config) String() string {
	// hideIfSet := func(v interface{}) string {
	// 	s := ""

	// 	switch typedV := v.(type) {
	// 	case string:
	// 		s = typedV
	// 	case []string:
	// 		s = strings.Join(typedV, ",")
	// 	case []byte:
	// 		s = string(typedV)
	// 	case fmt.Stringer:
	// 		if typedV != nil {
	// 			s = typedV.String()
	// 		}
	// 	}

	// 	if s != "" {
	// 		return "<hidden>"
	// 	}
	// 	return ""
	// }

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 1, 4, ' ', 0)
	for _, e := range []struct {
		k string
		v interface{}
	}{
		{"PORT", c.port},
	} {
		fmt.Fprintf(w, "%s\t%v\n", e.k, e.v)
	}
	w.Flush()
	return buf.String()
}

package chi

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jwalton/gchalk"
	"github.com/vixyninja/go-blocks/logx"
	"golang.org/x/mod/modfile"
)

type Server struct {
	router            *chi.Mux
	httpServer        *http.Server
	log               logx.Logx
	addr              string
	stopTimeout       time.Duration
	moduleName        string
	printRoutes       bool
	readHeaderTimeout time.Duration
	readTimeout       time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
	requestTimeout    time.Duration
	enableRequestID   bool
	enableRealIP      bool
	enableRecoverer   bool
}

type Option func(*Server)

func WithAddr(addr string) Option {
	return func(s *Server) { s.addr = addr }
}

func WithLogger(l logx.Logx) Option {
	return func(s *Server) { s.log = l }
}

func WithStopTimeout(d time.Duration) Option {
	return func(s *Server) { s.stopTimeout = d }
}

func WithPrintRoutes(enabled bool) Option {
	return func(s *Server) { s.printRoutes = enabled }
}

func WithReadHeaderTimeout(d time.Duration) Option {
	return func(s *Server) { s.readHeaderTimeout = d }
}

func WithReadTimeout(d time.Duration) Option {
	return func(s *Server) { s.readTimeout = d }
}

func WithWriteTimeout(d time.Duration) Option {
	return func(s *Server) { s.writeTimeout = d }
}

func WithIdleTimeout(d time.Duration) Option {
	return func(s *Server) { s.idleTimeout = d }
}

func WithRequestTimeout(d time.Duration) Option {
	return func(s *Server) { s.requestTimeout = d }
}

func WithRequestID(enabled bool) Option {
	return func(s *Server) { s.enableRequestID = enabled }
}

func WithRealIP(enabled bool) Option {
	return func(s *Server) { s.enableRealIP = enabled }
}

func WithRecoverer(enabled bool) Option {
	return func(s *Server) { s.enableRecoverer = enabled }
}
func NewServer(opts ...Option) *Server {
	s := &Server{
		router:            chi.NewRouter(),
		log:               logx.NewStdLogger(),
		addr:              ":4433",
		stopTimeout:       10 * time.Second,
		moduleName:        getModuleName(),
		readHeaderTimeout: 30 * time.Second,
		readTimeout:       60 * time.Second,
		writeTimeout:      60 * time.Second,
		idleTimeout:       90 * time.Second,
		requestTimeout:    60 * time.Second,
		enableRequestID:   true,
		enableRealIP:      true,
		enableRecoverer:   true,
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.enableRequestID {
		s.router.Use(middleware.RequestID)
	}
	if s.enableRealIP {
		s.router.Use(middleware.RealIP)
	}
	if s.enableRecoverer {
		s.router.Use(middleware.Recoverer)
	}
	if s.requestTimeout > 0 {
		s.router.Use(middleware.Timeout(s.requestTimeout))
	}

	return s
}

func (s *Server) Router() *chi.Mux { return s.router }

func (s *Server) HTTPServer() *http.Server {
	if s.httpServer == nil {
		s.httpServer = &http.Server{
			Addr:              s.addr,
			Handler:           s.router,
			ReadHeaderTimeout: s.readHeaderTimeout,
			ReadTimeout:       s.readTimeout,
			WriteTimeout:      s.writeTimeout,
			IdleTimeout:       s.idleTimeout,
		}
	}
	return s.httpServer
}

func (s *Server) Run() error {
	if s.printRoutes {
		s.PrintRoutes()
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) PrintRoutes(excludePatterns ...string) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	err := chi.Walk(s.router, func(method, route string, handler http.Handler, _ ...func(http.Handler) http.Handler) error {
		for _, pattern := range excludePatterns {
			if strings.HasPrefix(route, pattern) {
				return nil
			}
		}

		color := methodColor(method)
		handlerName := formatHandlerName(s.moduleName, handler)

		fmt.Fprintf(tw, "%s\t%-30s\t%s\n",
			color(fmt.Sprintf("%-7s", method)),
			route,
			handlerName,
		)
		return nil
	})

	_ = tw.Flush()

	if err != nil {
		s.log.Info(context.TODO(), fmt.Sprintf("route walker error: %v", err))
	}
}

func methodColor(method string) func(a ...string) string {
	switch method {
	case http.MethodGet:
		return gchalk.Green
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return gchalk.Yellow
	case http.MethodDelete:
		return gchalk.Red
	default:
		return gchalk.White
	}
}

func formatHandlerName(module string, handler http.Handler) string {
	v := reflect.ValueOf(handler)
	t := v.Type()

	var name string

	switch t.Kind() {
	case reflect.Func:
		name = runtime.FuncForPC(v.Pointer()).Name()
	default:
		if m, ok := t.MethodByName("ServeHTTP"); ok {
			name = runtime.FuncForPC(m.Func.Pointer()).Name()
		} else {
			name = t.String()
		}
	}

	name = filepath.Base(name)
	if module != "" {
		return gchalk.Blue(fmt.Sprintf("%s/%s", module, name))
	}
	return gchalk.Blue(name)
}

func getModuleName() string {
	b, err := os.ReadFile("go.mod")
	if err != nil {
		return ""
	}
	return modfile.ModulePath(b)
}

package chi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jwalton/gchalk"
	"github.com/vixyninja/go-blocks/logx"
	"golang.org/x/mod/modfile"
)

type Server struct {
	router      *chi.Mux
	httpServer  *http.Server
	log         logx.Logx
	addr        string
	stopTimeout time.Duration
	moduleName  string
	printRoutes bool
}

type Option func(*Server)

func WithAddr(addr string) Option            { return func(s *Server) { s.addr = addr } }
func WithLogger(l logx.Logx) Option          { return func(s *Server) { s.log = l } }
func WithStopTimeout(d time.Duration) Option { return func(s *Server) { s.stopTimeout = d } }
func WithPrintRoutes(enabled bool) Option {
	return func(s *Server) { s.printRoutes = enabled }
}
func NewServer(opts ...Option) *Server {
	s := &Server{
		router:      chi.NewRouter(),
		log:         logx.NewStdLogger(),
		addr:        ":4433",
		stopTimeout: 10 * time.Second,
		moduleName:  getModuleName(),
	}

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Router() *chi.Mux { return s.router }

func (s *Server) Run() error {
	return s.RunContext(context.Background())
}

func (s *Server) RunContext(ctx context.Context) error {
	if s.printRoutes {
		s.PrintRoutes()
	}
	s.httpServer = &http.Server{
		Addr:              s.addr,
		Handler:           s.router,
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       90 * time.Second,
	}

	// TODO: using hook
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.stopTimeout)
		defer cancel()
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			s.log.Info(context.TODO(), fmt.Sprintf("graceful shutdown failed: %v", err))
		}
	}()

	s.log.Info(context.TODO(), fmt.Sprintf("HTTP server listening on %s", s.addr))

	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Info(context.TODO(), fmt.Sprintf("listen and serve error: %v", err))
		return err
	}

	s.log.Info(context.TODO(), "server stopped")
	return nil
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

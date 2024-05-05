package main

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "officetime-api/api/v1/pb/api"
	"officetime-api/internal/config"
	"officetime-api/internal/db"
	departmentAdapter "officetime-api/internal/model/department/adapter"
	departmentApp "officetime-api/internal/model/department/app"
	departmentCommand "officetime-api/internal/model/department/app/command"
	departmentCommandQuery "officetime-api/internal/model/department/app/command_query"
	departmentQuery "officetime-api/internal/model/department/app/query"
	periodAdapter "officetime-api/internal/model/period/adapter"
	periodApp "officetime-api/internal/model/period/app"
	periodCommand "officetime-api/internal/model/period/app/command"
	periodCommandQuery "officetime-api/internal/model/period/app/command_query"
	periodQuery "officetime-api/internal/model/period/app/query"
	routerAdapter "officetime-api/internal/model/router/adapter"
	routerApp "officetime-api/internal/model/router/app"
	routerCommand "officetime-api/internal/model/router/app/command"
	routerCommandQuery "officetime-api/internal/model/router/app/command_query"
	routerQuery "officetime-api/internal/model/router/app/query"
	userAdapter "officetime-api/internal/model/user/adapter"
	userApp "officetime-api/internal/model/user/app"
	userCommand "officetime-api/internal/model/user/app/command"
	userCommandQuery "officetime-api/internal/model/user/app/command_query"
	userQuery "officetime-api/internal/model/user/app/query"
	weekendAdapter "officetime-api/internal/model/weekend/adapter"
	weekendApp "officetime-api/internal/model/weekend/app"
	weekendQuery "officetime-api/internal/model/weekend/app/query"
	"officetime-api/pkg/api"
	"officetime-api/pkg/api/endpoints"
	"officetime-api/pkg/api/transport"
	"os"
	"os/signal"
	"syscall"

	//"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/oklog/oklog/pkg/group"
	"net/http"
	"officetime-api/app/service"
	"strconv"
)

func main() {
	cfg := config.New()

	//dao.Db = db.GetDB()
	//model.Db = db.GetDB()
	//aggregator.Db = db.GetDB()

	/*fmt.Println("Setting up server, enabling CORS...")
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                    // All origins
		AllowedMethods: []string{"GET", "POST", "PATCH"}, // Allowing only get, just an example
	})*/

	var (
		//logger   log.Logger
		httpAddr = net.JoinHostPort("", cfg.HttpPort)
		grpcAddr = net.JoinHostPort("", cfg.GrpcPort)
	)

	//logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	//logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	routerRepository := routerAdapter.NewPgRouterRepository(db.GetDB())
	rApp := routerApp.Application{
		Commands: routerApp.Commands{
			UpdateRouter: routerCommand.NewUpdateRouterHandler(routerRepository),
			DeleteRouter: routerCommand.NewDeleteRouterHandler(routerRepository),
		},
		Queries: routerApp.Queries{
			GetRouter:  routerQuery.NewGetRouterHandler(routerRepository),
			GetRouters: routerQuery.NewGetRoutersHandler(routerRepository),
		},
		CommandsQueries: routerApp.CommandsQueries{
			CreateRouter: routerCommandQuery.NewCreateRouterHandler(routerRepository),
		},
	}

	periodRepository := periodAdapter.NewPgPeriodRepository(db.GetDB())
	pApp := periodApp.Application{
		Commands: periodApp.Commands{
			UpdatePeriod: periodCommand.NewUpdatePeriodHandler(periodRepository),
			DeletePeriod: periodCommand.NewDeletePeriodHandler(periodRepository),
		},
		Queries: periodApp.Queries{
			GetPeriod:  periodQuery.NewGetPeriodHandler(periodRepository),
			GetPeriods: periodQuery.NewGetPeriodsHandler(periodRepository),
		},
		CommandsQueries: periodApp.CommandsQueries{
			CreatePeriod: periodCommandQuery.NewCreatePeriodHandler(periodRepository),
		},
	}

	departmentRepository := departmentAdapter.NewPgDepartmentRepository(db.GetDB())
	dApp := departmentApp.Application{
		Commands: departmentApp.Commands{
			UpdateDepartment: departmentCommand.NewUpdateDepartmentHandler(departmentRepository),
			DeleteDepartment: departmentCommand.NewDeleteDepartmentHandler(departmentRepository),
		},
		Queries: departmentApp.Queries{
			GetDepartment:  departmentQuery.NewGetDepartmentHandler(departmentRepository),
			GetDepartments: departmentQuery.NewGetDepartmentsHandler(departmentRepository),
		},
		CommandsQueries: departmentApp.CommandsQueries{
			CreateDepartment: departmentCommandQuery.NewCreateDepartmentHandler(departmentRepository),
		},
	}

	weekendRepository := weekendAdapter.NewPgWeekendRepository(db.GetDB())
	wApp := weekendApp.Application{
		Queries: weekendApp.Queries{
			GetWeekends: weekendQuery.NewGetWeekendsHandler(weekendRepository),
		},
	}

	userRepository := userAdapter.NewPgUserRepository(db.GetDB())
	uApp := userApp.Application{
		Commands: userApp.Commands{
			UpdateUser: userCommand.NewUpdateUserHandler(userRepository),
			DeleteUser: userCommand.NewDeleteUserHandler(userRepository),
		},
		Queries: userApp.Queries{
			GetUser:                    userQuery.NewGetUserHandler(userRepository),
			GetUsers:                   userQuery.NewGetUsersHandler(userRepository),
			GetUserByEmail:             userQuery.NewGetUserByEmailHandler(userRepository),
			GetUsersByDepartment:       userQuery.NewGetUsersByDepartmentHandler(userRepository),
			GetUserPasswordHashByEmail: userQuery.NewGetUserPasswordHashByEmailHandler(userRepository),
		},
		CommandsQueries: userApp.CommandsQueries{
			CreateUser: userCommandQuery.NewCreateUserHandler(userRepository),
		},
	}

	var (
		s           = api.NewService(rApp, pApp, dApp, wApp, uApp)
		eps         = endpoints.NewEndpointSet(s)
		httpHandler = transport.NewHTTPHandler(eps)
		grpcServer  = transport.NewGRPCServer(eps)
	)

	// API Gateway
	var g group.Group
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			log.Fatal(err)
		}
		g.Add(func() error {
			log.Printf("Serving http address %s", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(err error) {
			httpListener.Close()
		})
	}
	{
		// The gRPC listener mounts the Go kit gRPC server we created.
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatal(err)
		}
		g.Add(func() error {
			log.Printf("Serving grpc address %s", grpcAddr)
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
			pb.RegisterApiServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	if err := g.Run(); err != nil {
		log.Fatal(err)
	}

	/*router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)

	router.HandleFunc("/api-service/token/refresh", dao.Refresh).Methods("POST")
	router.HandleFunc("/api-service/logout", dao.Logout).Methods("POST")

	router.Handle("/metrics", promhttp.Handler())

	router.Handle("/api-service/time/{id}/day/{date}", isAuthorized(dao.GetTimeDayAll)).Methods("GET")
	router.Handle("/api-service/time/{id}/period/{period}", isAuthorized(dao.GetTimeByPeriod)).Methods("GET")
	router.Handle("/api-service/time", isAuthorized(dao.CreateTime)).Methods("POST")

	router.Handle("/api-service/stat/periods-and-routers", isAuthorized(dao.GetStatByPeriodsAndRouters)).Methods("GET")
	router.Handle("/api-service/stat/departments/{date}", isAuthorized(dao.GetAllTimesDepartmentsByDate)).Methods("GET")
	router.Handle("/api-service/stat/working-period/{id}/period/{period}", isAuthorized(dao.GetStatWorkingPeriod)).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+cfg.HttpPort, c.Handler(router)))*/
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	cfg := config.New()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return []byte(cfg.AuthSigningKey), nil
			})

			if err != nil {
				if "Token is expired" == err.Error() {
					w.WriteHeader(http.StatusUnauthorized)
				}
				fmt.Fprintf(w, err.Error())
				return
			}

			au, err := service.ExtractTokenMetadata(r)
			if err != nil {
				if "Token is expired" == err.Error() {
					w.WriteHeader(http.StatusUnauthorized)
				}
				fmt.Fprintf(w, err.Error())
				return
			}
			if au == nil {
				fmt.Fprintf(w, err.Error())
				if "Token is expired" == err.Error() {
					w.WriteHeader(http.StatusUnauthorized)
				}
				return
			}

			id := mux.Vars(r)["id"]
			userId, _ := strconv.Atoi(id) // TODO: если обрабатывать ошибку, апи падает

			if au.UserId != uint64(userId) && !au.IsAdmin() {
				fmt.Fprintf(w, "Access denied")
				return
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

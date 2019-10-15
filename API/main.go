package main

import (
	"flag"
	"fmt"
	"github.com/samueldaviddelacruz/go-job-board/API/middleware"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samueldaviddelacruz/go-job-board/API/controllers"
	"github.com/samueldaviddelacruz/go-job-board/API/email"

	"github.com/samueldaviddelacruz/go-job-board/API/models"
)

func main() {

	boolPtr := flag.Bool("prod", false,
		"Provide this flag in production. This ensures that a config.json file is provided before the application starts")
	flag.Parse()
	appCfg := LoadConfig(*boolPtr)
	databaseConfig := appCfg.Database

	services, err := models.NewServices(
		models.WithGorm(
			databaseConfig.Dialect(),
			databaseConfig.ConnectionInfo()),
		models.WithLogMode(!appCfg.IsProd()),
		models.WithUser(appCfg.Pepper, appCfg.HMACKey),
		models.WithJobPost(),
		models.WithSkill(),
		models.WithOAuth(),
	)
	must(err)

	defer services.Close()
	must(services.DestructiveReset())
	//must(services.AutoMigrate())

	mgCfg := appCfg.Mailgun
	emailer := email.NewClient(
		email.WithSender("lenslocked-project-demo.net Support", "support@sandboxddba781be75b455ea3313563bb0b74b2.mailgun.org"),
		email.WithMailgun(mgCfg.Domain, mgCfg.APIKey),
	)

	r := mux.NewRouter()

	jobsC := controllers.NewJobs(services.JobPost, services.Skill)

	usersC := controllers.NewUsers(services.User, services.Skill)
	authC := controllers.NewAuth(services.User, emailer)
	privateKey, err := ioutil.ReadFile("./key.priv")
	must(err)

	requireJWT := middleware.RequireJWT{
		Secret: string(privateKey),
	}

	applyRoutes(r,
		Route{
			path:    "/signup",
			handler: authC.Create,
			method:  "POST",
		},
		Route{
			path:    "/login",
			handler: authC.Login,
			method:  "POST",
		},
		Route{
			path:    "/user/{id:[0-9]+}",
			handler: usersC.Update,
			method:  "PUT",
		},
		Route{
			path:    "/user/{id:[0-9]+}/company-profile",
			handler: requireJWT.ApplyFn(usersC.UpdateCompanyProfile),
			method:  "PUT",
		},
		Route{
			path:    "/user/{id:[0-9]+}/company-profile/add-skill",
			handler: usersC.AddCompanyProfileSkill,
			method:  "PUT",
		},
		Route{
			path:    "/user/{id:[0-9]+}/company-profile/remove-skill",
			handler: usersC.RemoveCompanyProfileSkill,
			method:  "PUT",
		},
		Route{
			path:    "/user/{id:[0-9]+}/company-profile/add-benefit",
			handler: usersC.AddCompanyProfileBenefit,
			method:  "PUT",
		},
		Route{
			path:    "/user/{id:[0-9]+}/company-profile/remove-benefit",
			handler: usersC.RemoveCompanyProfileBenefit,
			method:  "PUT",
		},
		Route{
			path:    "/user/{id:[0-9]+}/company-profile/update-benefit",
			handler: usersC.UpdateCompanyProfileBenefit,
			method:  "PUT",
		},
		Route{
			path:    "/jobs",
			handler: jobsC.List,
			method:  "GET",
		},
		Route{
			path:    "/jobs",
			handler: jobsC.Create,
			method:  "POST",
		},
		Route{
			path:    "/jobs/{id:[0-9]+}",
			handler: jobsC.Update,
			method:  "PUT",
		},
		Route{
			path:    "/jobs/{id:[0-9]+}",
			handler: jobsC.Delete,
			method:  "DELETE",
		},
		Route{
			path:    "/jobs/{id:[0-9]+}/add-skill",
			handler: jobsC.AddJobPostSkill,
			method:  "PUT",
		},
		Route{
			path:    "/jobs/{id:[0-9]+}/remove-skill",
			handler: jobsC.RemoveJobPostSkill,
			method:  "PUT",
		},
	)

	fmt.Printf("Running on port :%d", appCfg.Port)
	must(http.ListenAndServe(fmt.Sprintf(":%d", appCfg.Port), r))
}

type Route struct {
	path    string
	handler func(http.ResponseWriter, *http.Request)
	method  string
}

func applyRoutes(r *mux.Router, routes ...Route) {
	for _, route := range routes {
		r.HandleFunc(route.path, route.handler).Methods(route.method)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

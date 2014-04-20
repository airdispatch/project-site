package main

import (
	"flag"
	"fmt"
	"github.com/airdispatch/go-pressure"
	"os"
	"path/filepath"
)

// var WORKING_DIRECTORY string
// var TEMPLATE_DIRECTORY string
var PORT string

var flag_port = flag.String("port", "", "specify the port that the server should run on")
var blog_address = flag.String("blog_address", "", "specify the address to load the blog from")
var debug = flag.Bool("debug", false, "specify whether you want to debug the program or not")

func main() {
	flag.Parse()

	temp_port := ":" + os.Getenv("PORT")
	if temp_port == "" {
		temp_port = *flag_port
	}

	temp_wd, _ := os.Getwd()

	// Get Relevant Paths
	template_dir := filepath.Join(temp_wd, "templates")
	static_dir := filepath.Join(temp_wd, "static")

	// Create Server and Necessary Engines
	theServer := pressure.CreateServer(temp_port, *debug)
	tEng := theServer.CreateTemplateEngine(template_dir, "base.html")

	// Register Golang Import URLs
	theServer.RegisterURL(
		pressure.NewURLRoute("^/crypto$", &GolangFetchController{"airdispatch"}),
		pressure.NewURLRoute("^/errors$", &GolangFetchController{"airdispatch"}),
		pressure.NewURLRoute("^/identity$", &GolangFetchController{"airdispatch"}),
		pressure.NewURLRoute("^/message$", &GolangFetchController{"airdispatch"}),
		pressure.NewURLRoute("^/routing$", &GolangFetchController{"airdispatch"}),
		pressure.NewURLRoute("^/server$", &GolangFetchController{"airdispatch"}),
		pressure.NewURLRoute("^/server/server$", &GolangFetchController{"airdispatch"}),
		pressure.NewURLRoute("^/wire$", &GolangFetchController{"airdispatch"}),

		// Tracker Location
		pressure.NewURLRoute("^/tracker$", &GolangFetchController{"tracker"}),
		pressure.NewURLRoute("^/tracker/tracker$", &GolangFetchController{"tracker"}),
		pressure.NewURLRoute("^/tracker/wire$", &GolangFetchController{"tracker"}),
	)

	// Register URLs
	theServer.RegisterURL(
		pressure.NewURLRoute("^/project/airdispatch", &ProjectController{"protocol", tEng}),
		pressure.NewURLRoute("^/$", &HomepageController{tEng}),
		pressure.NewStaticFileRoute("^/static/", static_dir),
	)

	// Start the Server
	theServer.RunServer()
}

// Define Custom Controllers

type ProjectController struct {
	Project string
	tEng    *pressure.TemplateEngine
}

func (c *ProjectController) GetResponse(p *pressure.Request, l *pressure.Logger) (pressure.View, *pressure.HTTPError) {
	return c.tEng.NewTemplateView(fmt.Sprintf("projects/%v.html", c.Project), nil), nil
}

type GolangFetchController struct {
	packageName string
}

func (c *GolangFetchController) GetResponse(p *pressure.Request, l *pressure.Logger) (pressure.View, *pressure.HTTPError) {
	return pressure.NewHTMLView(
		`<html>
			<head>
				<meta name="go-import" content="airdispat.ch git https://github.com/airdispatch/` + c.packageName + `">
			</head>
		</html>`), nil
}

type HomepageController struct {
	tEng *pressure.TemplateEngine
}

func (c *HomepageController) GetResponse(p *pressure.Request, l *pressure.Logger) (pressure.View, *pressure.HTTPError) {
	return c.tEng.NewTemplateView("index.html", nil), nil
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"ghwebhook/ghwebhook"

	"github.com/go-playground/webhooks/v6/github"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func main() {
	/**
	### Initialize GitHub webhook parser ###
	-   This parser helps takes care about parsing the plain HTTP(S)
	    request to the GitHub event struct the application needs to work with.
	-   Also, it (optionally) can check for a secret that can
	    be defined for webhooks.
	-   In this case the parser will also verify that the incoming request contains the secret "my_secret"
	    otherwise it will error.
	*/
	hook, _ := github.New(github.Options.Secret("my_secret"))

	/**
	Create new Go Fiber application
	*/
	app := fiber.New()

	config, err := ghwebhook.GetConfig("config.yml")
	if err != nil {
		log.Fatal("Could not read config")
	}
	/**
	Create POST endpoint on the "/webhook" route to handle incoming webhooks
	*/
	app.Post("/webhook", func(c *fiber.Ctx) error {
		log.Println("Received webhook...")
		httpRequest := new(http.Request)

		err := fasthttpadaptor.ConvertRequest(c.Context(), httpRequest, true)
		if err != nil {
			log.Println("Error converting request", err)
		}
		payload, e := hook.Parse(httpRequest, github.ReleaseEvent)
		if e != nil {
			log.Println("Error parsing", e)
		}
		// fmt.Printf("Payload: %v \n", payload)
		switch payload.(type) {
		case github.CreatePayload:
			{
				createPayload := payload.(github.CreatePayload)
				fmt.Printf("%+v\n", createPayload)

			}
		case github.ReleasePayload:
			{
				releasePayload := payload.(github.ReleasePayload)
				println(releasePayload.Action)
				if releasePayload.Action == "published" {
					ghwebhook.ProcessNewRelease(config, releasePayload)
				}

			}
		case github.PushPayload:
			{
				pushPayload := payload.(github.PushPayload)
				fmt.Printf("%+v\n", pushPayload)

			}
		}
		return c.SendStatus(200)
	})

	app.Listen(":4567")
}

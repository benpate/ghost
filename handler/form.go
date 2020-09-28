package handler

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/ghost/service"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
)

// GetForm generates an HTML form for the requested stream/transition
func GetForm(factoryManager *service.FactoryManager) echo.HandlerFunc {

	return func(ctx echo.Context) error {

		// Get factory for this context
		factory, err := factoryManager.ByContext(ctx)

		if err != nil {
			return err
		}

		// Try to load required values
		streamService := factory.Stream()
		token := ctx.Param("token")
		transitionID := ctx.Param("transitionId")
		stream, err := streamService.LoadByToken(token)

		if err != nil {
			return derp.Report(derp.Wrap(err, "ghost.handler.GetTransition", "Cannot load Stream", token))
		}

		// Render the HTML
		result, err := factory.FormRenderer(ctx, stream, transitionID).Render()

		if err != nil {
			return derp.Report(derp.Wrap(err, "ghost.handler.GetTransition", "Error rendering form"))
		}

		// Success!
		return ctx.HTML(200, result)
	}
}

// PostForm returns an echo.HandlerFunc that accepts form posts
// and performs actions on streams based on the user's permissions.
func PostForm(factoryManager *service.FactoryManager) echo.HandlerFunc {

	return func(ctx echo.Context) error {

		// Get Factory and services required for this step
		factory, err := factoryManager.ByContext(ctx)

		if err != nil {
			return err
		}
		// Get parameters from context
		token := ctx.Param("token")
		transitionID := ctx.Param("transitionId")

		form := make(map[string]interface{})

		if err := ctx.Bind(&form); err != nil {
			return derp.Report(derp.Wrap(err, "ghost.handler.PostTransition", "Cannot load parse form data"))
		}

		spew.Dump(token, transitionID, form)

		streamService := factory.Stream()
		templateService := factory.Template()

		nextView := "default"

		// Load stream
		stream, err := streamService.LoadByToken(token)

		if err != nil {
			return derp.Report(derp.Wrap(err, "ghost.handler.PostTransition", "Cannot load stream", token))
		}

		// Load template
		template, err := templateService.Load(stream.Template)

		if err != nil {
			return derp.Report(derp.Wrap(err, "ghost.handler.PostTransition", "Cannot load template", stream))
		}

		// Execute transition
		if transition, err := template.Transition(stream.State, transitionID); err == nil {

			if err := streamService.Transition(stream, template, transitionID, form); err != nil {
				return derp.Report(derp.Wrap(err, "ghost.handler.PostTransition", "Error updating stream"))
			}

			nextView = transition.NextView
		}

		/// RENDER THE STREAM HERE
		ctx.QueryParams().Del("view")
		ctx.QueryParams().Add("view", nextView)
		result, err := factory.StreamRenderer(ctx, stream).Render()

		if err != nil {
			return derp.Report(derp.Wrap(err, "ghost.handler.GetStream", "Error rendering innerHTML"))
		}

		return ctx.HTML(http.StatusOK, result)
	}
}

package controllers

import (
	"net/http"

	"github.com/Francesco99975/reviews/internal/database"
	"github.com/Francesco99975/reviews/internal/helpers"
	"github.com/Francesco99975/reviews/internal/models"
	"github.com/Francesco99975/reviews/internal/repository"
	"github.com/Francesco99975/reviews/views"
	"github.com/Francesco99975/reviews/views/layouts"
	"github.com/labstack/echo/v4"
)

func Index() echo.HandlerFunc {
	return func(c echo.Context) error {
		data := models.GetDefaultSite("Home")

		data.CSRF = c.Get("csrf").(string)
		data.Nonce = c.Get("nonce").(string)

		// Configure tabs (or use DefaultTabs())
		tabs := layouts.DefaultTabs()

		// Determine active tab from route
		activeTab := "invoices" // Parse from request path

		props := layouts.TabLayoutProps{
			Site:      data,
			Tabs:      tabs,
			ActiveTab: activeTab,
		}

		ctx := c.Request().Context()

		tx, err := database.Pool().Begin(ctx)
		if err != nil {
			return helpers.SendReturnedGenericHTMLError(c, helpers.GenericError{Code: http.StatusInternalServerError, Message: err.Error(), UserMessage: "Server is not accessible to find this resource"}, nil)
		}
		defer database.HandleTransaction(ctx, tx, &err)

		repo := repository.New(tx)

		invoices, err := repo.GetAllInvoicesWithReview(ctx)
		if err != nil {
			return helpers.SendReturnedGenericHTMLError(c, helpers.GenericError{Code: http.StatusInternalServerError, Message: err.Error(), UserMessage: "Server is not accessible to find this resource"}, nil)
		}

		html := helpers.MustRenderHTML(views.Index(props, invoices))

		return c.Blob(http.StatusOK, "text/html; charset=utf-8", html)
	}
}

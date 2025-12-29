package controllers

import (
	"net/http"

	"github.com/Francesco99975/reviews/internal/database"
	"github.com/Francesco99975/reviews/internal/enums"
	"github.com/Francesco99975/reviews/internal/helpers"
	"github.com/Francesco99975/reviews/internal/models"
	"github.com/Francesco99975/reviews/internal/repository"
	"github.com/Francesco99975/reviews/views"
	"github.com/Francesco99975/reviews/views/components"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ReviewsTab() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		tx, err := database.Pool().Begin(ctx)
		if err != nil {
			return helpers.SendReturnedGenericHTMLError(c, helpers.GenericError{Code: http.StatusInternalServerError, Message: err.Error(), UserMessage: "Server is not accessible to find this resource"}, nil)
		}
		defer database.HandleTransaction(ctx, tx, &err)

		repo := repository.New(tx)

		reviews, err := repo.GetAllReviews(ctx)
		if err != nil {
			return helpers.SendReturnedGenericHTMLError(c, helpers.GenericError{Code: http.StatusInternalServerError, Message: err.Error(), UserMessage: "Server is not accessible to find this resource"}, nil)
		}

		csrf := c.Get("csrf").(string)

		html := helpers.MustRenderHTML(views.ReviewsTab(reviews, csrf))

		return c.Blob(http.StatusOK, "text/html; charset=utf-8", html)
	}
}

func SendReview() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload models.ReviewDTO

		err := c.Bind(&payload)
		if err != nil {
			return helpers.SendReturnedHTMLErrorMessage(c, helpers.ErrorMessage{Error: helpers.GenericError{Code: http.StatusBadRequest, Message: err.Error(), UserMessage: "Invalid data"}, Box: enums.Boxes.TOAST_TR, Persistance: "3000"}, nil)
		}

		ctx := c.Request().Context()

		tx, err := database.Pool().Begin(ctx)
		if err != nil {
			return helpers.SendReturnedHTMLErrorMessage(c, helpers.ErrorMessage{Error: helpers.GenericError{Code: http.StatusBadRequest, Message: err.Error(), UserMessage: "Cannot create this review"}, Box: enums.Boxes.TOAST_TR, Persistance: "3000"}, nil)
		}
		defer database.HandleTransaction(ctx, tx, &err)

		repo := repository.New(tx)

		UUID, err := uuid.Parse(payload.ID)
		if err != nil {
			return helpers.SendReturnedHTMLErrorMessage(c, helpers.ErrorMessage{Error: helpers.GenericError{Code: http.StatusNotFound, Message: err.Error(), UserMessage: "Could not parse invoice ID"}, Box: enums.Boxes.TOAST_TR, Persistance: "3000"}, nil)
		}

		_, err = repo.GetInvoiceByID(ctx, UUID)
		if err != nil {
			return helpers.SendReturnedHTMLErrorMessage(c, helpers.ErrorMessage{Error: helpers.GenericError{Code: http.StatusNotFound, Message: err.Error(), UserMessage: "Cannot find a related invoice with this ID"}, Box: enums.Boxes.TOAST_TR, Persistance: "3000"}, nil)
		}

		review, err := repo.CreateReview(ctx, repository.CreateReviewParams{ID: UUID, Content: payload.Content})
		if err != nil {
			return helpers.SendReturnedHTMLErrorMessage(c, helpers.ErrorMessage{Error: helpers.GenericError{Code: http.StatusNotFound, Message: err.Error(), UserMessage: "A review already exists for this invoice"}, Box: enums.Boxes.TOAST_TR, Persistance: "3000"}, nil)
		}

		html := helpers.MustRenderHTML(components.ReviewItem(review.ID, review.Content, *review.Created))
		html = append(html, helpers.MustRenderHTML(views.EmptyReviewState(false, true))...)

		return c.Blob(http.StatusOK, "text/html; charset=utf-8", html)
	}
}

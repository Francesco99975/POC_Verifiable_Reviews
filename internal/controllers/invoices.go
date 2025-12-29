package controllers

import (
	"net/http"
	"strconv"

	"github.com/Francesco99975/reviews/internal/database"
	"github.com/Francesco99975/reviews/internal/enums"
	"github.com/Francesco99975/reviews/internal/helpers"
	"github.com/Francesco99975/reviews/internal/models"
	"github.com/Francesco99975/reviews/internal/repository"
	"github.com/Francesco99975/reviews/views"
	"github.com/Francesco99975/reviews/views/components"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func InvoicesTab() echo.HandlerFunc {
	return func(c echo.Context) error {

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

		csrf := c.Get("csrf").(string)

		html := helpers.MustRenderHTML(views.InvoicesTab(invoices, csrf))

		return c.Blob(http.StatusOK, "text/html; charset=utf-8", html)
	}
}

func CreateInvoice() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload models.InvoiceDTO

		err := c.Bind(&payload)
		if err != nil {
			return helpers.SendReturnedHTMLErrorMessage(c, helpers.ErrorMessage{Error: helpers.GenericError{Code: http.StatusBadRequest, Message: err.Error(), UserMessage: "Invalid data"}, Box: enums.Boxes.TOAST_TR, Persistance: "3000"}, nil)
		}

		cleanStringTotal := helpers.NormalizeFloatStrToIntStr(payload.Total)

		intTotal, err := strconv.Atoi(cleanStringTotal)
		if err != nil {
			return helpers.SendReturnedHTMLErrorMessage(c, helpers.ErrorMessage{Error: helpers.GenericError{Code: http.StatusBadRequest, Message: err.Error(), UserMessage: "Invalid data"}, Box: enums.Boxes.TOAST_TR, Persistance: "3000"}, nil)
		}

		ctx := c.Request().Context()

		tx, err := database.Pool().Begin(ctx)
		if err != nil {
			return helpers.SendReturnedHTMLErrorMessage(c, helpers.ErrorMessage{Error: helpers.GenericError{Code: http.StatusBadRequest, Message: err.Error(), UserMessage: "Cannot create this invoice"}, Box: enums.Boxes.TOAST_TR, Persistance: "3000"}, nil)
		}
		defer database.HandleTransaction(ctx, tx, &err)

		repo := repository.New(tx)

		invoice, err := repo.CreateInvoice(ctx, repository.CreateInvoiceParams{
			ID:    uuid.New(),
			Total: int32(intTotal),
		})
		if err != nil {
			return helpers.SendReturnedHTMLErrorMessage(c, helpers.ErrorMessage{Error: helpers.GenericError{Code: http.StatusBadRequest, Message: err.Error(), UserMessage: "Cannot create this invoice"}, Box: enums.Boxes.TOAST_TR, Persistance: "3000"}, nil)
		}

		csrf := c.Get("csrf").(string)

		html := helpers.MustRenderHTML(components.InvoiceItem(invoice.ID, invoice.Total, *invoice.Created, nil, nil, csrf))
		html = append(html, helpers.MustRenderHTML(views.EmptyInvoiceState(false, true))...)

		return c.Blob(http.StatusCreated, "text/html; charset=utf-8", html)

	}
}

func DeleteInvoice() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		ctx := c.Request().Context()

		tx, err := database.Pool().Begin(ctx)
		if err != nil {
			return helpers.SendReturnedHTMLErrorMessage(c, helpers.ErrorMessage{Error: helpers.GenericError{Code: http.StatusBadRequest, Message: err.Error(), UserMessage: "Cannot create this invoice"}, Box: enums.Boxes.TOAST_TR, Persistance: "3000"}, nil)
		}
		defer database.HandleTransaction(ctx, tx, &err)

		repo := repository.New(tx)

		log.Debugf("ID: %s", id)

		UUID, err := uuid.Parse(id)
		if err != nil {
			return helpers.SendReturnedHTMLErrorMessage(c, helpers.ErrorMessage{Error: helpers.GenericError{Code: http.StatusNotFound, Message: err.Error(), UserMessage: "Invoice could not be deleted"}, Box: enums.Boxes.TOAST_TR, Persistance: "3000"}, nil)
		}

		_, err = repo.DeleteInvoice(ctx, UUID)
		if err != nil {
			return helpers.SendReturnedHTMLErrorMessage(c, helpers.ErrorMessage{Error: helpers.GenericError{Code: http.StatusNotFound, Message: err.Error(), UserMessage: "Invoice could not be deleted"}, Box: enums.Boxes.TOAST_TR, Persistance: "3000"}, nil)
		}

		invoiceCount, err := repo.CountInvoices(ctx)
		if err != nil {
			return helpers.SendReturnedHTMLErrorMessage(c, helpers.ErrorMessage{Error: helpers.GenericError{Code: http.StatusNotFound, Message: err.Error(), UserMessage: "Invoice could not be deleted, Invoices could not be counted"}, Box: enums.Boxes.TOAST_TR, Persistance: "3000"}, nil)
		}

		html := helpers.MustRenderHTML(views.EmptyInvoiceState(invoiceCount == 0, true))

		return c.Blob(http.StatusOK, "text/html; charset=utf-8", html)
	}
}

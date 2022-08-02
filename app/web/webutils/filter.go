package webutils

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vonmutinda/organono/app/forms"
	"github.com/vonmutinda/organono/app/utils"
)

func FilterFromContext(
	c *gin.Context,
) (*forms.Filter, error) {

	page, per, err := paginationFromContext(c)
	if err != nil {
		return &forms.Filter{}, err
	}

	filter := &forms.Filter{
		Page:   page,
		Per:    per,
		Term:   strings.TrimSpace(c.Query("term")),
		Status: strings.TrimSpace(c.Query("status")),
	}

	return filter, nil
}

func paginationFromContext(
	c *gin.Context,
) (int, int, error) {

	page := 1
	per := 20

	var err error

	pageQueryString := strings.TrimSpace(c.Query("page"))
	if pageQueryString != "" {
		page, err = strconv.Atoi(pageQueryString)
		if err != nil {
			return page, per, utils.NewErrorWithCode(
				err,
				utils.ErrorCodeInvalidArgument,
				"provided invalid page query string = [%v]",
				pageQueryString,
			)
		}
	}

	perQueryString := strings.TrimSpace(c.Query("per"))
	if perQueryString != "" {
		per, err = strconv.Atoi(perQueryString)
		if err != nil {
			return page, per, utils.NewErrorWithCode(
				err,
				utils.ErrorCodeInvalidArgument,
				"provided invalid per query string = [%v]",
				perQueryString,
			)
		}
	}

	return page, per, nil
}

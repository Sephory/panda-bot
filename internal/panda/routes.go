package panda

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

const ROUTE_POLL_RESULTS = "/api/polls/:pollId/poll_results"

func (p *Panda) getRoutes() []echo.Route {
	return []echo.Route{
		{
			Method: http.MethodGet,
			Path:   ROUTE_POLL_RESULTS,
			Handler: func(c echo.Context) error {
				results := p.database.GetResultsForPoll(c.PathParam("pollId"))
				return c.JSON(http.StatusOK, results)
			},
		},
	}
}

func (p *Panda) createRoutes(router *echo.Echo) {
	for _, r := range p.getRoutes() {
		router.AddRoute(r)
	}
}

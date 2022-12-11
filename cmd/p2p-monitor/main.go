package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/PacoXXD/p2p-monitor/pkg/models"
	"github.com/PacoXXD/p2p-monitor/pkg/usecase"
	"github.com/coocood/freecache"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
}

func main() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.00000",
		DataKey:         "p2p-monitor",
		PrettyPrint:     false,
	})
	log.SetLevel(log.InfoLevel)

	cacheSize := 100 * 1024 * 1024
	cache := freecache.NewCache(cacheSize)

	svc := usecase.NewUsecase(cache)
	// Echo instance
	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &MoitorContext{
				usecase: svc,
				Context: c,
			}
			return next(cc)
		}
	})

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/list", ListPeer)
	e.GET("/get", GetPeer)
	e.POST("/report", ReportPeer)

	// Start server
	e.Logger.Fatal(e.Start(":4012"))

}

type MoitorContext struct {
	echo.Context
	usecase usecase.MoitorUsecase
}

func (c *MoitorContext) GetMonitorUsecase() usecase.MoitorUsecase {
	return c.usecase
}

func (c *MoitorContext) SuccessResponse(data interface{}) error {
	return SuccessResponse(c, data)
}

func (c *MoitorContext) FailedResponse(msg string) error {
	return FailedResponse(c, msg)
}

func ListPeer(c echo.Context) error {
	ShareerKey := c.FormValue("share_key")
	log.WithFields(log.Fields{"share_key": ShareerKey}).Info("get one peer")

	ctx := c.(*MoitorContext)
	usecase := ctx.GetMonitorUsecase()

	newCtx, cancel := context.WithTimeout(ctx.Request().Context(), time.Second*2)
	defer cancel()

	peers, err := usecase.ListPeer(newCtx, ShareerKey)
	if err != nil || len(peers) < 1 {
		return ctx.FailedResponse("failed get peer")
	}

	return ctx.SuccessResponse(peers)

}

func ReportPeer(c echo.Context) error {
	PeerPort := c.FormValue("peer_port")
	TrackerURL := c.FormValue("tracker_url")
	ShareerKey := c.FormValue("share_key")
	ChatURL := c.FormValue("chat_url")
	Status := c.FormValue("status")
	PeerStatus := models.PeerStatus(Status)
	PeerIP := c.RealIP()

	log.WithFields(log.Fields{"peer_port": PeerPort, "peer_ip": PeerIP, "tracker_url": TrackerURL, "share_key": ShareerKey, "chat_url": ChatURL, "status": Status}).Info("peer report")

	ctx := c.(*MoitorContext)
	usecase := ctx.GetMonitorUsecase()

	newCtx, cancel := context.WithTimeout(ctx.Request().Context(), time.Second*2)
	defer cancel()

	err := usecase.ReportPeer(newCtx, PeerIP, PeerPort, TrackerURL, ChatURL, ShareerKey, PeerStatus)
	if err != nil {
		log.WithFields(log.Fields{"peer_ip": PeerIP, "peer_port": PeerPort}).Errorf("failed report peer.err:%s", err)
		return ctx.FailedResponse("failed report")
	}
	return ctx.SuccessResponse(nil)
}

func GetPeer(c echo.Context) error {
	ShareerKey := c.QueryParam("share_key")
	log.WithFields(log.Fields{"share_key": ShareerKey}).Info("get one peer")

	ctx := c.(*MoitorContext)
	usecase := ctx.GetMonitorUsecase()

	newCtx, cancel := context.WithTimeout(ctx.Request().Context(), time.Second*2)
	defer cancel()

	peer, err := usecase.GetPeer(newCtx, ShareerKey)
	if err != nil {
		return ctx.FailedResponse(fmt.Sprintf("failed get peer.err:%s", err))
	}

	if peer == nil {
		return ctx.FailedResponse("not found peer")
	}

	return ctx.SuccessResponse(peer)
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func SuccessResponse(c echo.Context, data interface{}) error {
	if data == nil {
		return c.JSON(http.StatusOK, &Response{
			Code: 200,
		})
	}
	return c.JSON(http.StatusOK, &Response{
		Code: 200,
		Data: data,
	})
}

func FailedResponse(c echo.Context, msg string) error {
	return c.JSON(http.StatusInternalServerError, &Response{
		Code:    10000,
		Message: msg,
	})
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

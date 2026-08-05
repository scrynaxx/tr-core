package rest

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/scrynaxx/tr-core/internal/echoutil/middleware"
	v1 "github.com/scrynaxx/tr-core/internal/modules/order/controller/rest/v1"
	"github.com/scrynaxx/tr-core/internal/modules/order/model"
	"github.com/scrynaxx/tr-core/internal/modules/order/usecase"

	"github.com/scrynaxx/tr-core/internal/security"
)

func RegisterRoutes(e *echo.Echo, authorizer security.Authorizer, dictionaryUc usecase.Dictionary) {
	e.Use(middleware.ErrorMapper(map[error]int{
		model.ErrCargoPackageNotFound:    http.StatusNotFound,
		model.ErrOfferingNotFound:        http.StatusNotFound,
		model.ErrPackingMaterialNotFound: http.StatusNotFound,
		model.ErrVehicleNotFound:         http.StatusNotFound,
		model.ErrVehicleExists:           http.StatusConflict,
	}))

	c := v1.New(dictionaryUc)
	adminGroupV1 := e.Group("/v1/admin/order", middleware.Auth(authorizer, security.ActorEmployee))

	adminGroupV1.POST("/dictionary/cargo-packages", c.CreateCargoPackage)
	adminGroupV1.GET("/dictionary/cargo-packages/:package_id", c.GetCargoPackage)
	adminGroupV1.GET("/dictionary/cargo-packages", c.ListCargoPackages)
	adminGroupV1.PUT("/dictionary/cargo-packages/:package_id", c.UpdateCargoPackage)
	adminGroupV1.POST("/dictionary/cargo-packages/:package_id/archive", c.ArchiveCargoPackage)
	adminGroupV1.POST("/dictionary/cargo-packages/:package_id/restore", c.RestoreCargoPackage)

	adminGroupV1.POST("/dictionary/offerings", c.CreateOffering)
	adminGroupV1.GET("/dictionary/offerings/:offering_id", c.GetOffering)
	adminGroupV1.GET("/dictionary/offerings", c.ListOfferings)
	adminGroupV1.PUT("/dictionary/offerings/:offering_id", c.UpdateOffering)
	adminGroupV1.POST("/dictionary/offerings/:offering_id/archive", c.ArchiveOffering)
	adminGroupV1.POST("/dictionary/offerings/:offering_id/restore", c.RestoreOffering)

	adminGroupV1.POST("/dictionary/packing-materials", c.CreatePackingMaterial)
	adminGroupV1.GET("/dictionary/packing-materials/:material_id", c.GetPackingMaterial)
	adminGroupV1.GET("/dictionary/packing-materials", c.ListPackingMaterials)
	adminGroupV1.PUT("/dictionary/packing-materials/:material_id", c.UpdatePackingMaterial)
	adminGroupV1.POST("/dictionary/packing-materials/:material_id/archive", c.ArchivePackingMaterial)
	adminGroupV1.POST("/dictionary/packing-materials/:material_id/restore", c.RestorePackingMaterial)

	adminGroupV1.POST("/dictionary/vehicles", c.CreateVehicle)
	adminGroupV1.GET("/dictionary/vehicles/:vehicle_id", c.GetVehicle)
	adminGroupV1.GET("/dictionary/vehicles", c.ListVehicles)
	adminGroupV1.PUT("/dictionary/vehicles/:vehicle_id", c.UpdateVehicle)
	adminGroupV1.POST("/dictionary/vehicles/:vehicle_id/archive", c.ArchiveVehicle)
	adminGroupV1.POST("/dictionary/vehicles/:vehicle_id/restore", c.RestoreVehicle)
}

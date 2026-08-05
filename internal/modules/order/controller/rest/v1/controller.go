package v1

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/scrynaxx/tr-core/internal/echoutil"
	"github.com/scrynaxx/tr-core/internal/modules/order/controller/rest/v1/contract"
	"github.com/scrynaxx/tr-core/internal/modules/order/model"
	"github.com/scrynaxx/tr-core/internal/modules/order/usecase"
)

type Controller struct {
	dictionaryUc usecase.Dictionary
}

func New(dictionaryUc usecase.Dictionary) *Controller {
	return &Controller{
		dictionaryUc: dictionaryUc,
	}
}

func (co *Controller) CreateCargoPackage(c *echo.Context) error {
	req, err := echoutil.Bind[contract.CreateCargoPackageRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.CreateCargoPackage(c.Request().Context(), req.Name); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (co *Controller) GetCargoPackage(c *echo.Context) error {
	req, err := echoutil.Bind[contract.GetCargoPackageRequest](c)
	if err != nil {
		return err
	}

	cargoPackage, err := co.dictionaryUc.GetCargoPackage(c.Request().Context(), req.PackageID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contract.ToCargoPackage(cargoPackage))
}

func (co *Controller) ListCargoPackages(c *echo.Context) error {
	pack, err := co.dictionaryUc.ListCargoPackages(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contract.ToListCargoPackagesResponse(pack))
}

func (co *Controller) UpdateCargoPackage(c *echo.Context) error {
	req, err := echoutil.Bind[contract.UpdateCargoPackageRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.UpdateCargoPackage(c.Request().Context(), req.PackageID, req.Name); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) ArchiveCargoPackage(c *echo.Context) error {
	req, err := echoutil.Bind[contract.ArchiveCargoPackageRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.ArchiveCargoPackage(c.Request().Context(), req.PackageID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) RestoreCargoPackage(c *echo.Context) error {
	req, err := echoutil.Bind[contract.RestoreCargoPackageRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.RestoreCargoPackage(c.Request().Context(), req.PackageID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) CreateOffering(c *echo.Context) error {
	req, err := echoutil.Bind[contract.CreateOfferingRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.CreateOffering(c.Request().Context(), req.Name, req.Price, req.Modifiers); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (co *Controller) GetOffering(c *echo.Context) error {
	req, err := echoutil.Bind[contract.GetOfferingRequest](c)
	if err != nil {
		return err
	}

	offering, err := co.dictionaryUc.GetOffering(c.Request().Context(), req.OfferingID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contract.ToOffering(offering))
}

func (co *Controller) ListOfferings(c *echo.Context) error {
	offerings, err := co.dictionaryUc.ListOfferings(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contract.ToListOfferingsResponse(offerings))
}

func (co *Controller) UpdateOffering(c *echo.Context) error {
	req, err := echoutil.Bind[contract.UpdateOfferingRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.UpdateOffering(c.Request().Context(), req.OfferingID, req.Name, req.Price, req.Modifiers); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) ArchiveOffering(c *echo.Context) error {
	req, err := echoutil.Bind[contract.ArchiveOfferingRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.ArchiveOffering(c.Request().Context(), req.OfferingID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) RestoreOffering(c *echo.Context) error {
	req, err := echoutil.Bind[contract.RestoreOfferingRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.RestoreOffering(c.Request().Context(), req.OfferingID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) CreatePackingMaterial(c *echo.Context) error {
	req, err := echoutil.Bind[contract.CreatePackingMaterialRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.CreatePackingMaterial(c.Request().Context(), req.Name, req.Price); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (co *Controller) GetPackingMaterial(c *echo.Context) error {
	req, err := echoutil.Bind[contract.GetPackingMaterialRequest](c)
	if err != nil {
		return err
	}

	material, err := co.dictionaryUc.GetPackingMaterial(c.Request().Context(), req.MaterialID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contract.ToPackingMaterial(material))
}

func (co *Controller) ListPackingMaterials(c *echo.Context) error {
	materials, err := co.dictionaryUc.ListPackingMaterials(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contract.ToListPackingMaterialsResponse(materials))
}

func (co *Controller) UpdatePackingMaterial(c *echo.Context) error {
	req, err := echoutil.Bind[contract.UpdatePackingMaterialRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.UpdatePackingMaterial(c.Request().Context(), req.MaterialID, req.Name, req.Price); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) ArchivePackingMaterial(c *echo.Context) error {
	req, err := echoutil.Bind[contract.ArchivePackingMaterialRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.ArchivePackingMaterial(c.Request().Context(), req.MaterialID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) RestorePackingMaterial(c *echo.Context) error {
	req, err := echoutil.Bind[contract.RestorePackingMaterialRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.RestorePackingMaterial(c.Request().Context(), req.MaterialID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) CreateVehicle(c *echo.Context) error {
	req, err := echoutil.Bind[contract.CreateVehicleRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.CreateVehicle(c.Request().Context(), model.VehicleInput{
		Type:               req.Type,
		RegistrationNumber: req.RegistrationNumber,
		LengthMeters:       req.LengthMeters,
		WidthMeters:        req.WidthMeters,
		HeightMeters:       req.HeightMeters,
		CapacityTonnes:     req.CapacityTonnes,
	}); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (co *Controller) GetVehicle(c *echo.Context) error {
	req, err := echoutil.Bind[contract.GetVehicleRequest](c)
	if err != nil {
		return err
	}

	vehicle, err := co.dictionaryUc.GetVehicle(c.Request().Context(), req.VehicleID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contract.ToVehicle(vehicle))
}

func (co *Controller) ListVehicles(c *echo.Context) error {
	vehicles, err := co.dictionaryUc.ListVehicles(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contract.ToListVehiclesResponse(vehicles))
}

func (co *Controller) UpdateVehicle(c *echo.Context) error {
	req, err := echoutil.Bind[contract.UpdateVehicleRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.UpdateVehicle(c.Request().Context(), req.VehicleID, model.VehicleInput{
		Type:               req.Type,
		RegistrationNumber: req.RegistrationNumber,
		LengthMeters:       req.LengthMeters,
		WidthMeters:        req.WidthMeters,
		HeightMeters:       req.HeightMeters,
		CapacityTonnes:     req.CapacityTonnes,
	}); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) ArchiveVehicle(c *echo.Context) error {
	req, err := echoutil.Bind[contract.ArchiveVehicleRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.ArchiveVehicle(c.Request().Context(), req.VehicleID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) RestoreVehicle(c *echo.Context) error {
	req, err := echoutil.Bind[contract.RestoreVehicleRequest](c)
	if err != nil {
		return err
	}

	if err = co.dictionaryUc.RestoreVehicle(c.Request().Context(), req.VehicleID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

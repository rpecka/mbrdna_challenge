package requests

type GetVehiclesResponse []Vehicle

type Vehicle struct {
	ID string `json:"id"`
	LicensePlate string `json:"licenseplate"`
	FINOrVIN string `json:"finorvin"`
}

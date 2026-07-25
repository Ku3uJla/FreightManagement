package filters

type DriverFilter struct {
	Status   *int
	Category *string
}

type AutoFilter struct {
	Capacity        *int
	LiftingCapacity *int
	Status          *int
}

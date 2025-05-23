package packagebase

type PackageImplementor interface {
	// Setters
	SetPorts(ports []int32) 	bool

	// Lifecycle
	Install() 	bool
	Run() 		bool
	Exit()
}

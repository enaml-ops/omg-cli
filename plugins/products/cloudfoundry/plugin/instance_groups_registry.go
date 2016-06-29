package cloudfoundry

// RegisterInstanceGrouperFactory registers an InstanceGrouperFactory.
// InstanceGrouperFactories should generally be registered in their package's
// init() function.
func RegisterInstanceGrouperFactory(igf InstanceGrouperFactory) {
	factories = append(factories, igf)
}

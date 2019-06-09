package conf

// OperatorsTable maps binary operators to corresponding list of functions.
// Functions should be provided in the environment to allow operator overloading.
type OperatorsTable map[string][]string

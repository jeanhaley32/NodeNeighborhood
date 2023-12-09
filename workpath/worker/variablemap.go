package worker

// Variable Map, used to pass variables to the task.

type (
	// Variable Map, used to pass variables to the task.
	VMap map[string]any
)

// Adds an entry to the variable map.
func (v *VMap) Set(key string, value any) {
	(*v)[key] = value
}

// Returns the value associated with the key.
func (v *VMap) Get(key string) any {
	return (*v)[key]
}

// Removes the entry associated with the key.
func (v *VMap) Delete(key string) {
	delete(*v, key)
}

// Appends multiple entries to the variable map.
func (v *VMap) SetKeys(m map[string]any) {
	for key, val := range m {
		(*v)[key] = val
	}
}

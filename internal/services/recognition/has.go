package recognition

func (r *recognition) Has(id string) bool {
	return r.taskManager.Get(id) != nil
}

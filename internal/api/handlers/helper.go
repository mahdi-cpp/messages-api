package handlers

type BuckUpdateChats struct {
	Ids []string `json:"ids"`
}

type BuckDeleteChats struct {
	Ids []string `json:"ids"`
}

package model

func Initial(tapEndpoint string, opAmpEndpoint string) model {
	tabs := []string{"Metrics", "Traces", "Logs"}
	return model{
		tapEndpoint:   tapEndpoint,
		opAmpEndpoint: opAmpEndpoint,
		tabs:          tabs,
		channel:       make(chan []byte),
	}
}

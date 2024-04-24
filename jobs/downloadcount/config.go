package downloadcount

type Config struct {
	Spec     string `json:"spec"     required:"true"`
	Endpoint string `json:"endpoint" required:"true"`
}

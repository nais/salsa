package commands

type Cfg struct {
    AttestCfg `mapstructure:",squash"`
}

type AttestCfg struct {
    KmsUrl string `mapstructure:"kmsUrl"`
}

package storage

type Registry interface {
	SaveKatanState(data []byte) error
	RestoreKatanState() ([]byte, error)
}

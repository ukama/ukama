package pkg

type MemStorage struct {
	data map[string]string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]string),
	}
}

func (m MemStorage) Get(key string) ([]byte, error) {
	if val, ok := m.data[key]; ok {
		return []byte(val), nil
	}

	return nil, nil
}

func (m MemStorage) Put(key string, value string) error {
	m.data[key] = value
	return nil
}

func (m MemStorage) Delete(key string) error {
	delete(m.data, key)
	return nil
}

package storage

type MemStorage struct {
	data map[string]simInfo
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]simInfo),
	}
}

func (m MemStorage) Get(key string) (simInfo, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}

	return simInfo{}, nil
}

func (m MemStorage) Put(key string, value simInfo) error {
	m.data[key] = value
	return nil
}

func (m MemStorage) Delete(key string) error {
	delete(m.data, key)
	return nil
}

package domain

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"github.com/google/uuid"
)

// Item represents the main domain entity
type Item struct {
	ID        string
	Name      string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewItem creates a new instance of Item
func NewItem(name string, active bool) Item {
	return Item{
		ID:        generateID(),
		Name:      name,
		Active:    active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (i Item) IsActive() bool {
	return i.Active
}

func (i Item) IsEmpty() bool {
	return i.ID == ""
}

func generateID() string {
	// MongoDB ObjectID tem 12 bytes.
	// Geramos 12 bytes aleatórios e os convertemos para uma string hexadecimal de 24 caracteres.
	objectIDBytes := make([]byte, 12)
	_, err := rand.Read(objectIDBytes)
	if err != nil {
		log.Printf("Error generating random bytes for ID: %v", err)
		// Em caso de erro, podemos retornar uma string vazia ou um ID de fallback, dependendo da política.
		// Por simplicidade aqui, retornamos uma string vazia e lidamos com isso no nível superior se necessário.
		return uuid.New().String()
	}

	return hex.EncodeToString(objectIDBytes)
}

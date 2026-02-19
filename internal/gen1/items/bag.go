package items

import (
	"fmt"

	"github.com/abravonunez/raracandy/internal/gen1/save"
)

const (
	MaxBagItems = 20
	MaxItemQty  = 99
)

// Item represents an item in the bag
type Item struct {
	ID       byte
	Quantity byte
	Name     string
}

// GetBagItems returns all items currently in the bag
func GetBagItems(s *save.Save) []Item {
	count := s.GetByte(save.OffsetBagCount)
	if count > MaxBagItems {
		count = MaxBagItems
	}

	items := make([]Item, 0, count)
	offset := save.OffsetBagItems

	for i := byte(0); i < count; i++ {
		id := s.GetByte(offset)
		qty := s.GetByte(offset + 1)

		items = append(items, Item{
			ID:       id,
			Quantity: qty,
			Name:     GetItemName(id),
		})

		offset += 2 // Each item takes 2 bytes (ID + quantity)
	}

	return items
}

// FindItemIndex finds the index of an item in the bag by ID
// Returns -1 if not found
func FindItemIndex(s *save.Save, itemID byte) int {
	count := s.GetByte(save.OffsetBagCount)
	offset := save.OffsetBagItems

	for i := byte(0); i < count; i++ {
		id := s.GetByte(offset)
		if id == itemID {
			return int(i)
		}
		offset += 2
	}

	return -1
}

// SetItemQuantity updates the quantity of an existing item in the bag
// If the item doesn't exist, it will be added to the bag
func SetItemQuantity(s *save.Save, itemID byte, quantity byte) error {
	if quantity > MaxItemQty {
		return fmt.Errorf("quantity %d exceeds maximum %d", quantity, MaxItemQty)
	}

	// Check if item exists
	idx := FindItemIndex(s, itemID)

	if idx >= 0 {
		// Item exists, update quantity
		offset := save.OffsetBagItems + (idx * 2) + 1
		return s.SetByte(offset, quantity)
	}

	// Item doesn't exist, add it
	return AddItem(s, itemID, quantity)
}

// AddItem adds a new item to the bag
func AddItem(s *save.Save, itemID byte, quantity byte) error {
	if quantity > MaxItemQty {
		return fmt.Errorf("quantity %d exceeds maximum %d", quantity, MaxItemQty)
	}

	count := s.GetByte(save.OffsetBagCount)

	if count >= MaxBagItems {
		return fmt.Errorf("bag is full (max %d items)", MaxBagItems)
	}

	// Calculate offset for new item (after last item)
	offset := save.OffsetBagItems + (int(count) * 2)

	// Set item ID and quantity
	if err := s.SetByte(offset, itemID); err != nil {
		return fmt.Errorf("failed to set item ID: %w", err)
	}
	if err := s.SetByte(offset+1, quantity); err != nil {
		return fmt.Errorf("failed to set item quantity: %w", err)
	}

	// Increment bag count
	if err := s.SetByte(save.OffsetBagCount, count+1); err != nil {
		return fmt.Errorf("failed to update bag count: %w", err)
	}

	// Add terminator byte (0xFF) after the new item
	if err := s.SetByte(offset+2, 0xFF); err != nil {
		return fmt.Errorf("failed to set terminator byte: %w", err)
	}

	return nil
}

// RemoveItem removes an item from the bag by ID
func RemoveItem(s *save.Save, itemID byte) error {
	idx := FindItemIndex(s, itemID)
	if idx < 0 {
		return fmt.Errorf("item not found in bag")
	}

	count := s.GetByte(save.OffsetBagCount)

	// Shift all items after the removed one
	for i := idx; i < int(count)-1; i++ {
		srcOffset := save.OffsetBagItems + ((i + 1) * 2)
		dstOffset := save.OffsetBagItems + (i * 2)

		id := s.GetByte(srcOffset)
		qty := s.GetByte(srcOffset + 1)

		s.SetByte(dstOffset, id)
		s.SetByte(dstOffset+1, qty)
	}

	// Decrement count
	s.SetByte(save.OffsetBagCount, count-1)

	// Add terminator at new end
	terminatorOffset := save.OffsetBagItems + (int(count-1) * 2)
	s.SetByte(terminatorOffset, 0xFF)

	return nil
}

package liborder

const (
	OrderAscending = "ascending"
	OrderDescending = "descending"
)

func ValidOrder(order string) bool {
	switch order{
	case OrderAscending:
		return true
	case OrderDescending:
		return true
	default:
		return false
	}
}

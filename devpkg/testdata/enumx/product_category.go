package enumx

// ProductCategory 产品类别
// +genx:enum
type ProductCategory int8

const (
	PRODUCT_CATEGORY_UNKNOWN ProductCategory = iota
	// PRODUCT_CATEGORY__CLOTHING 服饰
	PRODUCT_CATEGORY__CLOTHING
	// PRODUCT_CATEGORY__3C 3C数码
	PRODUCT_CATEGORY__3C
)

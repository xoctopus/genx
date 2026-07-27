package enumx

// ProductType 产品类型
// +genx:enum
// @enum map.Category=ProductCategory
type ProductType int8

const (
	PRODUCT_TYPE_UNKNOWN ProductType = iota
	// PRODUCT_TYPE__MEN_CLOTH 男装
	// @enum map.Category=PRODUCT_CATEGORY__CLOTHING
	PRODUCT_TYPE__MEN_CLOTH
	// PRODUCT_TYPE__SMART_PHONE 智能手机
	// @enum map.Category=PRODUCT_CATEGORY__3C
	PRODUCT_TYPE__SMART_PHONE
)

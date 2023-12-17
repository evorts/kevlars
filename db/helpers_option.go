/**
 * @Author: steven
 * @Description:
 * @File: helpers_option
 * @Date: 18/12/23 00.10
 */

package db

type IHelperOption interface {
	apply(h *helper)
}

type helperOption func(h *helper)

func (o helperOption) apply(h *helper) {
	o(h)
}

func WithPagination(page, limit int) IHelperOption {
	return helperOption(func(h *helper) {
		h.pagination = &Pagination{
			Page:  page,
			Limit: limit,
		}
		h.pagination.offset = h.pagination.calcOffset()
	})
}

func WithOrdersBy(v OrdersBy) IHelperOption {
	return helperOption(func(h *helper) {
		h.ordersBy = v
	})
}

func WithOrderBy(v OrderBy) IHelperOption {
	return helperOption(func(h *helper) {
		h.ordersBy = append(h.ordersBy, v)
	})
}

func WithFilters(v Filters) IHelperOption {
	return helperOption(func(h *helper) {
		h.filters = &v
	})
}

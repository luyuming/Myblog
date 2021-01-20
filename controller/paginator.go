package controller

import "math"

//分页器结构
type paginator struct {
	Total     int   //记录总数
	PageSize  int   //每页大小
	PageTotal int   //总页数
	PageNow   int   //当前页数
	LastPage  int   //上一页
	NextPage  int   //下一页
	PageNums  []int //显示页码
}

var defaultPageSize = 2 //默认页大小
var pageNum = 2         //显示页码数量

//创建分页器
func CreatePaginator(pageNow, pageSize, total int) paginator {
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	pager := &paginator{
		Total:     total,
		PageSize:  pageSize,
		PageTotal: int(math.Ceil(float64(total) / float64(pageSize))),
		PageNow:   pageNow,
	}
	if total <= 0 {
		pager.PageTotal = 1
		pager.PageNow = 1
		pager.LastPage = 1
		pager.NextPage = 1
		pager.PageNums = append(pager.PageNums, 1)
		return *pager
	}
	//分页边界处理
	if pager.PageNow > pager.PageTotal {
		pager.PageNow = pager.PageTotal
	} else if pager.PageNow < 1 {
		pager.PageNow = 1
	}
	//上一页与下一页
	pager.LastPage = pager.PageNow
	pager.NextPage = pager.PageNow
	if pager.PageNow > 1 {
		pager.LastPage = pager.PageNow - 1
	}
	if pager.PageNow < pager.PageTotal {
		pager.NextPage = pager.PageNow + 1
	}
	//显示页码
	var start, end int //开始页码与结束页码
	if pager.PageTotal <= pageNum {
		start = 1
		end = pager.PageTotal
	} else {
		before := pageNum / 2         //当前页前面页码数
		after := pageNum - before - 1 //当前页后面页码数
		start = pager.PageNow - before
		end = pager.PageNow + after
		if start < 1 { //当前页前面页码数不足
			start = 1
			end = pageNum
		} else if end > pager.PageTotal { //当前页后面页码数不足
			start = pager.PageTotal - pageNum + 1
			end = pager.PageTotal
		}
	}
	for i := start; i <= end; i++ {
		pager.PageNums = append(pager.PageNums, i)
	}
	return *pager
}

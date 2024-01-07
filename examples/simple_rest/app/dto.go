/**
 * @Author: steven
 * @Description:
 * @File: dto
 * @Date: 08/01/24 06.26
 */

package app

type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTodoRequest struct {
	CreateTodoRequest
	Id int `json:"id"`
}

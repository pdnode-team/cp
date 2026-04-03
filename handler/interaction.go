package handler

import (
	"cp-website/ent"
	"cp-website/ent/comment"
	"cp-website/ent/cp"
	"cp-website/ent/user"
	"cp-website/model"

	"github.com/labstack/echo/v4"
)

// ToggleLike 处理 POST /cp/:id/like
func ToggleLike(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		dbUser := c.Get("user").(*ent.User)
		cpID := ParseID(c.Param("id"))

		// 检查用户是否已经点赞过这个 CP
		liked, err := client.CP.Query().
			Where(cp.ID(cpID), cp.HasLikedByUsersWith(user.ID(dbUser.ID))).
			Exist(ctx)
		if err != nil {
			return err
		}

		if liked {
			// 如果已经点赞，则移除关联 (取消点赞)
			err = client.CP.UpdateOneID(cpID).RemoveLikedByUserIDs(dbUser.ID).Exec(ctx)
		} else {
			// 如果没点赞，则添加关联 (点赞)
			err = client.CP.UpdateOneID(cpID).AddLikedByUserIDs(dbUser.ID).Exec(ctx)
		}

		if err != nil {
			return err
		}

		return Success(c, map[string]interface{}{
			"liked": !liked, // 返回当前最终的状态
		}, 200)
	}
}

// CreateComment 处理 POST /cp/:id/comment
func CreateComment(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		dbUser := c.Get("user").(*ent.User)
		cpID := ParseID(c.Param("id"))

		req := new(model.CommentReq)
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(400, err.Error())
		}
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(400, err.Error())
		}

		// 开始构建评论
		builder := client.Comment.Create().
			SetContent(req.Content).
			SetAuthor(dbUser).
			SetCpID(cpID) // Ent 会自动基于边生成 SetCpID 方法

		// 如果前端传了 parent_id，说明这是一条回复
		if req.ParentID != nil {
			builder.SetParentID(*req.ParentID)
		}

		newComment, err := builder.Save(ctx)
		if err != nil {
			return err
		}

		return Success(c, newComment, 201)
	}
}

// GetComments 处理 GET /cp/:id/comment
func GetComments(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		cpID := ParseID(c.Param("id"))

		// 这里的查询十分强大：
		// 我们只查顶层评论 (没有 Parent 的评论)，然后通过 WithChildren 把回复一起带出来！
		comments, err := client.Comment.Query().
			Where(
				comment.HasCpWith(cp.ID(cpID)),
				comment.Not(comment.HasParent()), // 只查顶层评论
			).
			WithAuthor(). // 顺便把发评论的人的信息带上
			WithChildren(func(q *ent.CommentQuery) {
				q.WithAuthor() // 把回复的回复者的信息也带上
			}).
			Order(ent.Desc(comment.FieldCreatedAt)). // 按时间倒序
			All(ctx)

		if err != nil {
			return err
		}

		return Success(c, comments, 200)
	}
}

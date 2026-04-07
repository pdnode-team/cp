import { CommentSchema } from '#database/schema'
import { belongsTo } from '@adonisjs/lucid/orm'
import User from '#models/user'
import Cp from '#models/cp'
import type { BelongsTo } from '@adonisjs/lucid/types/relations'

export default class Comment extends CommentSchema {
  // 关联发布评论的用户
  @belongsTo(() => User, {
    foreignKey: 'userId', 
  })
  declare user: BelongsTo<typeof User>

  // 关联被评论的 CP
  @belongsTo(() => Cp, {
    foreignKey: 'cpId',
  })
  declare cp: BelongsTo<typeof Cp>
}
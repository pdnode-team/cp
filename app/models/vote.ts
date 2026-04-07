import { VoteSchema } from '#database/schema'
import { belongsTo } from '@adonisjs/lucid/orm'
import type { BelongsTo } from '@adonisjs/lucid/types/relations'
import User from '#models/user'
import Cp from '#models/cp'

export default class Vote extends VoteSchema {
  // 关联投票人
  @belongsTo(() => User, {
    foreignKey: 'userId',
  })
  declare user: BelongsTo<typeof User>

  // 关联获票的 CP
  @belongsTo(() => Cp, {
    foreignKey: 'cpId',
  })
  declare cp: BelongsTo<typeof Cp>
}
import { CpSchema } from '#database/schema'
import { belongsTo, hasMany } from '@adonisjs/lucid/orm'
import Comment from '#models/comment'
import Character from '#models/character'
import Vote from '#models/vote'


import type { BelongsTo, HasMany } from '@adonisjs/lucid/types/relations'

export default class Cp extends CpSchema {
    // 关联第一个角色
  @belongsTo(() => Character, {
    foreignKey: 'charOneId', // 这里的名称必须匹配你 Schema 里的变量名
  })
  declare charOne: BelongsTo<typeof Character>

  // 关联第二个角色
  @belongsTo(() => Character, {
    foreignKey: 'charTwoId',
  })
  declare charTwo: BelongsTo<typeof Character>

  // 获取该 CP 下的所有评论
  @hasMany(() => Comment, {
    foreignKey: 'cpId',
  })
  declare comments: HasMany<typeof Comment>

  // 获取该 CP 的所有投票记录
  @hasMany(() => Vote, {
    foreignKey: 'cpId',
  })
  declare votes: HasMany<typeof Vote>
}
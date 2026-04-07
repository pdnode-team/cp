import { CharacterSchema } from '#database/schema'
import Cp from '#models/cp'
import User from '#models/user'
import type { BelongsTo } from '@adonisjs/lucid/types/relations'
import { belongsTo, hasMany } from '@adonisjs/lucid/orm'
import type { HasMany } from '@adonisjs/lucid/types/relations'

export default class Character extends CharacterSchema {
  // 作为“角色1”参与的所有 CP
  @hasMany(() => Cp, {
    foreignKey: 'charOneId',
  })
  declare cpsAsPrimary: HasMany<typeof Cp>

  // 作为“角色2”参与的所有 CP
  @hasMany(() => Cp, {
    foreignKey: 'charTwoId',
  })
  declare cpsAsSecondary: HasMany<typeof Cp>

  @belongsTo(() => User)
  declare user: BelongsTo<typeof User>
}
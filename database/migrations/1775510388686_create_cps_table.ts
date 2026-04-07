import { BaseSchema } from '@adonisjs/lucid/schema'

export default class extends BaseSchema {
  protected tableName = 'cps'

  async up() {
    this.schema.createTable(this.tableName, (table) => {
      table.increments('id')
      table.string('name').notNullable()
      table.text('description').notNullable()
      table.string('image').notNullable()
      table.integer('char_one_id').unsigned().references('characters.id').onDelete('CASCADE')
      table.integer('char_two_id').unsigned().references('characters.id').onDelete('CASCADE')
      table.integer('vote_count').notNullable().defaultTo(0)
      table.string('tag_names').nullable()
      table.timestamp('created_at', { useTz: true })
      table.timestamp('updated_at', { useTz: true })
    })
  }

  async down() {
    this.schema.dropTable(this.tableName)
  }
}
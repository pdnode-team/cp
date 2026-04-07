import type { HttpContext } from '@adonisjs/core/http'
import Cp from '#models/cp' // 导入你的 CP 模型
import User from '#models/user'
// import db from '@adonisjs/lucid/services/db'
// import Character from '#models/character'

export default class CpsController {
  /**
   * 展示所有 CP 列表 (首页或列表页)
   * 对应路由: GET /
   */
  async index({ view }: HttpContext) {
    // 使用 Lucid 的魔法：预加载关联的角色信息
    // 这样你才能在页面上显示“角色A x 角色B”
    const cps = await Cp.query()
      .preload('charOne')
      .preload('charTwo')
      .orderBy('vote_count', 'desc') // 按票数降序排列，热门的在前

    return view.render('pages/explore', { cps, test:"test" })
  }

  async create({ view }: HttpContext) {

    const user = await User.query()
        .where('id', 1)
        .preload('characters') // 这里的名字必须匹配模型里 declare 的名字
        .firstOrFail()
    return view.render('pages/cp_create', { characters: user.characters })
  }

  /**
   * 展示单个 CP 详情页
   * 对应路由: GET /cps/:id
   */
  async show({ params, view }: HttpContext) {
    const cp = await Cp.query()
      .where('id', params.id)
      .preload('charOne')
      .preload('charTwo')
      .preload('comments', (query) => {
        query.preload('user') // 顺便加载评论的作者
      })
      .firstOrFail()

    return view.render('pages/cp_detail', { cp })
  }
}